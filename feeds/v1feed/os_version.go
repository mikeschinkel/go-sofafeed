package v1feed

type OSVersion struct {
	OSVersion        string             `json:"OSVersion"`
	Latest           OSVersionDetails   `json:"Latest"`
	SecurityReleases []SecurityReleases `json:"SecurityReleases"`
	SupportedModels  []SupportedModel   `json:"SupportedModels"`
}
