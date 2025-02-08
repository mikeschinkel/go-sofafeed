package v1feed

import (
	"time"
)

type XProtectPayloads struct {
	XProtect      string    `json:"com.apple.XProtectFramework.XProtect"`
	PluginService string    `json:"com.apple.XprotectFramework.PluginService"`
	ReleaseDate   time.Time `json:"ReleaseDate"`
}
