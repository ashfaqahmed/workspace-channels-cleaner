package model

import (
	"fmt"
	"strconv"
	"strings"

	"workspace-channel-cleaner/config"
	"workspace-channel-cleaner/slack"

	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
)

// AppState represents the current state of the application
type AppState int

const (
	MainMenu AppState = iota
	ConfigScreen
	FilterScreen
	ResultsScreen
	ConfirmationScreen
	SkipListScreen
	LoadingScreen
)

// model represents the main application model
type model struct {
	state AppState
	
	// Terminal dimensions
	width  int
	height int
	
	// Main menu
	cursor int
	choices []string
	
	// Configuration
	config   *config.AppConfig
	
	// Results
	channels []slack.ChannelInfo
	selected map[int]struct{}
	resultsOffset int // For pagination in results screen
	useSimpleView bool // Toggle between table and simple list view
	
	// Skip list
	skipList map[string]bool
	skipCursor int
	skipChoices []string
	skipListOffset int // For pagination
	skipListMode string // "view", "add", "remove"
	skipListInput string // For adding new channels
	
	// Config editing
	configMode string // "view", "edit"
	configCursor int
	configInput string
	editingField string // "days", "limit", "types", "verbose", "keyword"
	
	// Loading
	loadingMsg string
	
	// Error handling
	err error
	
	// Styles
	styles *Styles
}

// Styles defines the visual styling for the TUI
type Styles struct {
	title     lipgloss.Style
	subtitle  lipgloss.Style
	menu      lipgloss.Style
	cursor    lipgloss.Style
	selected  lipgloss.Style
	error     lipgloss.Style
	success   lipgloss.Style
	warning   lipgloss.Style
	info      lipgloss.Style
	border    lipgloss.Style
}

// NewStyles creates a new style configuration
func NewStyles() *Styles {
	return &Styles{
		title: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#3D0066")). // dark purple
			Background(lipgloss.Color("#E0D7F7")). // light purple background
			Padding(0, 1).
			Bold(true),
		
		subtitle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#555555")). // medium gray
			Italic(true),
		
		menu: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#222222")). // black
			Padding(0, 1),
		
		cursor: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#7D56F4")). // purple
			Bold(true),
		
		selected: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#04B575")). // green
			Bold(true),
		
		error: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF5F56")). // red
			Bold(true),
		
		success: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#04B575")). // green
			Bold(true),
		
		warning: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFBD2E")). // yellow
			Bold(true),
		
		info: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#5A9FD4")), // blue
		
		border: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#874BFD")).
			Padding(1, 2),
	}
}

// getResponsiveBorder returns a border style that adapts to terminal width
func (m model) getResponsiveBorder() lipgloss.Style {
	// Use full width minus 2 for padding
	width := m.width - 2
	if width < 60 {
		width = 60 // Minimum width
	}
	
	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#874BFD")).
		Padding(1, 2).
		Width(width)
}

// InitialModel creates the initial application model
func InitialModel() model {
	// Load configuration
	appConfig, err := config.LoadConfig(config.GetConfigPath())
	if err != nil {
		// Use default config if loading fails
		appConfig = config.DefaultConfig()
	}
	
	return model{
		state: MainMenu,
		width:  80,  // Default width
		height: 24,  // Default height
		choices: []string{
			"🔍 Find Stale Channels",
			"⚙️  Configuration",
			"📝 Edit Skip List",
			"🚪 Leave Channels",
			"❌ Exit",
		},
		selected: make(map[int]struct{}),
		styles:   NewStyles(),
		config:   appConfig,
		configMode: "view",
		configCursor: 0,
	}
}

// Init initializes the model
func (m model) Init() tea.Cmd {
	return nil
}


func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m.handleKeyPress(msg)
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	case errorMsg:
		m.err = msg
		return m, nil
	case channelsLoadedMsg:
		m.channels = msg.channels
		m.resultsOffset = 0 // Reset pagination
		m.cursor = 0 // Reset cursor
		m.state = ResultsScreen
		return m, nil
	case channelsLeftMsg:
		m.state = MainMenu
		return m, nil
	}
	return m, nil
}


func (m model) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch m.state {
	case MainMenu:
		return m.handleMainMenu(msg)
	case ConfigScreen:
		return m.handleConfigScreen(msg)
	case FilterScreen:
		return m.handleFilterScreen(msg)
	case ResultsScreen:
		return m.handleResultsScreen(msg)
	case ConfirmationScreen:
		return m.handleConfirmationScreen(msg)
	case SkipListScreen:
		return m.handleSkipListScreen(msg)
	}
	return m, nil
}


