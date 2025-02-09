package v1feed

// IPSW contains information about macOS IPSW installers
type IPSW struct {
	// URL is the download location for the IPSW
	URL string `json:"macos_ipsw_url"`
	// Build is the build number of the IPSW
	Build string `json:"macos_ipsw_build"`
	// Version is the macOS version of the IPSW
	Version string `json:"macos_ipsw_version"`
	// AppleSlug is Apple's identifier for this IPSW
	AppleSlug string `json:"macos_ipsw_apple_slug"`
}
