# Git Workflow Strategy

This document outlines the Git workflow strategy for the SentinelStacks project.

## Branching Strategy

We follow a trunk-based development approach with feature branches:

```
main (default branch)
├── feature/add-openai-integration
├── bugfix/fix-ollama-timeout
└── docs/update-provider-docs
```

### Branch Naming

- **Feature branches**: `feature/descriptive-name`
- **Bugfix branches**: `bugfix/issue-description`
- **Documentation branches**: `docs/topic-name`
- **Release branches**: `release/vX.Y.Z`

## Commit Guidelines

### Commit Message Format

```
<type>(<scope>): <subject>

<body>

<footer>
```

### Types

- **feat**: A new feature
- **fix**: A bug fix
- **docs**: Documentation changes
- **style**: Changes that don't affect code function (formatting, etc.)
- **refactor**: Code changes that neither fix bugs nor add features
- **test**: Adding or modifying tests
- **chore**: Changes to build process, tooling, etc.

### Example Commit Messages

```
feat(shim): add OpenAI integration

Implement OpenAI provider shim for GPT-4 and GPT-3.5 models.
Includes token counting and error handling.

Fixes #42
```

```
fix(cli): resolve panic when empty config file

The CLI would panic when encountering an empty config file.
Now it properly initializes default values.
```

## Pull Request Process

1. Create a branch from `main` using the naming convention above
2. Implement your changes with atomic commits following our commit guidelines
3. Ensure tests pass and add new ones as needed
4. Update documentation if required
5. Open a PR to `main` with a clear description of changes
6. Request reviews from at least one maintainer
7. Address review feedback
8. Squash and merge once approved

### PR Template

```markdown
## Description
Brief description of the changes

## Type of change
- [ ] Bug fix
- [ ] New feature
- [ ] Breaking change
- [ ] Documentation update

## How Has This Been Tested?
Describe the tests you ran

## Checklist
- [ ] My code follows the project's style guidelines
- [ ] I have added tests that prove my fix or feature works
- [ ] I have updated documentation as needed
- [ ] All tests pass locally and in CI
```

## Release Process

1. Create a release branch `release/vX.Y.Z` from `main`
2. Bump version numbers and update CHANGELOG.md
3. Create a PR from the release branch to `main`
4. After approval and merge, tag the release commit
5. Create a GitHub release with release notes

## Git Best Practices

- Keep branches short-lived (aim for < 1 week)
- Rebase feature branches on `main` frequently
- Write descriptive commit messages
- Squash commits when merging to maintain a clean history
- Never force push to `main`
- Use `git pull --rebase` to avoid merge commits

## GitHub Actions Integration

Our CI/CD pipeline automatically:

- Runs tests on every PR
- Builds and tests the project on pushes to `main`
- Deploys documentation on merges to `main`
- Builds and pushes Docker images for tagged releases 