func (m model) handleMainMenu(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		return m, tea.Quit
	case "up", "k":
		if m.cursor > 0 {
			m.cursor--
		}
	case "down", "j":
		if m.cursor < len(m.choices)-1 {
			m.cursor++
		}
	case "enter":
		return m.handleMenuSelection()
	}
	return m, nil
}


func (m model) handleMenuSelection() (tea.Model, tea.Cmd) {
	switch m.cursor {
	case 0: // Find Stale Channels
		m.state = FilterScreen
		return m, nil
	case 1: // Configuration
		m.state = ConfigScreen
		return m, nil
	case 2: // Edit Skip List
		return m.loadSkipList()
	case 3: // Leave Channels
		return m.loadChannels()
	case 4: // Exit
		return m, tea.Quit
	}
	return m, nil
}


func (m model) handleConfigScreen(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// If editing a specific field, handle that first
	if m.editingField != "" {
		return m.handleConfigFieldEdit(msg)
	}
	
	switch m.configMode {
	case "view":
		return m.handleConfigView(msg)
	case "edit":
		return m.handleConfigEdit(msg)
	}
	return m, nil
}

func (m model) handleConfigView(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		m.state = MainMenu
		return m, nil
	case "e":
		m.configMode = "edit"
		m.configCursor = 0
		return m, nil
	case "enter":
		m.state = MainMenu
		return m, nil
	}
	return m, nil
}

func (m model) handleConfigEdit(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		m.configMode = "view"
		m.configInput = ""
		m.editingField = ""
		return m, nil
	case "up", "k":
		if m.configCursor > 0 {
			m.configCursor--
		}
	case "down", "j":
		if m.configCursor < 4 { // 5 fields: days, limit, types, verbose, keyword
			m.configCursor++
		}
	case "enter":
		return m.startConfigFieldEdit()
	case "s":
		// Save configuration
		err := config.SaveConfig(config.GetConfigPath(), m.config)
		if err != nil {
			m.err = err
		}
		m.configMode = "view"
		return m, nil
	}
	return m, nil
}

func (m model) startConfigFieldEdit() (tea.Model, tea.Cmd) {
	switch m.configCursor {
	case 0: // Days
		m.editingField = "days"
		m.configInput = fmt.Sprintf("%d", m.config.Days)
	case 1: // Limit
		m.editingField = "limit"
		m.configInput = fmt.Sprintf("%d", m.config.Limit)
	case 2: // Types
		m.editingField = "types"
		m.configInput = strings.Join(m.config.Types, ",")
	case 3: // Verbose
		m.editingField = "verbose"
		m.configInput = fmt.Sprintf("%t", m.config.Verbose)
	case 4: // Keyword
		m.editingField = "keyword"
		m.configInput = m.config.Keyword
	}
	return m, nil
}

func (m model) handleConfigFieldEdit(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		m.editingField = ""
		m.configInput = ""
		return m, nil
	case "enter":
		return m.saveConfigField()
	case "backspace":
		if len(m.configInput) > 0 {
			m.configInput = m.configInput[:len(m.configInput)-1]
		}
	default:
		if len(msg.String()) == 1 {
			m.configInput += msg.String()
		}
	}
	return m, nil
}

func (m model) saveConfigField() (tea.Model, tea.Cmd) {
	switch m.editingField {
	case "days":
		if days, err := strconv.Atoi(m.configInput); err == nil && days > 0 {
			m.config.Days = days
		}
	case "limit":
		if limit, err := strconv.Atoi(m.configInput); err == nil && limit > 0 {
			m.config.Limit = limit
		}
	case "types":
		types := strings.Split(m.configInput, ",")
		validTypes := make([]string, 0)
		for _, t := range types {
			t = strings.TrimSpace(t)
			if t == "public" || t == "private" {
				validTypes = append(validTypes, t)
			}
		}
		if len(validTypes) > 0 {
			m.config.Types = validTypes
		}
	case "verbose":
		if m.configInput == "true" {
			m.config.Verbose = true
		} else if m.configInput == "false" {
			m.config.Verbose = false
		}
	case "keyword":
		m.config.Keyword = m.configInput
	}
	
	m.editingField = ""
	m.configInput = ""
	return m, nil
}

