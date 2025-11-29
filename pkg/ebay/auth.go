package ebay

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

const (
	sandboxTokenURL    = "https://api.sandbox.ebay.com/identity/v1/oauth2/token"
	productionTokenURL = "https://api.ebay.com/identity/v1/oauth2/token"
)

// AuthClient handles eBay OAuth authentication
type AuthClient struct {
	clientID     string
	clientSecret string
	refreshToken string
	environment  string
}

// NewAuthClient creates a new eBay auth client
func NewAuthClient(clientID, clientSecret, refreshToken, environment string) *AuthClient {
	return &AuthClient{
		clientID:     clientID,
		clientSecret: clientSecret,
		refreshToken: refreshToken,
		environment:  environment,
	}
}

type tokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	TokenType   string `json:"token_type"`
}

// GetAccessToken obtains an access token using the refresh token
func (a *AuthClient) GetAccessToken() (string, error) {
	tokenURL := sandboxTokenURL
	if a.environment == "production" {
		tokenURL = productionTokenURL
	}

	// Prepare request body
	data := url.Values{}
	data.Set("grant_type", "refresh_token")
	data.Set("refresh_token", a.refreshToken)

	req, err := http.NewRequest("POST", tokenURL, bytes.NewBufferString(data.Encode()))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	auth := base64.StdEncoding.EncodeToString([]byte(a.clientID + ":" + a.clientSecret))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "Basic "+auth)

	// Make request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to make token request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("token request failed with status %d: %s", resp.StatusCode, string(body))
	}

	// Parse response
	var tokenResp tokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return "", fmt.Errorf("failed to parse token response: %w", err)
	}

	return tokenResp.AccessToken, nil
}
