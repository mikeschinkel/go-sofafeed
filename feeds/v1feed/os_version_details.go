package v1feed

import (
	"time"
)

// OSVersionDetails contains detailed information about an OS release
type OSVersionDetails struct {
	// ProductVersion is the full version number
	ProductVersion string `json:"ProductVersion"`
	// Build is the build number
	Build string `json:"Build"`
	// ReleaseDate is when the version was released
	ReleaseDate time.Time `json:"ReleaseDate"`
	// ExpirationDate is when the version expires
	ExpirationDate time.Time `json:"ExpirationDate"`
	// SupportedDevices lists compatible device identifiers
	SupportedDevices []string `json:"SupportedDevices"`
	// SecurityInfo is a URL to security release notes
	SecurityInfo string `json:"SecurityInfo"`
	// CVEs maps CVE identifiers to their status
	CVEs CVEs `json:"CVEs"`
	// ActivelyExploitedCVEs lists CVEs known to be actively exploited
	ActivelyExploitedCVEs []string `json:"ActivelyExploitedCVEs"`
	// UniqueCVEsCount is the number of unique CVEs addressed
	UniqueCVEsCount int `json:"UniqueCVEsCount"`
}
