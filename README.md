# eBay Draft Listing Generator

[![Build and Test](https://github.com/korjavin/ebayDraftListing/actions/workflows/build.yml/badge.svg)](https://github.com/korjavin/ebayDraftListing/actions/workflows/build.yml)

A Golang CLI application that automatically generates eBay draft listings from photos using Google's Gemini AI.

## Features

- 🤖 AI-powered title and description generation using Gemini Flash API
- 📸 Support for multiple photo formats (JPEG, PNG, GIF, WebP)
- 🔐 eBay OAuth 2.0 authentication
- 📝 Automatic draft listing creation on eBay
- 🌍 Support for both sandbox and production environments

## Prerequisites

- Go 1.21 or higher
- Google Gemini API key
- eBay Developer account with OAuth credentials

## Installation

1. Clone the repository:
```bash
git clone https://github.com/korjavin/ebayDraftListing.git
cd ebayDraftListing
```

2. Install dependencies:
```bash
go mod download
```

3. Build the application:
```bash
go build -o ebay-listing cmd/main.go
```

## Configuration

Create a `.env` file or set the following environment variables:

```bash
# Gemini API Configuration
export GEMINI_API_KEY="your-gemini-api-key"
export PROMPT="Generate a compelling eBay listing title and description for the product shown in these images. Return the result in JSON format with 'Title' and 'Description' fields."

# eBay API Configuration
export EBAY_CLIENT_ID="your-ebay-client-id"
export EBAY_CLIENT_SECRET="your-ebay-client-secret"
export EBAY_REFRESH_TOKEN="your-ebay-refresh-token"
export EBAY_ENVIRONMENT="sandbox"  # or "production"
```

### Getting eBay Credentials

1. Create a developer account at [eBay Developers Program](https://developer.ebay.com/)
2. Create an application in the [Developer Dashboard](https://developer.ebay.com/my/keys)
3. Generate OAuth credentials (Client ID and Client Secret)
4. Obtain a refresh token using eBay's OAuth flow

### Getting Gemini API Key

1. Visit [Google AI Studio](https://makersuite.google.com/app/apikey)
2. Create a new API key
3. Copy the key to your environment variables

## Usage

Run the application with one or more photo paths as arguments:

```bash
./ebay-listing photo1.jpg photo2.png photo3.jpg
```

Or use `go run`:

```bash
go run cmd/main.go photo1.jpg photo2.png
```

### Example

```bash
export GEMINI_API_KEY="AIza..."
export PROMPT="Create an eBay listing for this product. Include a catchy title (max 80 chars) and detailed description. Format as JSON with Title and Description fields."
export EBAY_CLIENT_ID="YourApp-..."
export EBAY_CLIENT_SECRET="..."
export EBAY_REFRESH_TOKEN="v^1.1#..."
export EBAY_ENVIRONMENT="sandbox"

./ebay-listing ~/products/item1.jpg ~/products/item2.jpg
```

### Output

```
Processing 2 photo(s)...
Generating listing content with Gemini AI...

=== Generated Content ===
Title: Vintage Ceramic Vase - Hand Painted Floral Design
Description:
Beautiful vintage ceramic vase featuring hand-painted floral designs.
Excellent condition with no chips or cracks. Perfect for home decor.
Dimensions: 10" tall x 5" wide.
========================

Creating draft listing on eBay...

✓ Draft listing created successfully!
Offer ID: 1234567890
Environment: sandbox
```

## Project Structure

```
ebayDraftListing/
├── cmd/
│   └── main.go              # CLI entry point
├── pkg/
│   ├── config/
│   │   └── config.go        # Environment configuration
│   ├── gemini/
│   │   └── client.go        # Gemini API integration
│   ├── ebay/
│   │   ├── auth.go          # eBay OAuth authentication
│   │   └── listing.go       # eBay Listing API
│   └── models/
│       └── listing.go       # Data models
├── go.mod
├── go.sum
└── README.md
```

## Customizing the Prompt

The `PROMPT` environment variable controls how Gemini generates your listing content. You can customize it for different product types:

**For electronics:**
```bash
export PROMPT="Analyze these product images and create an eBay listing. Include technical specifications, condition, and key features. Return JSON with Title and Description fields."
```

**For collectibles:**
```bash
export PROMPT="Create an eBay listing for this collectible item. Highlight rarity, condition, and historical significance. Return JSON with Title and Description fields."
```

**For clothing:**
```bash
export PROMPT="Generate an eBay listing for this clothing item. Include size, material, brand, condition, and style details. Return JSON with Title and Description fields."
```

## Troubleshooting

### "GEMINI_API_KEY environment variable is required"
Make sure you've exported all required environment variables before running the application.

### "failed to get access token"
Verify your eBay credentials are correct and your refresh token hasn't expired.

### "failed to upload images"
Check that the photo paths are correct and the files exist.

### "request failed with status 401"
Your eBay OAuth token may have expired. Generate a new refresh token.

## Development

Run tests:
```bash
go test ./...
```

Format code:
```bash
go fmt ./...
```

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## Disclaimer

This application creates draft listings only. Review all generated content before publishing to ensure accuracy and compliance with eBay policies.
