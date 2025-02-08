package v1feed

type SupportedModel struct {
	Model       string        `json:"Model"`
	URL         string        `json:"URL"`
	Identifiers IdentifierMap `json:"Identifiers"`
}
