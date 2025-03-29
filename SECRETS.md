# Required GitHub Secrets

This document lists all the secrets required for the SentinelStacks CI/CD workflows to function properly.

## Release Workflow Secrets

### GPG Signing
- `GPG_PRIVATE_KEY`: The private GPG key for signing releases
- `GPG_PASSPHRASE`: The passphrase for the GPG key

### AWS Deployment
- `AWS_ACCESS_KEY_ID`: AWS access key for S3 and CloudFront operations
- `AWS_SECRET_ACCESS_KEY`: AWS secret key for S3 and CloudFront operations
- `AWS_S3_BUCKET`: The name of the S3 bucket for UI deployments
- `AWS_CLOUDFRONT_ID`: The CloudFront distribution ID for cache invalidation

## Testing Secrets

### Database Credentials (for E2E tests)
- `DB_USER`: Database user for testing (default: test)
- `DB_PASSWORD`: Database password for testing (default: test)

## Development Secrets

### API Keys (for development and testing)
- `OPENAI_API_KEY`: OpenAI API key for development and testing
- `ANTHROPIC_API_KEY`: Anthropic API key for development and testing

## How to Add Secrets

1. Go to your GitHub repository
2. Navigate to Settings > Secrets and variables > Actions
3. Click "New repository secret"
4. Add each secret with its corresponding value

## Security Notes

- Never commit these secrets to the repository
- Rotate secrets periodically
- Use GitHub's environment protection rules for production secrets
- Consider using OIDC for AWS authentication in production 