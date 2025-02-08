package v1feed

type InstallationApps struct {
	LatestUMA      UMA   `json:"LatestUMA"`
	AllPreviousUMA []UMA `json:"AllPreviousUMA"`
	LatestMacIPSW  IPSW  `json:"LatestMacIPSW"`
}