func (m model) handleFilterScreen(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		m.state = MainMenu
		return m, nil
	case "enter":
		return m.startChannelSearch()
	}
	return m, nil
}

func (m model) handleResultsScreen(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		m.state = MainMenu
		return m, nil
	case "up", "k":
		if m.cursor > 0 {
			m.cursor--
		}
		// Adjust offset for pagination
		if m.cursor < m.resultsOffset {
			m.resultsOffset = m.cursor
		}
	case "down", "j":
		if m.cursor < len(m.channels)-1 {
			m.cursor++
		}
		// Adjust offset for pagination (show 12 items per page)
		if m.cursor >= m.resultsOffset+12 {
			m.resultsOffset = m.cursor - 11
		}
	case " ":
		if _, ok := m.selected[m.cursor]; ok {
			delete(m.selected, m.cursor)
		} else {
			m.selected[m.cursor] = struct{}{}
		}
	case "enter":
		if len(m.selected) > 0 {
			m.state = ConfirmationScreen
		}
	case "pageup", "b":
		// Page up (move cursor and offset up by 12)
		if m.cursor >= 12 {
			m.cursor -= 12
			if m.resultsOffset >= 12 {
				m.resultsOffset -= 12
			} else {
				m.resultsOffset = 0
			}
		} else {
			m.cursor = 0
			m.resultsOffset = 0
		}
	case "pagedown", "f":
		// Page down (move cursor and offset down by 12)
		if m.cursor+12 < len(m.channels) {
			m.cursor += 12
			m.resultsOffset += 12
		} else {
			m.cursor = len(m.channels) - 1
			// Adjust offset to show the last page
			if len(m.channels) > 12 {
				m.resultsOffset = len(m.channels) - 12
			}
		}
	case "home", "g":
		// Go to first item
		m.cursor = 0
		m.resultsOffset = 0
	case "end", "G":
		// Go to last item
		m.cursor = len(m.channels) - 1
		if len(m.channels) > 12 {
			m.resultsOffset = len(m.channels) - 12
		}
	case "t":
		// Toggle between table and simple view
		m.useSimpleView = !m.useSimpleView
	}
	return m, nil
}

func (m model) handleConfirmationScreen(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		m.state = ResultsScreen
		return m, nil
	case "y", "Y":
		return m.leaveSelectedChannels()
	case "n", "N":
		m.state = ResultsScreen
		return m, nil
	}
	return m, nil
}

func (m model) handleSkipListScreen(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch m.skipListMode {
	case "view":
		return m.handleSkipListView(msg)
	case "add":
		return m.handleSkipListAdd(msg)
	case "remove":
		return m.handleSkipListRemove(msg)
	}
	return m, nil
}

func (m model) handleSkipListView(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		m.state = MainMenu
		return m, nil
	case "up", "k":
		if m.skipCursor > 0 {
			m.skipCursor--
		}
		// Adjust offset for pagination
		if m.skipCursor < m.skipListOffset {
			m.skipListOffset = m.skipCursor
		}
	case "down", "j":
		if m.skipCursor < len(m.skipChoices)-1 {
			m.skipCursor++
		}
		// Adjust offset for pagination
		if m.skipCursor >= m.skipListOffset+10 {
			m.skipListOffset = m.skipCursor - 9
		}
	case "a":
		m.skipListMode = "add"
		return m, nil
	case "d":
		if len(m.skipChoices) > 0 {
			m.skipListMode = "remove"
		}
		return m, nil
	case "enter":
		m.state = MainMenu
		return m, nil
	}
	return m, nil
}

func (m model) handleSkipListAdd(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		m.skipListMode = "view"
		m.skipListInput = ""
		return m, nil
	case "enter":
		if m.skipListInput != "" {
			// Add to skip list
			m.skipList[m.skipListInput] = true
			m.skipChoices = append(m.skipChoices, m.skipListInput)
			// Save to file
			err := slack.SaveSkipList("config/skiplist.json", m.skipList)
			if err != nil {
				m.err = err
			}
		}
		m.skipListMode = "view"
		m.skipListInput = ""
		return m, nil
	case "backspace":
		if len(m.skipListInput) > 0 {
			m.skipListInput = m.skipListInput[:len(m.skipListInput)-1]
		}
	default:
		if len(msg.String()) == 1 {
			m.skipListInput += msg.String()
		}
	}
	return m, nil
}

