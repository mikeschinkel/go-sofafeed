package v1feed

type Model struct {
	MarketingName string   `json:"MarketingName,omitempty"`
	SupportedOS   []string `json:"SupportedOS,omitempty"`
	OSVersions    []int    `json:"OSVersions,omitempty"`
}
