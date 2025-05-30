# Contributing to HAWS

Thank you for your interest in contributing to HAWS (Hugo on AWS)! This document provides guidelines and instructions for contributing.

## Development Setup

1. **Prerequisites**
   - Go 1.22 or later
   - AWS account for testing (optional for code-only contributions)

2. **Clone the repository**
   ```
   git clone https://github.com/dragosboca/haws.git
   cd haws
   ```

3. **Install dependencies**
   ```
   go mod download
   ```

4. **Build the project**
   ```
   go build -o haws
   ```

## Code Style and Guidelines

- Follow standard Go style guidelines as outlined in [Effective Go](https://golang.org/doc/effective_go) and [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- Run `go fmt ./...` before committing changes
- Use meaningful commit messages that describe the change

## Testing

- Write tests for any new functionality
- Ensure all tests pass before submitting a pull request:
  ```
  ./run_tests.sh
  ```

## Logging

HAWS uses a structured logging system:

- `logger.Debug()`: For verbose debug information
- `logger.Info()`: For general operational information
- `logger.Warn()`: For warning events
- `logger.Error()`: For error events
- `logger.Fatal()`: For fatal events (will exit the application)

Please use appropriate logging levels in your code.

## Pull Request Process

1. Create a new branch for your feature or bugfix
2. Implement your changes
3. Add tests for new functionality
4. Update documentation as needed
5. Ensure all tests pass
6. Submit a pull request with a clear description of the changes

## Code Reviews

All submissions require review. We use GitHub pull requests for this purpose.

## License

By contributing to HAWS, you agree that your contributions will be licensed under the same license as the project.

Thank you for your contributions!
