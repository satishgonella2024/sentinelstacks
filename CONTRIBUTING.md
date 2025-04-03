# Contributing to Sentinel Stacks

Thank you for your interest in contributing to Sentinel Stacks! This document provides guidelines and instructions for contributing to the project.

## Git Workflow

### Branch Structure
```
main (production)
└── develop (integration)
    ├── feature/* (new features)
    ├── release/* (release preparation)
    └── hotfix/* (urgent fixes)
```

### Branch Naming Conventions
- Feature branches: `feature/description` (e.g., `feature/memory-plugins`)
- Release branches: `release/v1.0.0`
- Hotfix branches: `hotfix/description`

### Workflow Process
1. New features branch from `develop`
2. Feature branches merge back to `develop`
3. Releases branch from `develop`
4. Hotfixes branch from `main`
5. `develop` merges to `main` for releases

### Commit Message Convention
```
type(scope): description

Types:
- feat: new feature
- fix: bug fix
- docs: documentation changes
- style: formatting, missing semicolons, etc.
- refactor: code refactoring
- test: adding tests
- chore: maintenance
```

## Development Process

1. **Create a Feature Branch**
   ```bash
   git checkout develop
   git pull origin develop
   git checkout -b feature/your-feature-name
   ```

2. **Make Changes**
   - Write code
   - Add tests
   - Update documentation
   - Follow the code style guidelines

3. **Commit Changes**
   ```bash
   git add .
   git commit -m "feat(scope): description of changes"
   ```

4. **Push Changes**
   ```bash
   git push origin feature/your-feature-name
   ```

5. **Create Pull Request**
   - Create a PR from your feature branch to `develop`
   - Ensure all tests pass
   - Request review from team members

## Code Style Guidelines

1. **Go Code**
   - Follow the [Effective Go](https://golang.org/doc/effective_go.html) guidelines
   - Use `gofmt` for formatting
   - Run `go vet` before committing

2. **Documentation**
   - Document all exported types and functions
   - Keep comments clear and concise
   - Update README.md for significant changes

3. **Testing**
   - Write unit tests for new features
   - Ensure test coverage doesn't decrease
   - Use table-driven tests where appropriate

## Pull Request Process

1. **Before Submitting**
   - Ensure all tests pass
   - Update documentation
   - Rebase on latest develop branch
   - Squash commits if necessary

2. **PR Description**
   - Describe the changes
   - Reference related issues
   - List any breaking changes
   - Provide testing instructions

3. **Review Process**
   - Address reviewer comments
   - Update PR as needed
   - Get approval from at least one reviewer

## Getting Help

- Open an issue for questions or problems
- Join our community chat
- Check the documentation

## License

By contributing to Sentinel Stacks, you agree that your contributions will be licensed under the project's license. 