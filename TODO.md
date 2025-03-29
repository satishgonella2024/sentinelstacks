# SentinelStacks TODO List

## CI/CD Tasks

### GitHub Actions Setup
- [ ] Add required secrets to GitHub repository:
  - `GPG_PRIVATE_KEY` for release signing
  - `GPG_PASSPHRASE` for release signing
  - `CODECOV_TOKEN` for test coverage reporting
  - `AWS_ACCESS_KEY_ID` for AWS deployments
  - `AWS_SECRET_ACCESS_KEY` for AWS deployments
  - `AWS_S3_BUCKET` for UI deployments
  - `AWS_CLOUDFRONT_ID` for cache invalidation

### Workflow Fixes
1. CLI Workflow:
   - [ ] Set up proper test environment
   - [ ] Configure test database for integration tests
   - [ ] Add proper test coverage thresholds

2. Desktop Workflow:
   - [ ] Add proper test setup for Tauri
   - [ ] Configure code signing for macOS and Windows builds
   - [ ] Set up auto-update mechanism

3. Web Workflow:
   - [ ] Add E2E testing setup
   - [ ] Configure proper staging deployments
   - [ ] Set up preview deployments for PRs

### Security & Protection Rules
- [ ] Set up branch protection rules:
  - Require status checks to pass
  - Require PR reviews
  - Enforce signed commits

- [ ] Configure environment protection:
  - Production environment rules
  - Required reviewers
  - Deployment branch restrictions

### Documentation
- [ ] Add proper CI/CD documentation
- [ ] Document release process
- [ ] Add troubleshooting guide for common CI issues

## Next Steps
1. Focus on core functionality development
2. Set up basic testing infrastructure
3. Implement key features
4. Return to CI/CD improvements once core features are stable

## Notes
- Current workflows are basic but functional
- Security measures need to be properly implemented
- Test coverage needs improvement
- Release process needs refinement 