func (m model) handleSkipListRemove(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		m.skipListMode = "view"
		return m, nil
	case "up", "k":
		if m.skipCursor > 0 {
			m.skipCursor--
		}
	case "down", "j":
		if m.skipCursor < len(m.skipChoices)-1 {
			m.skipCursor++
		}
	case "enter":
		if len(m.skipChoices) > 0 && m.skipCursor < len(m.skipChoices) {
			// Remove from skip list
			channelToRemove := m.skipChoices[m.skipCursor]
			delete(m.skipList, channelToRemove)
			
			// Remove from choices slice
			m.skipChoices = append(m.skipChoices[:m.skipCursor], m.skipChoices[m.skipCursor+1:]...)
			
			// Adjust cursor
			if m.skipCursor >= len(m.skipChoices) && len(m.skipChoices) > 0 {
				m.skipCursor = len(m.skipChoices) - 1
			}
			
			// Save to file
			err := slack.SaveSkipList("config/skiplist.json", m.skipList)
			if err != nil {
				m.err = err
			}
		}
		m.skipListMode = "view"
		return m, nil
	}
	return m, nil
}

// loadSkipList loads the skip list for editing
func (m model) loadSkipList() (tea.Model, tea.Cmd) {
	skipList, err := slack.LoadSkipList("config/skiplist.json")
	if err != nil {
		m.err = err
		return m, nil
	}
	
	m.skipList = skipList
	m.skipChoices = make([]string, 0, len(skipList))
	for name := range skipList {
		m.skipChoices = append(m.skipChoices, name)
	}
	m.skipCursor = 0
	m.skipListOffset = 0
	m.skipListMode = "view"
	m.skipListInput = ""
	m.state = SkipListScreen
	return m, nil
}

// loadChannels loads channels for leaving
func (m model) loadChannels() (tea.Model, tea.Cmd) {
	m.state = LoadingScreen
	m.loadingMsg = "Loading channels..."
	
	return m, func() tea.Msg {
		token := config.GetWorkspaceToken()
		cleaner := slack.NewCleaner(token, m.config.Limit, slack.GetChannelTypes(m.config.Types), m.config.Days, m.config.Keyword, m.config.Verbose)
		channels, err := cleaner.GetFilteredChannels()
		if err != nil {
			return errorMsg{err}
		}
		return channelsLoadedMsg{channels}
	}
}

func (m model) startChannelSearch() (tea.Model, tea.Cmd) {
	m.state = LoadingScreen
	m.loadingMsg = "Searching for stale channels..."
	
	return m, func() tea.Msg {
		token := config.GetWorkspaceToken()
		cleaner := slack.NewCleaner(token, m.config.Limit, slack.GetChannelTypes(m.config.Types), m.config.Days, m.config.Keyword, m.config.Verbose)
		channels, err := cleaner.GetFilteredChannels()
		if err != nil {
			return errorMsg{err}
		}
		return channelsLoadedMsg{channels}
	}
}

// leaveSelectedChannels leaves the selected channels
func (m model) leaveSelectedChannels() (tea.Model, tea.Cmd) {
	m.state = LoadingScreen
	m.loadingMsg = "Leaving selected channels..."
	
	selectedChannels := make([]slack.ChannelInfo, 0)
	for i := range m.selected {
		if i < len(m.channels) {
			selectedChannels = append(selectedChannels, m.channels[i])
		}
	}
	
	return m, func() tea.Msg {
		token := config.GetWorkspaceToken()
		cleaner := slack.NewCleaner(token, m.config.Limit, slack.GetChannelTypes(m.config.Types), m.config.Days, m.config.Keyword, m.config.Verbose)
		err := cleaner.LeaveChannels(selectedChannels)
		if err != nil {
			return errorMsg{err}
		}
		return channelsLeftMsg{}
	}
}

func (m model) View() string {
	switch m.state {
	case MainMenu:
		return m.renderMainMenu()
	case ConfigScreen:
		return m.renderConfigScreen()
	case FilterScreen:
		return m.renderFilterScreen()
	case ResultsScreen:
		return m.renderResultsScreen()
	case ConfirmationScreen:
		return m.renderConfirmationScreen()
	case SkipListScreen:
		return m.renderSkipListScreen()
	case LoadingScreen:
		return m.renderLoadingScreen()
	}
	return ""
}

