# Contributing to Slack Channel Cleaner

Thank you for your interest in contributing to Slack Channel Cleaner! This document provides guidelines and information for contributors.

## 🤝 How to Contribute

### Reporting Issues

Before creating bug reports, please check the existing issues to see if the problem has already been reported. When creating a bug report, please include:

- **Clear and descriptive title**
- **Detailed description** of the problem
- **Steps to reproduce** the issue
- **Expected behavior** vs **actual behavior**
- **Environment details** (OS, Go version, terminal type)
- **Screenshots** if applicable

### Feature Requests

We welcome feature requests! Please include:

- **Clear description** of the feature
- **Use case** and why it would be useful
- **Proposed implementation** (if you have ideas)

### Pull Requests

1. **Fork** the repository
2. **Create** a feature branch (`git checkout -b feature/amazing-feature`)
3. **Make** your changes
4. **Test** your changes thoroughly
5. **Commit** your changes (`git commit -m 'Add some amazing feature'`)
6. **Push** to the branch (`git push origin feature/amazing-feature`)
7. **Open** a Pull Request

## 🛠️ Development Setup

### Prerequisites

- Go 1.24 or higher
- Git
- A Slack workspace for testing

### Local Development

1. **Clone your fork**:
   ```bash
   git clone https://github.com/YOUR_USERNAME/workspace-channel-cleaner.git
cd workspace-channel-cleaner
   ```

2. **Set up environment**:
   ```bash
   cp example.env .env
   # Edit .env with your Slack token
   cp config/app.example.json config/app.json
   cp config/skiplist.example.json config/skiplist.json
   ```

3. **Install dependencies**:
   ```bash
   go mod tidy
   ```

4. **Build and test**:
   ```bash
   go build -o slack-cleaner-tui main.go
   ./slack-cleaner-tui
   ```

## 📋 Coding Guidelines

### Go Code Style

- Follow [Effective Go](https://golang.org/doc/effective_go.html)
- Use `gofmt` to format your code
- Keep functions small and focused
- Add comments for exported functions and complex logic
- Use meaningful variable and function names

### TUI Development

- **Responsive Design**: Ensure the UI works on different terminal sizes
- **User Experience**: Keep navigation intuitive and consistent
- **Error Handling**: Provide clear error messages to users
- **Performance**: Avoid blocking operations in the UI thread

### Testing

- Add tests for new functionality
- Ensure existing tests pass
- Test on different terminal types if possible
- Test with various Slack workspace configurations

## 🏗️ Project Structure

```
workspace-channel-cleaner/
├── main.go              # Application entry point
├── model/
│   └── model.go         # TUI model and state management
├── slack/
│   └── slack_client.go  # Workspace API integration
├── config/
│   ├── env.go          # Environment configuration
│   └── config.go       # Configuration management
├── ui/
│   └── actions.go      # UI action handlers
└── README.md           # Documentation
```

## 🔧 Key Components

### Model (`model/model.go`)
- **AppState**: Manages different UI screens
- **Styles**: Defines visual styling with Lip Gloss
- **Message Types**: Communication between components
- **State Handlers**: Keyboard input processing

### Slack Client (`slack/slack_client.go`)
- **Cleaner**: Main workspace API wrapper
- **Rate Limiting**: Automatic API rate limit handling
- **Channel Filtering**: Smart filtering logic
- **Skip List Management**: Load/save protected channels

### Configuration (`config/`)
- **Environment Loading**: Secure token management
- **Skip List**: JSON-based channel protection
- **Validation**: Input validation and error handling

## 🐛 Debugging

### Enable Debug Mode
```bash
DEBUG=1 go run main.go
```

### Common Issues
- **Rate Limiting**: The app handles this automatically, but you can monitor it
- **Token Issues**: Ensure your Slack token has the required scopes
- **UI Problems**: Test on different terminal emulators

## 📝 Documentation

- Update README.md for user-facing changes
- Add comments for complex code
- Update example files if configuration changes
- Document new features and options

## 🚀 Release Process

1. **Version Bump**: Update version in relevant files
2. **Changelog**: Document changes in CHANGELOG.md
3. **Testing**: Test on different platforms
4. **Tag**: Create a git tag for the release
5. **Publish**: Create a GitHub release

## 📞 Getting Help

- **Issues**: Use GitHub issues for bugs and feature requests
- **Discussions**: Use GitHub Discussions for questions and ideas
- **Code Review**: All PRs will be reviewed and feedback provided

## 🙏 Recognition

Contributors will be recognized in:
- The project README
- Release notes
- GitHub contributors list

Thank you for contributing to Slack Channel Cleaner! 🎉 