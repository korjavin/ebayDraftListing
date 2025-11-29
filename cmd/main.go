package main

import (
	"fmt"
	"os"

	"github.com/korjavin/ebayDraftListing/pkg/config"
	"github.com/korjavin/ebayDraftListing/pkg/ebay"
	"github.com/korjavin/ebayDraftListing/pkg/gemini"
	"github.com/korjavin/ebayDraftListing/pkg/models"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	// Check for photo arguments
	if len(os.Args) < 2 {
		return fmt.Errorf("usage: %s <photo1> [photo2] [photo3] ...", os.Args[0])
	}

	photoPaths := os.Args[1:]

	// Validate photo paths
	for _, path := range photoPaths {
		if _, err := os.Stat(path); err != nil {
			return fmt.Errorf("photo not found: %s", path)
		}
	}

	fmt.Printf("Processing %d photo(s)...\n", len(photoPaths))

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	// Generate listing content using Gemini
	fmt.Println("Generating listing content with Gemini AI...")
	geminiClient := gemini.NewClient(cfg.GeminiAPIKey, cfg.Prompt)
	content, err := geminiClient.GenerateListingContent(photoPaths)
	if err != nil {
		return fmt.Errorf("failed to generate listing content: %w", err)
	}

	fmt.Printf("\n=== Generated Content ===\n")
	fmt.Printf("Title: %s\n", content.Title)
	fmt.Printf("Description:\n%s\n", content.Description)
	fmt.Printf("========================\n\n")

	// Create draft listing on eBay
	fmt.Println("Creating draft listing on eBay...")
	authClient := ebay.NewAuthClient(
		cfg.EbayClientID,
		cfg.EbayClientSecret,
		cfg.EbayRefreshToken,
		cfg.EbayEnvironment,
	)

	listingClient := ebay.NewListingClient(authClient)

	draftListing := &models.DraftListing{
		Title:       content.Title,
		Description: content.Description,
		PhotoPaths:  photoPaths,
	}

	offerID, err := listingClient.CreateDraftListing(draftListing)
	if err != nil {
		return fmt.Errorf("failed to create draft listing: %w", err)
	}

	fmt.Printf("\n✓ Draft listing created successfully!\n")
	fmt.Printf("Offer ID: %s\n", offerID)
	fmt.Printf("Environment: %s\n", cfg.EbayEnvironment)

	return nil
}