func (m model) renderMainMenu() string {
	var b strings.Builder
	
	// Center the title
	title := m.styles.title.Render("🔧 Workspace Channel Cleaner")
	titleWidth := lipgloss.Width(title)
	titlePadding := (m.width - titleWidth - 4) / 2 // -4 for border padding
	if titlePadding < 0 {
		titlePadding = 0
	}
	
	b.WriteString(strings.Repeat(" ", titlePadding))
	b.WriteString(title)
	b.WriteString("\n\n")
	
	// Center the menu items
	menuWidth := 0
	for _, choice := range m.choices {
		if lipgloss.Width(choice) > menuWidth {
			menuWidth = lipgloss.Width(choice)
		}
	}
	menuWidth += 2 // For cursor space
	
	menuPadding := (m.width - menuWidth - 4) / 2
	if menuPadding < 0 {
		menuPadding = 0
	}
	
	for i, choice := range m.choices {
		cursor := " "
		if m.cursor == i {
			cursor = m.styles.cursor.Render(">")
		}
		
		style := m.styles.menu
		if m.cursor == i {
			style = m.styles.selected
		}
		
		b.WriteString(strings.Repeat(" ", menuPadding))
		b.WriteString(fmt.Sprintf("%s %s\n", cursor, style.Render(choice)))
	}
	
	b.WriteString("\n")
	instructions := m.styles.subtitle.Render("Use ↑↓ to navigate, Enter to select, q to quit")
	instPadding := (m.width - lipgloss.Width(instructions) - 4) / 2
	if instPadding < 0 {
		instPadding = 0
	}
	b.WriteString(strings.Repeat(" ", instPadding))
	b.WriteString(instructions)
	
	if m.err != nil {
		b.WriteString("\n\n")
		errorMsg := m.styles.error.Render("Error: " + m.err.Error())
		errPadding := (m.width - lipgloss.Width(errorMsg) - 4) / 2
		if errPadding < 0 {
			errPadding = 0
		}
		b.WriteString(strings.Repeat(" ", errPadding))
		b.WriteString(errorMsg)
	}
	
	return m.getResponsiveBorder().Render(b.String())
}

func (m model) renderConfigScreen() string {
	// If editing a specific field, show the edit screen
	if m.editingField != "" {
		return m.renderConfigFieldEdit()
	}
	
	switch m.configMode {
	case "view":
		return m.renderConfigView()
	case "edit":
		return m.renderConfigEdit()
	}
	return m.getResponsiveBorder().Render("")
}

func (m model) renderConfigView() string {
	var b strings.Builder
	
	b.WriteString(m.styles.title.Render("⚙️  Configuration"))
	b.WriteString("\n\n")
	
	b.WriteString(fmt.Sprintf("Days: %d\n", m.config.Days))
	b.WriteString(fmt.Sprintf("Limit: %d\n", m.config.Limit))
	b.WriteString(fmt.Sprintf("Types: %s\n", strings.Join(m.config.Types, ", ")))
	b.WriteString(fmt.Sprintf("Verbose: %t\n", m.config.Verbose))
	b.WriteString(fmt.Sprintf("Keyword: %s\n", m.config.Keyword))
	
	b.WriteString("\n")
	b.WriteString(m.styles.subtitle.Render("Press 'e' to edit, Enter to return to main menu"))
	
	return m.getResponsiveBorder().Render(b.String())
}

func (m model) renderConfigEdit() string {
	var b strings.Builder
	
	b.WriteString(m.styles.title.Render("⚙️  Edit Configuration"))
	b.WriteString("\n\n")
	
	fields := []string{
		fmt.Sprintf("Days: %d", m.config.Days),
		fmt.Sprintf("Limit: %d", m.config.Limit),
		fmt.Sprintf("Types: %s", strings.Join(m.config.Types, ",")),
		fmt.Sprintf("Verbose: %t", m.config.Verbose),
		fmt.Sprintf("Keyword: %s", m.config.Keyword),
	}
	
	for i, field := range fields {
		cursor := " "
		if m.configCursor == i {
			cursor = m.styles.cursor.Render(">")
		}
		b.WriteString(fmt.Sprintf("%s %s\n", cursor, field))
	}
	
	b.WriteString("\n")
	b.WriteString(m.styles.subtitle.Render("Use ↑↓ to navigate, Enter to edit, 's' to save, q to cancel"))
	
	return m.getResponsiveBorder().Render(b.String())
}

