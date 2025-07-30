# Security Policy

## Supported Versions

Use this section to tell people about which versions of your project are currently being supported with security updates.

| Version | Supported          |
| ------- | ------------------ |
| 1.0.x   | :white_check_mark: |
| < 1.0   | :x:                |

## Reporting a Vulnerability

We take security vulnerabilities seriously. If you discover a security vulnerability in Slack Channel Cleaner, please follow these steps:

### 1. **DO NOT** create a public GitHub issue
Security vulnerabilities should be reported privately to avoid potential exploitation.

### 2. Email the maintainer
Send an email to the maintainer with the following information:
- **Subject**: `[SECURITY] Vulnerability in Slack Channel Cleaner`
- **Description**: Detailed description of the vulnerability
- **Steps to reproduce**: Clear steps to reproduce the issue
- **Impact**: Potential impact of the vulnerability
- **Suggested fix**: If you have ideas for fixing the issue

### 3. Response timeline
- **Initial response**: Within 48 hours
- **Assessment**: Within 1 week
- **Fix timeline**: Depends on severity and complexity

### 4. Disclosure
- Security vulnerabilities will be disclosed through GitHub Security Advisories
- Patches will be released as soon as possible
- Users will be notified through releases and documentation updates

## Security Best Practices

### For Users
1. **Keep your Slack token secure**
   - Never commit your `.env` file to version control
   - Use environment variables in production
   - Rotate tokens regularly

2. **Review skip lists**
   - Ensure important channels are in your skip list
   - Regularly review and update the skip list

3. **Use the latest version**
   - Always use the latest stable release
   - Subscribe to security notifications

### For Contributors
1. **Code review**
   - All code changes are reviewed for security implications
   - Pay special attention to API integrations and user input

2. **Dependencies**
   - Keep dependencies updated
   - Monitor for known vulnerabilities
   - Use `go mod tidy` and `go mod verify`

3. **Testing**
   - Test with various input scenarios
   - Include security-focused test cases
   - Test error handling and edge cases

## Known Security Considerations

### Workspace API Token
- The application requires a workspace API token with specific scopes
- Tokens are stored in environment variables or `.env` files
- Never log or display tokens in the application

### Rate Limiting
- The application handles workspace API rate limits automatically
- This prevents potential API abuse and account suspension

### Input Validation
- All user inputs are validated before processing
- Configuration files are validated for correct format and values

### File Permissions
- Configuration files should have appropriate permissions
- `.env` files should be readable only by the user running the application

## Security Features

### Built-in Protections
- **Skip List**: Protects important channels from accidental removal
- **Confirmation**: Double-confirmation for destructive actions
- **Error Handling**: Graceful handling of API errors and rate limits
- **Input Validation**: Validation of all configuration and user inputs

### API Security
- **Rate Limiting**: Automatic handling of workspace API rate limits
- **Token Security**: Secure token storage and usage
- **Error Handling**: Proper error handling without exposing sensitive information

## Reporting Security Issues

If you find a security vulnerability, please report it to:

**Email**: [Your email address]
**Subject**: `[SECURITY] Vulnerability in Slack Channel Cleaner`

Please include:
- Detailed description of the vulnerability
- Steps to reproduce
- Potential impact
- Suggested fix (if available)

## Security Updates

Security updates will be:
- Released as patch versions (e.g., 1.0.1, 1.0.2)
- Announced through GitHub releases
- Documented in the changelog
- Tagged with security advisories when appropriate

---

Thank you for helping keep Slack Channel Cleaner secure! 🔒 