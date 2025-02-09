package v1feed

import (
	"time"
)

// XProtectPayloads contains XProtect framework update information
type XProtectPayloads struct {
	// XProtect is the framework version
	XProtect string `json:"com.apple.XProtectFramework.XProtect"`
	// PluginService is the plugin service version
	PluginService string `json:"com.apple.XprotectFramework.PluginService"`
	// ReleaseDate is when the update was released
	ReleaseDate time.Time `json:"ReleaseDate"`
}
