# Contributing to incident.io MCP Server

First off, thank you for considering contributing to this project! 

## Code of Conduct

This project and everyone participating in it is governed by our Code of Conduct. By participating, you are expected to uphold this code.

## How Can I Contribute?

### Reporting Bugs

Before creating bug reports, please check existing issues as you might find out that you don't need to create one. When you are creating a bug report, please include as many details as possible:

* Use a clear and descriptive title
* Describe the exact steps which reproduce the problem
* Provide specific examples to demonstrate the steps
* Describe the behavior you observed after following the steps
* Explain which behavior you expected to see instead and why
* Include logs and error messages

### Suggesting Enhancements

Enhancement suggestions are tracked as GitHub issues. When creating an enhancement suggestion, please include:

* Use a clear and descriptive title
* Provide a step-by-step description of the suggested enhancement
* Provide specific examples to demonstrate the steps
* Describe the current behavior and explain which behavior you expected to see instead
* Explain why this enhancement would be useful

### Pull Requests

* Fill in the required template
* Do not include issue numbers in the PR title
* Follow the Go style guide
* Include thoughtfully-worded, well-structured tests
* Document new code
* End all files with a newline

## Development Process

1. Fork the repo and create your branch from `main`
2. If you've added code that should be tested, add tests
3. If you've changed APIs, update the documentation
4. Ensure the test suite passes
5. Make sure your code lints
6. Issue that pull request!

### Local Development

```bash
# Clone your fork
git clone https://github.com/your-username/incidentio-mcp-golang.git
cd incidentio-mcp-golang

# Add upstream remote
git remote add upstream https://github.com/twentworth12/incidentio-mcp-golang.git

# Install dependencies
go mod download

# Run tests
make test

# Run linter
golangci-lint run

# Build
make build
```

### Testing

* Write unit tests for new functionality
* Ensure all tests pass: `make test`
* Add integration tests if applicable
* Test with real incident.io API when possible

### Code Style

* Follow standard Go conventions
* Use `gofmt` to format your code
* Follow the [Effective Go](https://go.dev/doc/effective_go) guidelines
* Keep functions small and focused
* Write descriptive variable and function names
* Add comments for exported functions and types

### Commit Messages

* Use the present tense ("Add feature" not "Added feature")
* Use the imperative mood ("Move cursor to..." not "Moves cursor to...")
* Limit the first line to 72 characters or less
* Reference issues and pull requests liberally after the first line

### Documentation

* Update the README.md with details of changes to the interface
* Update the TESTING.md if you add new test cases
* Add or update comments in the code
* Update tool descriptions if functionality changes

## Project Structure

```
.
├── cmd/              # Application entrypoints
├── internal/         # Private application code
│   ├── incidentio/   # incident.io API client
│   ├── server/       # MCP server implementation
│   └── tools/        # MCP tool implementations
├── pkg/              # Public libraries
└── test/             # Additional test files
```

## Adding New Tools

1. Create the API client method in `internal/incidentio/`
2. Add types to `internal/incidentio/types.go`
3. Create the tool wrapper in `internal/tools/`
4. Register the tool in `internal/server/server.go`
5. Add tests for both the API client and tool
6. Update the README with the new tool

## Questions?

Feel free to open an issue with your question or reach out to the maintainers.