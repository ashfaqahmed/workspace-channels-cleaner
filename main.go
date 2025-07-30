/*
 * Workspace Channel Cleaner
 * 
 * Copyright (c) 2024 Ashfaque Ali
 * 
 * This software is NOT affiliated with, endorsed by, or sponsored by Slack Technologies, Inc.
 * Slack is a registered trademark of Slack Technologies, Inc.
 * 
 * This project uses Slack's publicly available API in accordance with their terms of service.
 * Users are responsible for ensuring their use complies with Slack's terms and applicable laws.
 * 
 * For complete legal information, see DISCLAIMER.md
 */

package main

import (
	"fmt"
	"log"
	"os"

	"workspace-channel-cleaner/config"
	"workspace-channel-cleaner/model"

	"github.com/charmbracelet/bubbletea"
)

func main() {
	if err := config.LoadEnvironment(); err != nil {
		log.Printf("Warning: Could not load .env file: %v", err)
	}

	if err := config.ValidateToken(); err != nil {
		fmt.Printf("❌ %s\n", err.Error())
		fmt.Println("Please set your SLACK_API_TOKEN in the .env file or environment variables.")
		os.Exit(1)
	}

	p := tea.NewProgram(
		model.InitialModel(),
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
		tea.WithFPS(60),
	)

	if _, err := p.Run(); err != nil {
		fmt.Printf("❌ Error running TUI: %v\n", err)
		os.Exit(1)
	}
}
