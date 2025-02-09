package v1feed

// UMA contains information about Universal macOS Applications
type UMA struct {
	// Title is the name of the installer
	Title string `json:"title"`
	// Version is the OS version
	Version string `json:"version"`
	// Build is the build number
	Build string `json:"build"`
	// AppleSlug is Apple's identifier for this installer
	AppleSlug string `json:"apple_slug"`
	// URL is the download location
	URL string `json:"url"`
}
