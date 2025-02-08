package sofafeed_test

import (
	"bytes"
	"context"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/mikeschinkel/go-sofafeed"
	"github.com/mikeschinkel/go-sofafeed/feeds/v1feed"
)

//goland:noinspection SpellCheckingInspection
const (
	expectedUpdateHash = "60f722820a1bc63e15b540c3e6eb4f528fb3a3237b8ecd5bfb9a0669f60bafb7"
)

var (
	expectedXProtect = v1feed.XProtectPayloads{
		XProtect:      "149",
		PluginService: "74",
		ReleaseDate:   time.Date(2024, 12, 17, 18, 3, 20, 0, time.UTC),
	}
	expectedPlistConfig = v1feed.XProtectPlistConfigData{
		XProtect:    "5287",
		ReleaseDate: time.Date(2025, 2, 5, 18, 35, 11, 0, time.UTC),
	}
	expectedOSVersions = []struct {
		version string
		build   string
	}{
		{"Sequoia 15", "15.3"},
		{"Sonoma 14", "14.7.3"},
		{"Monterey 12", "12.7.6"},
	}
	expectedModels = []struct {
		id            string
		marketingName string
		osVersions    []int
	}{
		{
			"Mac15,13",
			"MacBook Air (15-inch, M3, 2024)",
			[]int{15, 14},
		},
		{
			"Mac14,15",
			"MacBook Air (15-inch, M2, 2023)",
			[]int{15, 14, 13},
		},
		{
			"MacBookAir10,1",
			"MacBook Air (M1, 2020)",
			[]int{15, 14, 13, 12},
		},
	}
	expectedLatestUMA = v1feed.UMA{
		Title:     "macOS Sequoia",
		Version:   "15.3",
		Build:     "24D60",
		AppleSlug: "072-08251",
		URL:       "https://swcdn.apple.com/content/downloads/11/60/072-08251-A_VJ2TWIZ7CM/mrkdsdl45umr07vwe27n8hfvubidmxgcbk/InstallAssistant.pkg",
	}
	expectedPreviousUMAs = []v1feed.UMA{
		{
			Title:     "macOS Sequoia",
			Version:   "15.2",
			Build:     "24C101",
			AppleSlug: "072-44286",
			URL:       "https://swcdn.apple.com/content/downloads/52/42/072-44286-A_45NJHEFEDY/nxk2stn9edtjgfr3dqsh046dx4ti2h7w2s/InstallAssistant.pkg",
		},
		{
			Title:     "macOS Sequoia",
			Version:   "15.1.1",
			Build:     "24B91",
			AppleSlug: "072-30111",
			URL:       "https://swcdn.apple.com/content/downloads/21/19/072-30111-A_4V7Y0VVH1Q/ie1hmy1uaj094z769s4zqmdaojp2vk4dkj/InstallAssistant.pkg",
		},
		{
			Title:     "macOS Sequoia",
			Version:   "15.1.1",
			Build:     "24B2091",
			AppleSlug: "072-29965",
			URL:       "https://swcdn.apple.com/content/downloads/14/63/072-29965-A_JR3Q5P9LW7/qtu6yl2s14vidrhfnunfi1t7wid5dzxajm/InstallAssistant.pkg",
		},
	}
	expectedIPSW = v1feed.IPSW{
		URL:       "https://updates.cdn-apple.com/2025WinterFCS/fullrestores/072-08269/7CAAB9F7-E970-428D-8764-4CD7BCD105CD/UniversalMac_15.3_24D60_Restore.ipsw",
		Build:     "24D60",
		Version:   "15.3",
		AppleSlug: "072-08269",
	}
)

