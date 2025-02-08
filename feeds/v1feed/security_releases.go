package v1feed

import (
	"time"
)

type SecurityReleases struct {
	UpdateName               string    `json:"UpdateName"`
	ProductName              string    `json:"ProductName"`
	ProductVersion           string    `json:"ProductVersion"`
	ReleaseDate              time.Time `json:"ReleaseDate"`
	ReleaseType              string    `json:"ReleaseType"`
	SecurityInfo             string    `json:"SecurityInfo"`
	SupportedDevices         []string  `json:"SupportedDevices"`
	CVEs                     CVEMap    `json:"CVEs"`
	ActivelyExploitedCVEs    []string  `json:"ActivelyExploitedCVEs"`
	UniqueCVEsCount          int       `json:"UniqueCVEsCount"`
	DaysSincePreviousRelease int       `json:"DaysSincePreviousRelease"`
}
