package v1feed

// OSVersion contains information about a specific OS version
type OSVersion struct {
	// OSVersion is the major version number
	OSVersion string `json:"OSVersion"`
	// Latest contains information about the most recent release
	Latest OSVersionDetails `json:"Latest"`
	// SecurityReleases contains information about security updates
	SecurityReleases []SecurityReleases `json:"SecurityReleases"`
	// SupportedModels lists compatible device models (macOS-specific)
	SupportedModels []SupportedModel `json:"SupportedModels,omitempty"`
}
