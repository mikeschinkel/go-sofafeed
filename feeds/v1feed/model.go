package v1feed

// Model contains information about a specific device model
type Model struct {
	// MarketingName is the consumer-facing name of the device
	MarketingName string `json:"MarketingName"`
	// SupportedOS lists compatible OS versions
	SupportedOS []string `json:"SupportedOS"`
	// OSVersions lists compatible OS version numbers
	OSVersions []int `json:"OSVersions"`
}
