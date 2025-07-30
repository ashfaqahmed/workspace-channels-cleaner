package slack

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/slack-go/slack"
)

type ChannelInfo struct {
	ID       string
	Name     string
	LastSeen time.Time
	Type     string
}

type Cleaner struct {
	API          *slack.Client
	SkipChannels map[string]bool
	Limit        int
	Types        []string
	Cutoff       time.Time
	Keyword      string
	Verbose      bool
}

// NewCleaner creates a new Slack cleaner instance
func NewCleaner(token string, limit int, types []string, days int, keyword string, verbose bool) *Cleaner {
	api := slack.New(token)
	skipChannels, _ := LoadSkipList("config/skiplist.json")
	
	cutoff := time.Now().AddDate(0, 0, -days)
	
	return &Cleaner{
		API:          api,
		SkipChannels: skipChannels,
		Limit:        limit,
		Types:        types,
		Cutoff:       cutoff,
		Keyword:      keyword,
		Verbose:      verbose,
	}
}

// LoadSkipList loads the skip list from a JSON file
func LoadSkipList(path string) (map[string]bool, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return make(map[string]bool), nil // Return empty map if file doesn't exist
	}
	
	var list []string
	if err := json.Unmarshal(data, &list); err != nil {
		return nil, fmt.Errorf("failed to parse skip list: %w", err)
	}
	
	skipMap := make(map[string]bool)
	for _, name := range list {
		skipMap[strings.TrimSpace(name)] = true
	}
	return skipMap, nil
}

// SaveSkipList saves the skip list to a JSON file
func SaveSkipList(path string, skipList map[string]bool) error {
	var list []string
	for name := range skipList {
		list = append(list, name)
	}
	
	data, err := json.MarshalIndent(list, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal skip list: %w", err)
	}
	
	return os.WriteFile(path, data, 0644)
}

// GetFilteredChannels retrieves and filters channels based on criteria
func (c *Cleaner) GetFilteredChannels() ([]ChannelInfo, error) {
	var results []ChannelInfo
	var wg sync.WaitGroup
	chMutex := sync.Mutex{}
	semaphore := make(chan struct{}, 5)
	cursor := ""

	for {
		var channels []slack.Channel
		var nextCursor string
		var err error

		for {
			params := &slack.GetConversationsParameters{
				Limit:           c.Limit,
				ExcludeArchived: true,
				Cursor:          cursor,
				Types:           c.Types,
			}

			channels, nextCursor, err = c.API.GetConversations(params)
			if err != nil {
				if rateErr := c.handleRateLimit(err); rateErr != nil {
					return nil, rateErr
				}
				continue
			}
			break
		}

		for _, ch := range channels {
			if !ch.IsMember || c.SkipChannels[ch.Name] {
				continue
			}
			if c.Keyword != "" && !strings.Contains(ch.Name, c.Keyword) {
				continue
			}

			wg.Add(1)
			go func(ch slack.Channel) {
				defer wg.Done()
				semaphore <- struct{}{}
				defer func() { <-semaphore }()

				var lastTime time.Time
				var history *slack.GetConversationHistoryResponse
				
				for attempt := 0; attempt < 2; attempt++ {
					history, err = c.API.GetConversationHistory(&slack.GetConversationHistoryParameters{
						ChannelID: ch.ID,
						Limit:     1,
					})
					if err != nil {
						if rateErr := c.handleRateLimit(err); rateErr != nil {
							return
						}
						continue
					}
					break
				}
				
				if history == nil || len(history.Messages) == 0 {
					return
				}

				tsFloat, err := strconv.ParseFloat(history.Messages[0].Timestamp, 64)
				if err != nil {
					return
				}
				lastTime = time.Unix(int64(tsFloat), 0)
				
				// Only include if older than cutoff
				if !lastTime.Before(c.Cutoff) {
					return
				}

				chMutex.Lock()
				channelType := "public"
				if ch.IsPrivate {
					channelType = "private"
				}
				results = append(results, ChannelInfo{
					ID:       ch.ID,
					Name:     ch.Name,
					LastSeen: lastTime,
					Type:     channelType,
				})
				chMutex.Unlock()
				time.Sleep(1 * time.Second)
			}(ch)
		}

		if nextCursor == "" {
			break
		}
		cursor = nextCursor
	}

	wg.Wait()
	return results, nil
}

// LeaveChannels leaves the specified channels
func (c *Cleaner) LeaveChannels(channels []ChannelInfo) error {
	for i, ch := range channels {
		if c.Verbose {
			fmt.Printf("➡️  [%d/%d] Leaving #%s (ID: %s)...\n", i+1, len(channels), ch.Name, ch.ID)
		}
		
		_, err := c.API.LeaveConversation(ch.ID)
		if err != nil {
			if rateErr := c.handleRateLimit(err); rateErr != nil {
				return fmt.Errorf("failed to leave #%s: %w", ch.Name, rateErr)
			}
			return fmt.Errorf("failed to leave #%s: %w", ch.Name, err)
		}
		
		if c.Verbose {
			fmt.Printf("✅ Left #%s\n", ch.Name)
		}
		
		time.Sleep(1 * time.Second)
	}
	return nil
}


func (c *Cleaner) handleRateLimit(err error) error {
	if strings.Contains(err.Error(), "rate_limited") {
		if c.Verbose {
			fmt.Println("⚠️  Hit rate limit. Waiting 30s...")
		}
		time.Sleep(30 * time.Second)
		return nil
	}
	
	if rateErr, ok := err.(*slack.RateLimitedError); ok {
		wait := time.Duration(rateErr.RetryAfter) * time.Second
		if wait <= 0 {
			wait = 30 * time.Second
		}
		if c.Verbose {
			fmt.Printf("⏳ Rate limit hit. Waiting %v before retrying...\n", wait)
		}
		time.Sleep(wait)
		return nil
	}
	
	return err
}

// GetChannelTypes converts user-friendly types to workspace API types
func GetChannelTypes(types []string) []string {
	var result []string
	for _, t := range types {
		t = strings.TrimSpace(strings.ToLower(t))
		switch t {
		case "public":
			result = append(result, "public_channel")
		case "private":
			result = append(result, "private_channel")
		}
	}
	if len(result) == 0 {
		result = []string{"public_channel"} // Default to public
	}
	return result
}
