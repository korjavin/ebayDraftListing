package config

import (
	"fmt"
	"os"
)

// Config holds all configuration values from environment variables
type Config struct {
	// Gemini API configuration
	GeminiAPIKey string
	Prompt       string

	// eBay API configuration
	EbayClientID     string
	EbayClientSecret string
	EbayRefreshToken string
	EbayEnvironment  string // "production" or "sandbox"
}

// Load reads configuration from environment variables
func Load() (*Config, error) {
	config := &Config{
		GeminiAPIKey:     os.Getenv("GEMINI_API_KEY"),
		Prompt:           os.Getenv("EBAY_PROMPT"),
		EbayClientID:     os.Getenv("EBAY_CLIENT_ID"),
		EbayClientSecret: os.Getenv("EBAY_CLIENT_SECRET"),
		EbayRefreshToken: os.Getenv("EBAY_REFRESH_TOKEN"),
		EbayEnvironment:  os.Getenv("EBAY_ENVIRONMENT"),
	}

	// Validate required fields
	if config.GeminiAPIKey == "" {
		return nil, fmt.Errorf("GEMINI_API_KEY environment variable is required")
	}

	if config.Prompt == "" {
		return nil, fmt.Errorf("EBAY_PROMPT environment variable is required")
	}

	if config.EbayClientID == "" {
		return nil, fmt.Errorf("EBAY_CLIENT_ID environment variable is required")
	}

	if config.EbayClientSecret == "" {
		return nil, fmt.Errorf("EBAY_CLIENT_SECRET environment variable is required")
	}

	if config.EbayRefreshToken == "" {
		return nil, fmt.Errorf("EBAY_REFRESH_TOKEN environment variable is required")
	}

	// Default to sandbox if not specified
	if config.EbayEnvironment == "" {
		config.EbayEnvironment = "sandbox"
	}

	return config, nil
}