func (m model) renderConfigFieldEdit() string {
	var b strings.Builder
	
	b.WriteString(m.styles.title.Render(fmt.Sprintf("⚙️  Edit %s", strings.Title(m.editingField))))
	b.WriteString("\n\n")
	
	b.WriteString(fmt.Sprintf("Current value: %s\n", m.configInput))
	b.WriteString("New value: ")
	
	// Show input with cursor
	if m.configInput == "" {
		b.WriteString(m.styles.cursor.Render("_"))
	} else {
		b.WriteString(m.configInput)
		b.WriteString(m.styles.cursor.Render("_"))
	}
	
	b.WriteString("\n\n")
	b.WriteString(m.styles.subtitle.Render("Type new value and press Enter to save, q to cancel"))
	
	return m.getResponsiveBorder().Render(b.String())
}

func (m model) renderFilterScreen() string {
	var b strings.Builder
	
	b.WriteString(m.styles.title.Render("🔍 Filter Configuration"))
	b.WriteString("\n\n")
	
	b.WriteString(fmt.Sprintf("Days: %d\n", m.config.Days))
	b.WriteString(fmt.Sprintf("Keyword: %s\n", m.config.Keyword))
	b.WriteString(fmt.Sprintf("Limit: %d\n", m.config.Limit))
	b.WriteString(fmt.Sprintf("Types: %s\n", strings.Join(m.config.Types, ", ")))
	
	b.WriteString("\n")
	b.WriteString(m.styles.subtitle.Render("Press Enter to start search"))
	
	return m.getResponsiveBorder().Render(b.String())
}

func (m model) renderResultsScreen() string {
	var b strings.Builder
	
	b.WriteString(m.styles.title.Render("📋 Channel Results"))
	b.WriteString("\n\n")
	
	if len(m.channels) == 0 {
		b.WriteString(m.styles.info.Render("No channels found matching your criteria."))
		b.WriteString("\n\n")
		b.WriteString(m.styles.subtitle.Render("Press q to return to main menu"))
		return m.getResponsiveBorder().Render(b.String())
	}
	
	b.WriteString(fmt.Sprintf("Found %d channel(s):\n\n", len(m.channels)))
	
	// Show paginated results (12 items per page to ensure headers are visible)
	start := m.resultsOffset
	end := start + 12
	if end > len(m.channels) {
		end = len(m.channels)
	}
	
	visibleChannels := m.channels[start:end]
	
	// Calculate column widths based on terminal width
	availableWidth := m.width - 10 // Account for border and padding
	selectColWidth := 8  // Fixed width for selection column
	nameColWidth := (availableWidth - selectColWidth) * 3 / 5 // 60% for name
	dateColWidth := (availableWidth - selectColWidth) * 2 / 5  // 40% for date
	
	if nameColWidth < 15 {
		nameColWidth = 15
	}
	if dateColWidth < 20 {
		dateColWidth = 20
	}
	
	// Create table rows for visible channels only
	var rows [][]string
	for i, ch := range visibleChannels {
		globalIndex := start + i
		cursor := " "
		if m.cursor == globalIndex {
			cursor = m.styles.cursor.Render(">")
		}
		
		checked := " "
		if _, ok := m.selected[globalIndex]; ok {
			checked = m.styles.selected.Render("✓")
		}
		
		lastSeen := "No messages"
		if !ch.LastSeen.IsZero() {
			lastSeen = ch.LastSeen.Format("2006-01-02 15:04:05")
		}
		
		// Truncate name if too long
		name := ch.Name
		if len(name) > nameColWidth-3 {
			name = name[:nameColWidth-6] + "..."
		}
		
		rows = append(rows, []string{
			fmt.Sprintf("%s [%s]", cursor, checked),
			fmt.Sprintf("#%s", name),
			lastSeen,
		})
	}
	
	// Create table with better styling
	t := table.New().
		Border(lipgloss.RoundedBorder()).
		BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("#874BFD"))).
		Headers("Select", "Channel", "Last Activity").
		Rows(rows...).
		StyleFunc(func(row, col int) lipgloss.Style {
			switch col {
			case 0:
				return lipgloss.NewStyle().Width(selectColWidth).Align(lipgloss.Center)
			case 1:
				return lipgloss.NewStyle().Width(nameColWidth)
			case 2:
				return lipgloss.NewStyle().Width(dateColWidth)
			default:
				return lipgloss.NewStyle()
			}
		})
	
	if m.useSimpleView {
		b.WriteString(m.renderSimpleListView(visibleChannels, start))
	} else {
		b.WriteString(t.Render())
	}
	
	// Show pagination info
	if len(m.channels) > 12 {
		currentPage := (m.resultsOffset / 12) + 1
		totalPages := (len(m.channels) + 11) / 12
		b.WriteString(fmt.Sprintf("\n%s\n", m.styles.info.Render(fmt.Sprintf("Page %d/%d (Showing %d-%d of %d)", currentPage, totalPages, start+1, end, len(m.channels)))))
	}
	
	b.WriteString("\n")
	b.WriteString(m.styles.subtitle.Render("Use ↑↓ to navigate, Space to select, Enter to leave selected"))
	b.WriteString("\n")
	b.WriteString(m.styles.subtitle.Render("Page Up/Down (b/f), Home/End (g/G), 't' to toggle view, q to quit"))
	
	return m.getResponsiveBorder().Render(b.String())
}

