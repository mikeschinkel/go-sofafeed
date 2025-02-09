package v1feed

// InstallationApps contains information about macOS installer applications
type InstallationApps struct {
	// LatestUMA contains information about the latest Universal macOS Installer
	LatestUMA UMA `json:"LatestUMA"`
	// AllPreviousUMA contains historical Universal macOS Installer information
	AllPreviousUMA []UMA `json:"AllPreviousUMA"`
	// LatestMacIPSW contains information about the latest macOS IPSW installer
	LatestMacIPSW IPSW `json:"LatestMacIPSW"`
}
