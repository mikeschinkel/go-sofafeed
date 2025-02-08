package v1feed

import (
	"time"
)

type OSVersionDetails struct {
	ProductVersion        string    `json:"ProductVersion"`
	Build                 string    `json:"Build"`
	ReleaseDate           time.Time `json:"ReleaseDate"`
	ExpirationDate        time.Time `json:"ExpirationDate"`
	SupportedDevices      []string  `json:"SupportedDevices"`
	SecurityInfo          string    `json:"SecurityInfo"`
	CVEs                  CVEMap    `json:"CVEs"`
	ActivelyExploitedCVEs []string  `json:"ActivelyExploitedCVEs"`
	UniqueCVEsCount       int       `json:"UniqueCVEsCount"`
}