func (m model) renderConfirmationScreen() string {
	var b strings.Builder
	
	b.WriteString(m.styles.title.Render("⚠️  Confirmation"))
	b.WriteString("\n\n")
	
	selectedCount := len(m.selected)
	b.WriteString(fmt.Sprintf("Are you sure you want to leave %d channel(s)?\n\n", selectedCount))
	
	selectedChannels := make([]string, 0)
	for i := range m.selected {
		if i < len(m.channels) {
			selectedChannels = append(selectedChannels, "#"+m.channels[i].Name)
		}
	}
	
	for _, name := range selectedChannels {
		b.WriteString(fmt.Sprintf("  • %s\n", name))
	}
	
	b.WriteString("\n")
	b.WriteString(m.styles.warning.Render("This action cannot be undone!"))
	b.WriteString("\n\n")
	b.WriteString(m.styles.subtitle.Render("Press y to confirm, n to cancel"))
	
	return m.styles.border.Render(b.String())
}

func (m model) renderSkipListScreen() string {
	var b strings.Builder
	
	switch m.skipListMode {
	case "view":
		return m.renderSkipListView()
	case "add":
		return m.renderSkipListAdd()
	case "remove":
		return m.renderSkipListRemove()
	}
	return m.getResponsiveBorder().Render(b.String())
}

func (m model) renderSkipListView() string {
	var b strings.Builder
	
	b.WriteString(m.styles.title.Render("📝 Skip List"))
	b.WriteString("\n\n")
	
	if len(m.skipChoices) == 0 {
		b.WriteString(m.styles.info.Render("No channels in skip list."))
		b.WriteString("\n\n")
		b.WriteString(m.styles.subtitle.Render("Press 'a' to add channels, Enter to return to main menu"))
		return m.getResponsiveBorder().Render(b.String())
	}
	
	// Show paginated results (10 items per page)
	start := m.skipListOffset
	end := start + 10
	if end > len(m.skipChoices) {
		end = len(m.skipChoices)
	}
	
	visibleChoices := m.skipChoices[start:end]
	
	for i, name := range visibleChoices {
		cursor := " "
		if m.skipCursor == start+i {
			cursor = m.styles.cursor.Render(">")
		}
		b.WriteString(fmt.Sprintf("%s %s\n", cursor, name))
	}
	
	// Show pagination info
	if len(m.skipChoices) > 10 {
		currentPage := (m.skipListOffset / 10) + 1
		totalPages := (len(m.skipChoices) + 9) / 10
		b.WriteString(fmt.Sprintf("\n%s\n", m.styles.info.Render(fmt.Sprintf("Page %d/%d", currentPage, totalPages))))
	}
	
	b.WriteString("\n")
	b.WriteString(m.styles.subtitle.Render("Use ↑↓ to navigate, 'a' to add, 'd' to delete, Enter to return"))
	
	return m.getResponsiveBorder().Render(b.String())
}

func (m model) renderSkipListAdd() string {
	var b strings.Builder
	
	b.WriteString(m.styles.title.Render("➕ Add Channel to Skip List"))
	b.WriteString("\n\n")
	
	b.WriteString("Enter channel name (without #):\n")
	b.WriteString(fmt.Sprintf("> %s", m.skipListInput))
	
	b.WriteString("\n\n")
	b.WriteString(m.styles.subtitle.Render("Type channel name and press Enter to add, q to cancel"))
	
	return m.getResponsiveBorder().Render(b.String())
}

