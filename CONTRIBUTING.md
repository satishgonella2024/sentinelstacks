# Contributing to SentinelStacks

Thank you for your interest in contributing to SentinelStacks! This document provides guidelines and instructions for contributing to the project.

## Code of Conduct

Please read and follow our [Code of Conduct](CODE_OF_CONDUCT.md) to foster an inclusive and respectful community.

## Getting Started

1. Fork the repository on GitHub
2. Clone your fork locally: `git clone https://github.com/your-username/sentinel.git`
3. Add the upstream repository: `git remote add upstream https://github.com/sentinelstacks/sentinel.git`
4. Create a new branch for your feature or bugfix: `git checkout -b feature/my-feature`
5. Make your changes, following our code style and commit guidelines
6. Push your branch to your fork: `git push origin feature/my-feature`
7. Open a Pull Request from your fork to the main repository

## Development Setup

Make sure you have the following installed:
- Go (version 1.21 or later)
- Git

Install dependencies:
```bash
go mod download
```

Build the project:
```bash
go build -o bin/sentinel ./cmd/sentinel
```

## Git Workflow

Please follow our [Git workflow strategy](docs/development/git_strategy.md) when making contributions. Key points:

- Use feature branches named according to our convention
- Follow our commit message format
- Keep PRs focused on a single change
- Rebase your branch on upstream main before submitting a PR

## Testing

All contributions should include appropriate tests. Run the tests to ensure your changes don't break existing functionality:

```bash
go test ./...
```

For more detailed testing instructions, see the [testing guide](docs/development/testing.md).

## Documentation

If your changes require updates to documentation, please include them in the same PR. Documentation is written in Markdown and stored in the `docs/` directory.

## Submitting a Pull Request

1. Make sure your code is properly tested
2. Ensure your commits follow our commit guidelines
3. Update documentation as necessary
4. Fill out the PR template completely
5. Request a review from a maintainer

## Getting Help

If you have questions or need help with the contribution process, you can:
- Open an issue with your question
- Ask in our community chat/forum
- Reach out to the maintainers directly

## License

By contributing to SentinelStacks, you agree that your contributions will be licensed under the project's [MIT License](LICENSE). 