package ebay

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/korjavin/ebayDraftListing/pkg/models"
)

const (
	sandboxInventoryURL    = "https://api.sandbox.ebay.com/sell/inventory/v1"
	productionInventoryURL = "https://api.ebay.com/sell/inventory/v1"
)

// ListingClient handles eBay listing operations
type ListingClient struct {
	authClient *AuthClient
	baseURL    string
}

// NewListingClient creates a new eBay listing client
func NewListingClient(authClient *AuthClient) *ListingClient {
	baseURL := sandboxInventoryURL
	if authClient.environment == "production" {
		baseURL = productionInventoryURL
	}

	return &ListingClient{
		authClient: authClient,
		baseURL:    baseURL,
	}
}

type inventoryItem struct {
	Product      product      `json:"product"`
	Condition    string       `json:"condition"`
	Availability availability `json:"availability"`
}

type product struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	ImageURLs   []string `json:"imageUrls,omitempty"`
	Aspects     aspects  `json:"aspects,omitempty"`
}

type aspects struct {
	Brand []string `json:"Brand,omitempty"`
}

type availability struct {
	ShipToLocationAvailability shipToLocationAvailability `json:"shipToLocationAvailability"`
}

type shipToLocationAvailability struct {
	Quantity int `json:"quantity"`
}

type offer struct {
	SKU                 string          `json:"sku"`
	MarketplaceID       string          `json:"marketplaceId"`
	Format              string          `json:"format"`
	ListingDescription  string          `json:"listingDescription,omitempty"`
	PricingSummary      pricingSummary  `json:"pricingSummary"`
	ListingPolicies     listingPolicies `json:"listingPolicies"`
	CategoryID          string          `json:"categoryId"`
	MerchantLocationKey string          `json:"merchantLocationKey,omitempty"`
}

type pricingSummary struct {
	Price price `json:"price"`
}

type price struct {
	Value    string `json:"value"`
	Currency string `json:"currency"`
}

type listingPolicies struct {
	FulfillmentPolicyID string `json:"fulfillmentPolicyId,omitempty"`
	PaymentPolicyID     string `json:"paymentPolicyId,omitempty"`
	ReturnPolicyID      string `json:"returnPolicyId,omitempty"`
}

// CreateDraftListing creates a draft listing on eBay
func (l *ListingClient) CreateDraftListing(listing *models.DraftListing) (string, error) {
	// Get access token
	accessToken, err := l.authClient.GetAccessToken()
	if err != nil {
		return "", fmt.Errorf("failed to get access token: %w", err)
	}

	// Upload images first
	imageURLs, err := l.uploadImages(accessToken, listing.PhotoPaths)
	if err != nil {
		return "", fmt.Errorf("failed to upload images: %w", err)
	}

	// Generate a SKU
	sku := fmt.Sprintf("DRAFT-%d", len(listing.Title))

	// Create inventory item
	if err := l.createInventoryItem(accessToken, sku, listing, imageURLs); err != nil {
		return "", fmt.Errorf("failed to create inventory item: %w", err)
	}

	// Create offer (draft listing)
	offerID, err := l.createOffer(accessToken, sku, listing)
	if err != nil {
		return "", fmt.Errorf("failed to create offer: %w", err)
	}

	return offerID, nil
}

// createInventoryItem creates an inventory item
func (l *ListingClient) createInventoryItem(accessToken, sku string, listing *models.DraftListing, imageURLs []string) error {
	item := inventoryItem{
		Product: product{
			Title:       listing.Title,
			Description: listing.Description,
			ImageURLs:   imageURLs,
		},
		Condition: "NEW",
		Availability: availability{
			ShipToLocationAvailability: shipToLocationAvailability{
				Quantity: 1,
			},
		},
	}

	jsonData, err := json.Marshal(item)
	if err != nil {
		return fmt.Errorf("failed to marshal inventory item: %w", err)
	}

	url := fmt.Sprintf("%s/inventory_item/%s", l.baseURL, sku)
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Content-Language", "en-US")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("request failed with status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

// createOffer creates an offer (draft listing)
func (l *ListingClient) createOffer(accessToken, sku string, listing *models.DraftListing) (string, error) {
	offer := offer{
		SKU:                sku,
		MarketplaceID:      "EBAY_US",
		Format:             "FIXED_PRICE",
		ListingDescription: listing.Description,
		PricingSummary: pricingSummary{
			Price: price{
				Value:    "9.99",
				Currency: "USD",
			},
		},
		CategoryID:      "111422", // Default category
		ListingPolicies: listingPolicies{},
	}

	jsonData, err := json.Marshal(offer)
	if err != nil {
		return "", fmt.Errorf("failed to marshal offer: %w", err)
	}

	url := fmt.Sprintf("%s/offer", l.baseURL)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Content-Language", "en-US")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return "", fmt.Errorf("request failed with status %d: %s", resp.StatusCode, string(body))
	}

	// Parse response to get offer ID
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("failed to parse response: %w", err)
	}

	offerID, ok := result["offerId"].(string)
	if !ok {
		return "", fmt.Errorf("offer ID not found in response")
	}

	return offerID, nil
}

// uploadImages uploads images and returns their URLs
func (l *ListingClient) uploadImages(accessToken string, photoPaths []string) ([]string, error) {
	var imageURLs []string

	// For eBay, we need to use the Trading API or host images externally
	// This is a simplified version that converts images to base64 data URLs
	// In production, you would want to use eBay's picture services or host images externally

	for _, photoPath := range photoPaths {
		// Read image file
		data, err := os.ReadFile(photoPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read image %s: %w", photoPath, err)
		}

		// Determine MIME type
		mimeType := "image/jpeg"
		ext := filepath.Ext(photoPath)
		switch ext {
		case ".png":
			mimeType = "image/png"
		case ".jpg", ".jpeg":
			mimeType = "image/jpeg"
		case ".gif":
			mimeType = "image/gif"
		}

		// Create data URL
		encoded := base64.StdEncoding.EncodeToString(data)
		dataURL := fmt.Sprintf("data:%s;base64,%s", mimeType, encoded)
		imageURLs = append(imageURLs, dataURL)
	}

	return imageURLs, nil
}
