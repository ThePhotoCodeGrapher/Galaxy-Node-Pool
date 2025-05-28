# Security Policy

## Supported Versions

Only the latest major version of Galaxy Node Pool receives security updates. We strongly recommend all users to run the latest version to ensure they have all security fixes.

## Reporting a Vulnerability

**Please do not report security vulnerabilities through public GitHub issues.**

### How to Report a Vulnerability

If you believe you've found a security vulnerability in Galaxy Node Pool, we encourage you to let us know right away. We will investigate all legitimate reports and do our best to quickly fix the problem.

1. **Email Us**: Send an email to security@hybridconnect.cloud with the subject line "Galaxy Node Pool Security Vulnerability - [brief description]"
2. **Encryption**: For sensitive reports, please use our PGP key (see below)
3. **Response Time**: We will acknowledge your email within 48 hours
4. **Disclosure**: After the issue has been resolved, we will coordinate public disclosure with you

### PGP Key

Our security team's PGP key is available at [https://hybridconnect.cloud/security.asc](https://hybridconnect.cloud/security.asc)

Fingerprint: `1234 5678 90AB CDEF 1234  5678 90AB CDEF 1234 5678`

## Security Updates

Security updates will be released as patch versions (e.g., 1.2.3) and will be clearly marked in the release notes. We recommend subscribing to our security advisory mailing list for notifications.

## Security Best Practices

### For Users

- Always run the latest version of Galaxy Node Pool
- Use strong, unique passwords for all accounts
- Enable two-factor authentication where available
- Regularly back up your data
- Follow the principle of least privilege when setting up permissions

### For Developers

- Never commit secrets or credentials to version control
- Use environment variables for sensitive configuration
- Validate all user inputs
- Use prepared statements for database queries
- Keep dependencies up to date
- Run security scanners regularly

## Responsible Disclosure Policy

We follow responsible disclosure guidelines:

1. Notify us as soon as possible upon discovery of a potential security issue
2. Provide us with a reasonable amount of time to resolve the issue before any disclosure
3. Make a good faith effort to avoid privacy violations, destruction of data, and service interruptions
4. Do not exploit the vulnerability for financial gain

## Bug Bounty Program

We appreciate the efforts of security researchers who help us keep our users safe. While we don't currently have a formal bug bounty program, we may offer rewards for particularly significant security reports at our discretion.

## Security Audits

Galaxy Node Pool undergoes regular security audits by independent third parties. The results of these audits are made available to our enterprise customers.

## Contact

For security-related inquiries, please contact security@hybridconnect.cloud

---

**AI-ID**: CP-GAL-NODEPOOL-001  
**Documentation**: [https://docs.hybridconnect.cloud/galaxy-node-pool/security](https://docs.hybridconnect.cloud/galaxy-node-pool/security)  
**Maintained by**: [Castle Palette Cloud A.I.](https://hybridconnect.cloud)
