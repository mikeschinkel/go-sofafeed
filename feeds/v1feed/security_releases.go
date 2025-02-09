package v1feed

import (
	"time"
)

// SecurityReleases contains information about security updates
type SecurityReleases struct {
	// UpdateName is the user-facing name of the update
	UpdateName string `json:"UpdateName"`
	// ProductName is the name of the operating system
	ProductName string `json:"ProductName"`
	// ProductVersion is the full version number
	ProductVersion string `json:"ProductVersion"`
	// ReleaseDate is when the update was released
	ReleaseDate time.Time `json:"ReleaseDate"`
	// ReleaseType indicates the type of release
	ReleaseType string `json:"ReleaseType"`
	// SecurityInfo is a URL to security release notes
	SecurityInfo string `json:"SecurityInfo"`
	// SupportedDevices lists compatible device identifiers
	SupportedDevices []string `json:"SupportedDevices"`
	// CVEs maps CVE identifiers to their status
	CVEs CVEs `json:"CVEs"`
	// ActivelyExploitedCVEs lists CVEs known to be actively exploited
	ActivelyExploitedCVEs []string `json:"ActivelyExploitedCVEs"`
	// UniqueCVEsCount is the number of unique CVEs addressed
	UniqueCVEsCount int `json:"UniqueCVEsCount"`
	// DaysSincePreviousRelease is days elapsed since last release
	DaysSincePreviousRelease int `json:"DaysSincePreviousRelease"`
}
