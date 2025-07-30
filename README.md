# 🔧 Slack Channel Cleaner

[![Go Version](https://img.shields.io/badge/Go-1.24+-blue.svg)](https://golang.org/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/ashfaqahmed/workspace-channels-cleaner)](https://goreportcard.com/report/github.com/ashfaqahmed/workspace-channels-cleaner)

> A terminal app I built to help clean up stale channels in Slack workspaces. It's got a nice UI, protects important channels, and makes it easy to manage your workspace without accidentally leaving channels you need.

> 💻 **Looking for a simple command-line CLI version?** Check out [workspace-channel-cleaner-cli](https://github.com/ashfaqahmed/workspace-channel-cleaner-cli) for a lightweight, automation-friendly CLI tool perfect for scripts and batch processing.

## ✨ Features

### 🎨 Nice Terminal Interface
- **Clean Design**: Simple, colorful interface that's easy to navigate
- **Menu System**: Keyboard-driven menus that feel natural
- **Live Updates**: See what's happening in real-time
- **Error Handling**: Clear error messages when something goes wrong

### 🔍 Channel Filtering
- **Find Old Channels**: Look for channels with no activity for N+ days (default: 30 days)
- **Keyword Filtering**: Filter channels by name keywords (can be set programmatically)
- **Channel Types**: Choose public, private, or both (default: public only)
- **Custom Limits**: Set how many channels to check at once (default: 30)

### 🛡️ Safety Features
- **Skip List**: Protect important channels from being left accidentally
- **Double Check**: Confirms before you leave any channels
- **Membership Check**: Only shows channels you're actually in
- **Rate Limits**: Handles Slack's API limits automatically

### 📝 Configuration
- **Skip List File**: Stores protected channels in `config/skiplist.json`
- **Environment Setup**: Keep your token secure in `.env` file
- **Good Defaults**: Sensible settings out of the box

### 📄 Navigation
- **Pagination**: Shows 12 items per page with easy navigation
- **Two Views**: Switch between table and simple list views
- **Quick Keys**: Page up/down, home/end, and arrow keys
- **Selection Memory**: Your selections stay when you move between pages

## 🚀 Quick Start

### 📸 Screenshots

**Main Menu**
![Main Menu](https://ashfaq.dev/github/workspace-channels-cleaner/main-menu.png)

**Configuration Screen**
![Configuration](https://ashfaq.dev/github/workspace-channels-cleaner/configuration.jpg)

**Channel Results with Pagination**
![Channel Results](https://ashfaq.dev/github/workspace-channels-cleaner/channel_leave.jpg)

**Skip List Management**
![Skip List](https://ashfaq.dev/github/workspace-channels-cleaner/skip_list.jpg)

### Prerequisites
- Go 1.24 or higher
- Workspace API token with required scopes

### Installation

1. **Clone the repository**:
   ```bash
   git clone https://github.com/ashfaqahmed/workspace-channels-cleaner
cd workspace-channels-cleaner
   ```

2. **Install dependencies**:
   ```bash
   go mod tidy
   ```

3. **Set up your workspace token**:
   ```bash
   cp example.env .env
   # Edit .env and add your WORKSPACE_API_TOKEN
   ```

4. **Set up configuration files**:
   ```bash
   cp config/app.example.json config/app.json
   cp config/skiplist.example.json config/skiplist.json
   ```

5. **Build and run**:
   ```bash
   go build -o workspace-cleaner-tui main.go
./workspace-cleaner-tui
   ```

## 🎮 Usage

### Main Menu Navigation
- **↑/↓ or k/j**: Navigate menu items
- **Enter**: Select menu item
- **q**: Quit application
- **Ctrl+C**: Force quit

### Available Options

#### 🔍 Find Stale Channels
- Shows current filter settings (days, limit, types, keyword)
- Search for stale channels based on configuration
- View results with last message timestamps
- Select channels to leave

#### ⚙️ Configuration
- View current configuration settings
- Edit configuration values (days, limit, types, verbose)
- Save configuration to `config/app.json`
- All changes are persistent across sessions

#### 📝 Edit Skip List
- View protected channels in skip list with pagination (10 items per page)
- Add new channels to the skip list
- Remove channels from the skip list
- Channels in this list are never processed

#### 🚪 Leave Channels
- Directly load and leave channels based on current settings
- Bypasses the search step

### Results Screen
- **↑/↓**: Navigate through channels
- **Space**: Select/deselect channel
- **Enter**: Leave selected channels
- **q**: Return to main menu
- **Page Up/Down (b/f)**: Jump 12 items up or down
- **Home/End (g/G)**: Go to first or last item
- **Toggle View (t)**: Switch between table and simple list view
- **Pagination**: Shows 12 items per page with page info
- **Responsive Table**: Automatically adjusts column widths based on terminal size
- **Smart Truncation**: Long channel names are truncated with "..." for better display

### Configuration Screen
- **e**: Enter edit mode
- **↑/↓**: Navigate between fields (in edit mode)
- **Enter**: Edit selected field (in edit mode)
- **s**: Save all changes (in edit mode)
- **q**: Cancel editing and return to view mode
- **Enter**: Return to main menu (in view mode)

### Confirmation Screen
- **y**: Confirm leaving selected channels
- **n**: Cancel and return to results
- **q**: Cancel and return to results

## ⚙️ Configuration

### Environment Variables
Create a `.env` file in the project root:

```env
# Required: Your workspace API token
WORKSPACE_API_TOKEN=xoxp-your-token-here

# Optional: Enable debug mode
DEBUG=1
```

### Application Configuration
The application configuration is stored in `config/app.json`:

```json
{
  "days": 30,
  "limit": 30,
  "types": ["public"],
  "verbose": false
}
```

**Configuration Editor Features:**
- **View Mode**: Press `e` to enter edit mode
- **Navigation**: Use ↑/↓ to select fields to edit
- **Field Editing**: Press Enter to edit a field, type new value, press Enter to save
- **Save Changes**: Press `s` to save all changes to the config file
- **Validation**: Automatic validation of input values

**Available Settings:**
- **Days**: Number of days of inactivity (minimum: 1)
- **Limit**: API request limit (minimum: 1)
- **Types**: Channel types to process (`public`, `private`, or both)
- **Verbose**: Enable detailed output (`true`/`false`)

### Skip List
The skip list is stored in `config/skiplist.json` and contains channels that should never be processed:

```json
[
  "general",
  "team-infra", 
  "support-team-requests",
  "company-announcements"
]
```

**Skip List Editor Features:**
- **Navigation**: Use ↑/↓ to browse through channels
- **Pagination**: Automatically paginates long lists (10 items per page)
- **Add Channels**: Press 'a' to add new channels to the skip list
- **Remove Channels**: Press 'd' to delete channels from the skip list
- **Save Changes**: All changes are automatically saved to the JSON file

## 🔧 Workspace API Requirements

Your workspace token must have the following scopes:
- `channels:history` - Read channel message history
- `groups:history` - Read private channel history  
- `conversations.list` - List all channels
- `conversations.leave` - Leave channels

### Getting Your Workspace Token

1. Go to [Slack API Apps](https://api.slack.com/apps)
2. Create a new app or select an existing one
3. Go to "OAuth & Permissions"
4. Add the required scopes listed above
5. Install the app to your workspace
6. Copy the "Bot User OAuth Token" (starts with `xoxb-`) or "User OAuth Token" (starts with `xoxp-`)

## 🎨 UI Features

### Color Scheme
- **Dark Purple**: Primary theme color for titles and selection (optimized for white backgrounds)
- **Green**: Success states and selected items
- **Red**: Error messages
- **Yellow**: Warning messages
- **Blue**: Information text
- **Dark Gray**: Subtitles and secondary text (high contrast)
- **Black**: Menu items and main text (excellent readability)

### Responsive Design
- **Full-Width Layout**: Automatically expands to fill the entire terminal width
- **Dynamic Sizing**: Adapts to terminal window resizing in real-time
- **Smart Column Layout**: Results table adjusts column widths based on available space
- **Centered Content**: Menu items and titles are automatically centered
- **Minimum Width**: Maintains readability with minimum width constraints
- **Text Truncation**: Long channel names are intelligently truncated with ellipsis

## 🛠️ Development

### Project Structure
```
workspace-channels-cleaner/
├── main.go              # Application entry point
├── model/
│   └── model.go         # TUI model and state management
├── slack/
│   └── slack_client.go  # Slack API integration
├── config/
│   ├── env.go          # Environment configuration
│   ├── config.go       # Configuration management
│   ├── app.example.json # Example configuration
│   └── skiplist.example.json # Example skip list

├── example.env         # Example environment file
├── .gitignore          # Git ignore rules
└── README.md           # This file
```

### Building from Source

```bash
# Clone the repository
git clone https://github.com/ashfaqahmed/workspace-channels-cleaner
cd workspace-channels-cleaner

# Install dependencies
go mod tidy

# Build the application
go build -o workspace-cleaner-tui main.go

# Run the application
./workspace-cleaner-tui
```

### Dependencies
- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - TUI framework
- [Lip Gloss](https://github.com/charmbracelet/lipgloss) - Styling library
- [Slack Go SDK](https://github.com/slack-go/slack) - Workspace API client
- [Godotenv](https://github.com/joho/godotenv) - Environment variable loading

## 🔒 Security

- **Token Protection**: Never logged or displayed
- **Skip List**: Protects important channels
- **Confirmation**: Double-confirmation for destructive actions
- **Error Handling**: Graceful failure without data loss

## 🐛 Troubleshooting

### Common Issues

**"WORKSPACE_API_TOKEN not set"**
- Ensure your `.env` file exists and contains the token
- Check that the token has the required scopes

**"Rate limit hit"**
- The application automatically handles rate limits
- Wait for the retry mechanism to complete

**"No channels found"**
- Check your filter settings
- Verify you're a member of channels
- Review the skip list for protected channels

**"Cannot navigate results"**
- Use arrow keys for single item navigation
- Use `b`/`f` for page navigation
- Use `g`/`G` for home/end navigation
- Press `t` to toggle between view modes

### Debug Mode
Set the `DEBUG` environment variable for additional logging:
```bash
DEBUG=1 go run main.go
```

## 🤝 Contributing

Feel free to contribute! If you want to add something big, open an issue first so we can talk about it.

1. Fork the repo
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### Development Guidelines
- Follow Go coding standards
- Add tests for new features
- Update docs when needed
- Keep the UI responsive and user-friendly

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## ⚖️ Legal Disclaimers

**Important**: Please read our [Legal Disclaimer](DISCLAIMER.md) for important information about third-party trademarks, copyrights, and terms of use.

### Copyright Disclaimer

This project is **NOT** affiliated with, endorsed by, or sponsored by Slack Technologies, Inc. or any of its subsidiaries.

### Third-Party Trademarks and Copyrights
- **Slack** is a registered trademark of Slack Technologies, Inc.
- **Slack API** and related services are owned by Slack Technologies, Inc.
- This project uses the official [Slack Go SDK](https://github.com/slack-go/slack) which is subject to its own license terms.

### Fair Use
This project is developed for educational and productivity purposes, using Slack's publicly available API in accordance with their [API Terms of Service](https://slack.com/terms-of-service/api). The use of Slack's API and SDK is subject to Slack's own terms and conditions.

### No Warranty
This project is provided "as is" without any warranties. Users are responsible for ensuring their use complies with Slack's terms of service and applicable laws.

## 🙏 Acknowledgments

- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - TUI framework
- [Lip Gloss](https://github.com/charmbracelet/lipgloss) - Styling library
- [Slack Go SDK](https://github.com/slack-go/slack) - Slack API client
- The open-source community for inspiration and tools

---

## ☕ Support My Work

If this tool saved you time or effort, consider buying me a coffee.
Your support helps me keep building and maintaining open-source projects like this!

You can either scan the QR code below or click the link to tip me:

👉 [buymeacoffee.com/ashfaqueali](https://buymeacoffee.com/ashfaqueali)

**Buy Me a Coffee QR**

<img src="https://ashfaqsolangi.com/images/bmc_qr.png" alt="Buy Me a Coffee QR" width="220" height="220" />

---

**Happy channel cleaning! 🧹✨**

*Made with ❤️ by [Ashfaque Ali](https://github.com/ashfaqahmed)* 