// Helper function to compare slices
func sliceEqual(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// testData holds the parsed feed that all tests can use
var testData *v1feed.Feed

// init loads the test data once for all tests
func init() {
	j, err := os.ReadFile("testdata/sofafeed-v1.json")
	if err != nil {
		panic("Failed to read test data: " + err.Error())
	}
	testData, err = sofafeed.Parse(j)
	if err != nil {
		panic("Failed to parse test data: " + err.Error())
	}
}

func TestParseFeed(t *testing.T) {
	got := testData

	// Test top-level Feed fields
	if got.UpdateHash != expectedUpdateHash {
		t.Errorf("UpdateHash = %v, want %v", got.UpdateHash, expectedUpdateHash)
	}

	// Test XProtectPayloads
	if got.XProtectPayloads != expectedXProtect {
		t.Errorf("XProtectPayloads = %v, want %v", got.XProtectPayloads, expectedXProtect)
	}

	if got.XProtectPlistConfigData != expectedPlistConfig {
		t.Errorf("XProtectPlistConfigData = %v, want %v", got.XProtectPlistConfigData, expectedPlistConfig)
	}

	// Test OSVersions
	if len(got.OSVersions) == 0 {
		t.Fatal("OSVersions is empty")
	}

	for i, expected := range expectedOSVersions {
		found := false
		for _, osv := range got.OSVersions {
			if osv.OSVersion == expected.version {
				found = true
				if osv.Latest.ProductVersion != expected.build {
					t.Errorf("OS Version %s build = %v, want %v", expected.version, osv.Latest.ProductVersion, expected.build)
				}
			}
		}
		if !found {
			t.Errorf("OS Version #%d %s not found", i, expected.version)
		}
	}

	for _, expected := range expectedModels {
		model, exists := got.Models[expected.id]
		if !exists {
			t.Errorf("Model %s not found", expected.id)
			continue
		}
		if model.MarketingName != expected.marketingName {
			t.Errorf("Model %s marketing name = %v, want %v", expected.id, model.MarketingName, expected.marketingName)
		}
		if !sliceEqual(model.OSVersions, expected.osVersions) {
			t.Errorf("Model %s OS versions = %v, want %v", expected.id, model.OSVersions, expected.osVersions)
		}
	}

	// Test InstallationApps
	latestUMA := got.InstallationApps.LatestUMA
	if latestUMA != expectedLatestUMA {
		t.Errorf("LatestUMA = %v, want %v", latestUMA, expectedLatestUMA)
	}

	// Test three previous UMAs
	if len(got.InstallationApps.AllPreviousUMA) < 3 {
		t.Fatal("AllPreviousUMA has fewer than 3 entries")
	}

	for i, expected := range expectedPreviousUMAs {
		if got.InstallationApps.AllPreviousUMA[i] != expected {
			t.Errorf("AllPreviousUMA[%d] = %v, want %v", i, got.InstallationApps.AllPreviousUMA[i], expected)
		}
	}

	if got.InstallationApps.LatestMacIPSW != expectedIPSW {
		t.Errorf("LatestMacIPSW = %v, want %v", got.InstallationApps.LatestMacIPSW, expectedIPSW)
	}
}

// TestCVEMapping validates the parsing of CVE data from various OS versions
func TestCVEMapping(t *testing.T) {
	// Test cases with known CVEs from different OS versions and releases
	testCases := []struct {
		osVersion  string
		cve        string
		exploited  bool
		inLatest   bool   // Whether to look in Latest or SecurityReleases
		releaseVer string // Which release version to find it in, if not Latest
	}{
		{
			osVersion: "Sequoia 15",
			cve:       "CVE-2025-24085",
			exploited: true,
			inLatest:  true,
		},
		{
			osVersion:  "Sonoma 14",
			cve:        "CVE-2024-23225",
			exploited:  true,
			inLatest:   false,
			releaseVer: "14.4",
		},
		{
			osVersion: "Monterey 12",
			cve:       "CVE-2024-40783",
			exploited: false,
			inLatest:  true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.osVersion+"/"+tc.cve, func(t *testing.T) {
			// Find the OS version
			var osv *v1feed.OSVersion
			for i := range testData.OSVersions {
				if testData.OSVersions[i].OSVersion == tc.osVersion {
					osv = &testData.OSVersions[i]
					break
				}
			}
			if osv == nil {
				t.Fatalf("OS version %s not found", tc.osVersion)
			}

			if tc.inLatest {
				// Check Latest section
				exploited, exists := osv.Latest.CVEs[tc.cve]
				if !exists {
					t.Errorf("CVE %s not found in %s Latest", tc.cve, tc.osVersion)
				}
				if exploited != tc.exploited {
					t.Errorf("CVE %s exploitation status = %v, want %v", tc.cve, exploited, tc.exploited)
				}

				foundInList := false
				for _, exploitedCVE := range osv.Latest.ActivelyExploitedCVEs {
					if exploitedCVE == tc.cve {
						foundInList = true
						break
					}
				}
				if tc.exploited != foundInList {
					t.Errorf("CVE %s in Latest.ActivelyExploitedCVEs = %v, want %v", tc.cve, foundInList, tc.exploited)
				}
			} else {
				// Look through SecurityReleases for specific version
				var found bool
				for _, release := range osv.SecurityReleases {
					if release.ProductVersion == tc.releaseVer {
						exploited, exists := release.CVEs[tc.cve]
						if !exists {
							continue
						}
						found = true
						if exploited != tc.exploited {
							t.Errorf("CVE %s exploitation status in %s = %v, want %v",
								tc.cve, tc.releaseVer, exploited, tc.exploited)
						}

						foundInList := false
						for _, exploitedCVE := range release.ActivelyExploitedCVEs {
							if exploitedCVE == tc.cve {
								foundInList = true
								break
							}
						}
						if tc.exploited != foundInList {
							t.Errorf("CVE %s in %s ActivelyExploitedCVEs = %v, want %v",
								tc.cve, tc.releaseVer, foundInList, tc.exploited)
						}
						break
					}
				}
				if !found {
					t.Errorf("CVE %s not found in %s version %s", tc.cve, tc.osVersion, tc.releaseVer)
				}
			}
		})
	}
}

// TestSupportedDevices validates the parsing of supported device lists
func TestSupportedDevices(t *testing.T) {
	// Test cases for specific devices that should be supported across different OS versions
	testCases := []struct {
		osVersion     string
		deviceID      string
		shouldSupport bool
	}{
		{"Sequoia 15", "Mac-1E7E29AD0135F9BC", true},  // Known supported device
		{"Sonoma 14", "Mac-937A206F2EE63C01", true},   // Known supported device
		{"Monterey 12", "Mac-FFE5EF870D7BA81A", true}, // Known supported device
		{"Sequoia 15", "Mac-INVALID-ID", false},       // Invalid device
		{"Sonoma 14", "MacBookPro11,1", false},        // Unsupported device
		{"Monterey 12", "iMac14,1", false},            // Unsupported device
	}

	for _, tc := range testCases {
		t.Run(tc.osVersion+"/"+tc.deviceID, func(t *testing.T) {
			// Find the OS version
			var osv *v1feed.OSVersion
			for i := range testData.OSVersions {
				if testData.OSVersions[i].OSVersion == tc.osVersion {
					osv = &testData.OSVersions[i]
					break
				}
			}
			if osv == nil {
				t.Fatalf("OS version %s not found", tc.osVersion)
			}

			// Check if device is supported
			found := false
			for _, device := range osv.Latest.SupportedDevices {
				if device == tc.deviceID {
					found = true
					break
				}
			}

			if found != tc.shouldSupport {
				t.Errorf("Device %s support status = %v, want %v", tc.deviceID, found, tc.shouldSupport)
			}
		})
	}
}

// TestErrorHandling validates the parser's handling of malformed JSON
func TestErrorHandling(t *testing.T) {
	testCases := []struct {
		name    string
		json    string
		wantErr bool
	}{
		{
			name:    "Empty JSON",
			json:    "",
			wantErr: true,
		},
		{
			name:    "Invalid JSON",
			json:    "{not valid json}",
			wantErr: true,
		},
		{
			name:    "Missing required field",
			json:    `{"OSVersions": []}`,
			wantErr: false, // Should not error as struct fields are optional in Go
		},
		{
			name:    "Invalid date format",
			json:    `{"XProtectPayloads": {"ReleaseDate": "not-a-date"}}`,
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := sofafeed.ParseString(tc.json)
			if (err != nil) != tc.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tc.wantErr)
			}
		})
	}
}

