package v1feed

type UMA struct {
	Title     string `json:"title"`
	Version   string `json:"version"`
	Build     string `json:"build"`
	AppleSlug string `json:"apple_slug"`
	URL       string `json:"url"`
}
