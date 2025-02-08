package v1feed

type Feed struct {
	UpdateHash              string                  `json:"UpdateHash"`
	OSVersions              []OSVersion             `json:"OSVersions"`
	XProtectPayloads        XProtectPayloads        `json:"XProtectPayloads"`
	XProtectPlistConfigData XProtectPlistConfigData `json:"XProtectPlistConfigData"`
	Models                  ModelMap                `json:"Models"`
	InstallationApps        InstallationApps        `json:"InstallationApps"`
}
