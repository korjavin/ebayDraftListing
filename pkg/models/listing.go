package models

// ListingContent represents the generated content for a listing
type ListingContent struct {
	Title       string
	Description string
}

// DraftListing represents the data needed to create an eBay draft listing
type DraftListing struct {
	Title       string
	Description string
	PhotoPaths  []string
}
