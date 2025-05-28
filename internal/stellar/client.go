package stellar

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/protocols/horizon"
	"github.com/stellar/go/txnbuild"
)

// StellarClient handles interactions with the Stellar network
type StellarClient struct {
	client      *horizonclient.Client
	networkPass string
	poolAccount *keypair.Full
}

// NewStellarClient creates a new Stellar client
func NewStellarClient(horizonURL, networkPassphrase, poolSeed string) (*StellarClient, error) {
	client := &horizonclient.Client{
		HorizonURL: horizonURL,
	}

	// Load pool account from seed
	poolAccount, err := keypair.ParseFull(poolSeed)
	if err != nil {
		return nil, fmt.Errorf("failed to parse pool seed: %v", err)
	}

	return &StellarClient{
		client:      client,
		networkPass: networkPassphrase,
		poolAccount: poolAccount,
	}, nil
}

// RegisterPoolWithMainNet registers this pool with the main net via a Stellar transaction
func (s *StellarClient) RegisterPoolWithMainNet(mainNetAccount, poolDomain string, fee string) error {
	// Load pool account
	poolAccount, err := s.loadAccount(s.poolAccount.Address())
	if err != nil {
		return err
	}

	// Create transaction
	tx, err := txnbuild.NewTransaction(
		txnbuild.TransactionParams{
			SourceAccount:        poolAccount,
			IncrementSequenceNum: true,
			Operations: []txnbuild.Operation{
				&txnbuild.Payment{
					Destination: mainNetAccount,
					Amount:      fee,
					Asset:       txnbuild.NativeAsset{},
				},
			},
			Memo:       txnbuild.MemoText(fmt.Sprintf("GALAXY_POOL_REG:%s", poolDomain)),
			Timebounds: txnbuild.NewTimeout(300),
			BaseFee:    txnbuild.MinBaseFee,
		},
	)
	if err != nil {
		return err
	}

	// Sign transaction
	tx, err = tx.Sign(s.networkPass, s.poolAccount)
	if err != nil {
		return err
	}

	// Submit transaction
	resp, err := s.client.SubmitTransaction(tx)
	if err != nil {
		return err
	}

	log.Printf("Pool registered with main net. Transaction ID: %s", resp.ID)
	return nil
}

// ProcessNodeRegistrationFee processes a registration fee from a node
func (s *StellarClient) ProcessNodeRegistrationFee(nodeAccount, nodeID string, fee string) error {
	// Check if payment was received
	payments, err := s.getPaymentsTo(s.poolAccount.Address(), nodeAccount)
	if err != nil {
		return err
	}

	// Verify payment
	for _, payment := range payments {
		if payment.Type == "payment" && payment.Amount == fee {
			log.Printf("Registration fee received from node %s (account: %s)", nodeID, nodeAccount)
			return nil
		}
	}

	return fmt.Errorf("registration fee not received from node %s (account: %s)", nodeID, nodeAccount)
}

// DistributeStakerRewards distributes rewards to stakers
func (s *StellarClient) DistributeStakerRewards(stakerAccounts []string, rewardPerStaker string) error {
	// Load pool account
	poolAccount, err := s.loadAccount(s.poolAccount.Address())
	if err != nil {
		return err
	}

	// Create operations for each staker
	operations := make([]txnbuild.Operation, 0, len(stakerAccounts))
	for _, staker := range stakerAccounts {
		operations = append(operations, &txnbuild.Payment{
			Destination: staker,
			Amount:      rewardPerStaker,
			Asset:       txnbuild.NativeAsset{},
		})
	}

	// Create transaction
	tx, err := txnbuild.NewTransaction(
		txnbuild.TransactionParams{
			SourceAccount:        poolAccount,
			IncrementSequenceNum: true,
			Operations:           operations,
			Memo:                 txnbuild.MemoText("GALAXY_POOL_REWARDS"),
			Timebounds:           txnbuild.NewTimeout(300),
			BaseFee:              txnbuild.MinBaseFee,
		},
	)
	if err != nil {
		return err
	}

	// Sign transaction
	tx, err = tx.Sign(s.networkPass, s.poolAccount)
	if err != nil {
		return err
	}

	// Submit transaction
	resp, err := s.client.SubmitTransaction(tx)
	if err != nil {
		return err
	}

	log.Printf("Rewards distributed to %d stakers. Transaction ID: %s", len(stakerAccounts), resp.ID)
	return nil
}

// VerifyNodeIdentity verifies a node's identity using Stellar signatures
func (s *StellarClient) VerifyNodeIdentity(nodeID, nodeAccount, signature, challenge string) (bool, error) {
	// Get the node's public key
	kp, err := keypair.Parse(nodeAccount)
	if err != nil {
		return false, err
	}

	// Verify the signature
	valid := kp.Verify([]byte(challenge), []byte(signature))
	return valid, nil
}

// Helper function to load an account
func (s *StellarClient) loadAccount(address string) (*horizon.Account, error) {
	account, err := s.client.AccountDetail(horizonclient.AccountRequest{
		AccountID: address,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to load account %s: %v", address, err)
	}
	return &account, nil
}

// Helper function to get payments to an account
func (s *StellarClient) getPaymentsTo(recipient, sender string) ([]horizon.Payment, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	payments, err := s.client.Payments(horizonclient.PaymentRequest{
		ForAccount:    recipient,
		Cursor:        "now",
		Order:         horizonclient.OrderDesc,
		Limit:         10,
		IncludeFailed: false,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get payments: %v", err)
	}

	var result []horizon.Payment
	for _, payment := range payments.Embedded.Records {
		if payment.From == sender {
			result = append(result, payment)
		}
	}

	return result, nil
}
