package v1feed

import (
	"time"
)

// XProtectPlistConfigData contains XProtect configuration information
type XProtectPlistConfigData struct {
	// XProtect is the configuration version
	XProtect string `json:"com.apple.XProtect"`
	// ReleaseDate is when the configuration was released
	ReleaseDate time.Time `json:"ReleaseDate"`
}
