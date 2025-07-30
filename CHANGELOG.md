# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Initial public release
- Beautiful TUI interface with Bubble Tea
- Interactive channel filtering and selection
- Configuration management system
- Skip list protection for important channels
- Pagination and navigation for large lists
- Dual view modes (table and simple list)
- Responsive design for different terminal sizes
- Rate limit handling for workspace API
- Error handling and user feedback

### Changed
- Converted from CLI to TUI application
- Moved hardcoded configurations to external files
- Improved user experience with interactive menus

### Fixed
- Header visibility issues in results screen
- Pagination navigation for large channel lists
- Color scheme for better visibility on white backgrounds

## [1.0.0] - 2024-12-19

### Added
- Initial release of Slack Channel Cleaner TUI
- Complete TUI interface with Bubble Tea framework
- Workspace API integration with rate limiting
- Configuration management with JSON files
- Skip list protection system
- Pagination and navigation features
- Dual view modes (table and simple list)
- Responsive design for terminal resizing
- Error handling and user feedback
- Comprehensive documentation

### Features
- **Main Menu**: Interactive menu system with keyboard navigation
- **Channel Search**: Find stale channels with configurable filters
- **Results Display**: Paginated results with selection capabilities
- **Configuration Editor**: In-app configuration management
- **Skip List Editor**: Add/remove protected channels
- **Confirmation System**: Double-confirmation for destructive actions
- **Navigation**: Multiple navigation options (arrows, page up/down, home/end)
- **View Toggle**: Switch between table and simple list views

### Technical
- Built with Go 1.24+
- Uses Bubble Tea for TUI framework
- Lip Gloss for styling
- Slack Go SDK for API integration
- Godotenv for environment management
- MIT License
- Comprehensive error handling
- Rate limit management
- Responsive design

---

## Version History

### v1.0.0 (Current)
- Initial public release
- Complete TUI functionality
- All core features implemented
- Production-ready application

---

## Migration Guide

### From CLI to TUI
This is the first public release, so there's no migration needed. The application has been completely rewritten from a simple CLI to a full-featured TUI application.

---

## Support

For support and questions:
- Create an issue on GitHub
- Check the README for documentation
- Review the troubleshooting section

---

*This changelog follows the [Keep a Changelog](https://keepachangelog.com/) format.* 