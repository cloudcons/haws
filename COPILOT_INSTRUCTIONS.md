# Copilot Instructions for HAWS

This project follows best practices for Go development and CI/CD. Please adhere to the following guidelines when using GitHub Copilot or similar AI coding assistants:

## General Guidelines
- All code must be idiomatic Go, following effective Go and community conventions.
- Use explicit context.Context propagation for all operations that may block, perform I/O, or interact with AWS/cloud APIs.
- Do not embed context.Context in structs; always pass it as a parameter.
- All AWS SDK calls must be mockable for tests. Use interfaces and dependency injection where appropriate.
- All tests must be CI-friendly: do not make real AWS calls. Use mocks/fakes for all cloud interactions.
- Ensure all new code is covered by unit tests, including error and edge cases.
- Use Go modules for dependency management. Run `go mod tidy` before committing.
- All CLI entry points must propagate a single context.Context instance through the call chain.

## Release & Build
- Use GoReleaser for cross-platform builds and releases. See `.goreleaser.yml` for configuration.
- All binaries must be statically linked (CGO_ENABLED=0).

## Documentation
- Update README.md and code comments for all new features and changes.
- Document all exported functions and types.

## Pull Requests
- All PRs must pass CI, including tests and linting.
- Include a clear description of changes and reasoning.

---

For more details, see the project README and `.goreleaser.yml`.