func (m model) renderSkipListRemove() string {
	var b strings.Builder
	
	b.WriteString(m.styles.title.Render("🗑️  Remove Channel from Skip List"))
	b.WriteString("\n\n")
	
	if len(m.skipChoices) == 0 {
		b.WriteString(m.styles.info.Render("No channels to remove."))
		b.WriteString("\n\n")
		b.WriteString(m.styles.subtitle.Render("Press q to return"))
		return m.getResponsiveBorder().Render(b.String())
	}
	
	// Show paginated results (10 items per page)
	start := m.skipListOffset
	end := start + 10
	if end > len(m.skipChoices) {
		end = len(m.skipChoices)
	}
	
	visibleChoices := m.skipChoices[start:end]
	
	for i, name := range visibleChoices {
		cursor := " "
		if m.skipCursor == start+i {
			cursor = m.styles.cursor.Render(">")
		}
		b.WriteString(fmt.Sprintf("%s %s\n", cursor, name))
	}
	
	// Show pagination info
	if len(m.skipChoices) > 10 {
		currentPage := (m.skipListOffset / 10) + 1
		totalPages := (len(m.skipChoices) + 9) / 10
		b.WriteString(fmt.Sprintf("\n%s\n", m.styles.info.Render(fmt.Sprintf("Page %d/%d", currentPage, totalPages))))
	}
	
	b.WriteString("\n")
	b.WriteString(m.styles.warning.Render("Press Enter to remove selected channel"))
	b.WriteString("\n")
	b.WriteString(m.styles.subtitle.Render("Use ↑↓ to navigate, q to cancel"))
	
	return m.getResponsiveBorder().Render(b.String())
}

func (m model) renderLoadingScreen() string {
	var b strings.Builder
	
	b.WriteString(m.styles.title.Render("⏳ Loading"))
	b.WriteString("\n\n")
	
	b.WriteString(m.styles.info.Render(m.loadingMsg))
	b.WriteString("\n")
	b.WriteString(m.styles.subtitle.Render("Please wait..."))
	
	return m.getResponsiveBorder().Render(b.String())
}

// Message types for communication between components
type errorMsg struct {
	err error
}

func (e errorMsg) Error() string {
	return e.err.Error()
}

type channelsLoadedMsg struct {
	channels []slack.ChannelInfo
}

type channelsLeftMsg struct {}

func (m model) renderSimpleListView(visibleChannels []slack.ChannelInfo, start int) string {
	var b strings.Builder
	
	// Add a simple border
	b.WriteString("┌")
	b.WriteString(strings.Repeat("─", m.width-6))
	b.WriteString("┐\n")
	
	// Header
	b.WriteString("│ ")
	b.WriteString(m.styles.title.Render("Select"))
	b.WriteString(" │ ")
	b.WriteString(m.styles.title.Render("Channel"))
	b.WriteString(strings.Repeat(" ", m.width-30))
	b.WriteString(" │ ")
	b.WriteString(m.styles.title.Render("Last Activity"))
	b.WriteString(" │\n")
	
	// Separator
	b.WriteString("├")
	b.WriteString(strings.Repeat("─", m.width-6))
	b.WriteString("┤\n")
	
	// Channel rows
	for i, ch := range visibleChannels {
		globalIndex := start + i
		cursor := " "
		if m.cursor == globalIndex {
			cursor = m.styles.cursor.Render(">")
		}
		
		checked := " "
		if _, ok := m.selected[globalIndex]; ok {
			checked = m.styles.selected.Render("✓")
		}
		
		lastSeen := "No messages"
		if !ch.LastSeen.IsZero() {
			lastSeen = ch.LastSeen.Format("2006-01-02 15:04:05")
		}
		
		// Truncate name if too long
		name := ch.Name
		maxNameWidth := m.width - 50 // Leave space for other columns
		if len(name) > maxNameWidth {
			name = name[:maxNameWidth-3] + "..."
		}
		
		b.WriteString("│ ")
		b.WriteString(fmt.Sprintf("%s [%s]", cursor, checked))
		b.WriteString(" │ ")
		b.WriteString(fmt.Sprintf("#%-*s", maxNameWidth, name))
		b.WriteString(" │ ")
		b.WriteString(lastSeen)
		b.WriteString(" │\n")
	}
	
	// Bottom border
	b.WriteString("└")
	b.WriteString(strings.Repeat("─", m.width-6))
	b.WriteString("┘\n")
	
	return b.String()
}
