package v1feed

// SupportedModel contains information about device compatibility
type SupportedModel struct {
	// Model is the device identifier
	Model string `json:"Model"`
	// URL is a link to device information
	URL string `json:"URL"`
	// Identifiers maps various device identifiers
	Identifiers Identifiers `json:"Identifiers"`
}
