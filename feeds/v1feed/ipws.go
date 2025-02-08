package v1feed

type IPSW struct {
	URL       string `json:"macos_ipsw_url"`
	Build     string `json:"macos_ipsw_build"`
	Version   string `json:"macos_ipsw_version"`
	AppleSlug string `json:"macos_ipsw_apple_slug"`
}