func TestFetch(t *testing.T) {
	// Create a test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify headers
		if r.Header.Get("Accept") != "application/json" {
			t.Errorf("Expected Accept: application/json header, got %s", r.Header.Get("Accept"))
		}
		if ua := r.Header.Get("User-Agent"); !strings.HasPrefix(ua, "go-sofafeed/") {
			t.Errorf("Expected User-Agent to start with go-sofafeed/, got %s", ua)
		}

		// Send test response
		w.Header().Set("Content-Type", "application/json")
		mustWrite(w, []byte(`{"test": "data"}`))
	}))
	defer ts.Close()

	// Temporarily override the default endpoint for testing
	originalEndpoint := sofafeed.Endpoint
	sofafeed.Endpoint = ts.URL
	defer func() { sofafeed.Endpoint = originalEndpoint }()

	// Test with nil client (should use default)
	ctx := context.Background()
	data, err := sofafeed.Fetch(ctx, nil)
	if err != nil {
		t.Fatalf("Fetch() error = %v", err)
	}
	if !bytes.Contains(data, []byte(`"test"`)) {
		t.Errorf("Fetch() response doesn't contain expected data")
	}
}

func mustWrite(w io.Writer, data []byte) {
	n, err := w.Write(data)
	if err != nil {
		slog.Error("Failed to write", "data", data, "error", err)
	}
	if n != len(data) {
		slog.Error("Failed to write all data", "data", data, "bytes_written", n)
	}
}
