package v1feed

import (
	"time"
)

type XProtectPlistConfigData struct {
	XProtect    string    `json:"com.apple.XProtect"`
	ReleaseDate time.Time `json:"ReleaseDate"`
}
