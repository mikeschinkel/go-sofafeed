package sofafeed_test

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/mikeschinkel/go-sofafeed"
	"github.com/mikeschinkel/go-sofafeed/feeds"
	"github.com/mikeschinkel/go-sofafeed/feeds/v1feed"
)

//goland:noinspection SpellCheckingInspection
const (
	expectedUpdateHash = "60f722820a1bc63e15b540c3e6eb4f528fb3a3237b8ecd5bfb9a0669f60bafb7"
)

func TestMacOSParseFeed(t *testing.T) {
	var (
		expectedXProtect = &v1feed.XProtectPayloads{
			XProtect:      "149",
			PluginService: "74",
			ReleaseDate:   time.Date(2024, 12, 17, 18, 3, 20, 0, time.UTC),
		}
		expectedPlistConfig = &v1feed.XProtectPlistConfigData{
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
	got := getTestData(sofafeed.MacOS)

	// Test top-level Feed fields
	if got.UpdateHash != expectedUpdateHash {
		t.Errorf("UpdateHash = %v, want %v", got.UpdateHash, expectedUpdateHash)
	}

	// Test XProtectPayloads
	if *got.XProtectPayloads != *expectedXProtect {
		t.Errorf("XProtectPayloads = %v, want %v", got.XProtectPayloads, expectedXProtect)
	}

	if *got.XProtectPlistConfigData != *expectedPlistConfig {
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

	if got.Models == nil {
		t.Errorf("No models found")
	} else {
		for _, expected := range expectedModels {
			model, exists := (*got.Models)[expected.id]
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

func TestIOSParseFeed(t *testing.T) {
	var (
		expectedUpdateHash = "382363738da23412c3ea17a630f5b4815ed1728d24140fd34198363e5c973ad8"
		expectedOSVersions = []struct {
			version string
			latest  struct {
				productVersion string
				build          string
			}
		}{
			{
				version: "18",
				latest: struct {
					productVersion string
					build          string
				}{
					productVersion: "18.3",
					build:          "22D63",
				},
			},
			{
				version: "17",
				latest: struct {
					productVersion string
					build          string
				}{
					productVersion: "17.7.2",
					build:          "21H6221",
				},
			},
			{
				version: "16",
				latest: struct {
					productVersion string
					build          string
				}{
					productVersion: "16.7.10",
					build:          "20H6350",
				},
			},
		}
		expectedSecurityReleases = []struct {
			version          string
			uniqueCVEsCount  int
			activeExploitCVE string // One example CVE known to be actively exploited
		}{
			{
				version:          "18.3",
				uniqueCVEsCount:  28,
				activeExploitCVE: "CVE-2025-24085",
			},
			{
				version:          "18.2.1",
				uniqueCVEsCount:  0,
				activeExploitCVE: "",
			},
			{
				version:          "18.2",
				uniqueCVEsCount:  38,
				activeExploitCVE: "",
			},
		}
		expectedDevices = []string{
			"iPad11,1",
			"iPad11,2",
			"iPad11,3",
		}
	)

	got := getTestData(sofafeed.IOS)

	if got.UpdateHash != expectedUpdateHash {
		t.Errorf("UpdateHash = %v, want %v", got.UpdateHash, expectedUpdateHash)
	}

	// Test OS Versions - we'll check iOS 18, 17, and 16

	if len(got.OSVersions) == 0 {
		t.Fatal("OSVersions is empty")
	}

	for _, expected := range expectedOSVersions {
		found := false
		for _, osv := range got.OSVersions {
			if osv.OSVersion == expected.version {
				found = true
				if osv.Latest.ProductVersion != expected.latest.productVersion {
					t.Errorf("OS Version %s product version = %v, want %v",
						expected.version, osv.Latest.ProductVersion, expected.latest.productVersion)
				}
				if osv.Latest.Build != expected.latest.build {
					t.Errorf("OS Version %s build = %v, want %v",
						expected.version, osv.Latest.Build, expected.latest.build)
				}
			}
		}
		if !found {
			t.Errorf("OS Version %s not found", expected.version)
		}
	}

	// Test Security Releases for iOS 18
	var ios18 *v1feed.OSVersion
	for i := range got.OSVersions {
		if got.OSVersions[i].OSVersion == "18" {
			ios18 = &got.OSVersions[i]
			break
		}
	}

	if ios18 == nil {
		t.Fatal("iOS 18 version not found")
	}

	for _, expected := range expectedSecurityReleases {
		found := false
		for _, release := range ios18.SecurityReleases {
			if release.ProductVersion == expected.version {
				found = true
				if release.UniqueCVEsCount != expected.uniqueCVEsCount {
					t.Errorf("iOS %s UniqueCVEsCount = %v, want %v",
						expected.version, release.UniqueCVEsCount, expected.uniqueCVEsCount)
				}
				if expected.activeExploitCVE != "" {
					foundExploit := false
					for _, cve := range release.ActivelyExploitedCVEs {
						if cve == expected.activeExploitCVE {
							foundExploit = true
							break
						}
					}
					if !foundExploit {
						t.Errorf("iOS %s missing expected actively exploited CVE %s",
							expected.version, expected.activeExploitCVE)
					}
				}
				break
			}
		}
		if !found {
			t.Errorf("Security release %s not found", expected.version)
		}
	}

	// Test Supported Devices for iOS 18
	for _, expectedDevice := range expectedDevices {
		found := false
		for _, device := range ios18.Latest.SupportedDevices {
			if device == expectedDevice {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Device %s not found in supported devices", expectedDevice)
		}
	}
}

// TestCVEMapping validates the parsing of CVE data from various OS versions
func TestCVEMapping(t *testing.T) {
	// Test cases with known CVEs from different OS versions and releases
	testCases := []struct {
		feedType   sofafeed.FeedType
		osVersion  string
		cve        string
		exploited  bool
		inLatest   bool   // Whether to look in Latest or SecurityReleases
		releaseVer string // Which release version to find it in, if not Latest
	}{
		{
			feedType:  sofafeed.MacOS,
			osVersion: "Sequoia 15",
			cve:       "CVE-2025-24085",
			exploited: true,
			inLatest:  true,
		},
		{
			feedType:   sofafeed.MacOS,
			osVersion:  "Sonoma 14",
			cve:        "CVE-2024-23225",
			exploited:  true,
			inLatest:   false,
			releaseVer: "14.4",
		},
		{
			feedType:  sofafeed.MacOS,
			osVersion: "Monterey 12",
			cve:       "CVE-2024-40783",
			exploited: false,
			inLatest:  true,
		},
		{
			feedType:   sofafeed.IOS,
			osVersion:  "18",
			cve:        "CVE-2025-24085",
			exploited:  true,
			inLatest:   false,
			releaseVer: "18.3",
		},
		{
			feedType:   sofafeed.IOS,
			osVersion:  "17",
			cve:        "CVE-2024-44308",
			exploited:  true,
			inLatest:   false,
			releaseVer: "17.7.2",
		},
	}
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%s/%s/%s", tc.feedType, tc.osVersion, tc.cve), func(t *testing.T) {
			testData := getTestData(tc.feedType)
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
		feedType      sofafeed.FeedType
		osVersion     string
		deviceID      string
		shouldSupport bool
	}{
		{sofafeed.MacOS, "Sequoia 15", "Mac-1E7E29AD0135F9BC", true},  // Known supported device
		{sofafeed.MacOS, "Sonoma 14", "Mac-937A206F2EE63C01", true},   // Known supported device
		{sofafeed.MacOS, "Monterey 12", "Mac-FFE5EF870D7BA81A", true}, // Known supported device
		{sofafeed.MacOS, "Sequoia 15", "Mac-INVALID-ID", false},       // Invalid device
		{sofafeed.MacOS, "Sonoma 14", "MacBookPro11,1", false},        // Unsupported device
		{sofafeed.MacOS, "Monterey 12", "iMac14,1", false},            // Unsupported device
		{sofafeed.IOS, "18", "iPad11,1", true},                        // Current-gen iPad mini, officially supported
		{sofafeed.IOS, "18", "iPad11,2", true},                        // iPad mini variant with cellular, supported
		{sofafeed.IOS, "18", "iPad99,99", false},                      // Non-existent device ID, should not be supported
	}
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%s/%s/%s", tc.feedType, tc.osVersion, tc.deviceID), func(t *testing.T) {
			testData := getTestData(tc.feedType)
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
			feed := v1feed.NewFeed()
			_, err := feed.ParseString(tc.json)
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

	// Test with nil client (should use default)
	feed := v1feed.NewFeed()
	ctx := context.Background()
	data, err := feed.Fetch(ctx, &feeds.FetchArgs{
		FeedURL: ts.URL,
	})
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

func getTestData(ft sofafeed.FeedType) *v1feed.Feed {
	return testData[ft].ParseResult().Result().(*v1feed.Feed)
}

// testData holds the parsed feed that all tests can use
var testData = make(map[sofafeed.FeedType]feeds.Feed, 2)

// init loads the test data once for all tests
func init() {
	loadTestData(sofafeed.IOS)
	loadTestData(sofafeed.MacOS)
}

func loadTestData(ft sofafeed.FeedType) {
	j, err := os.ReadFile(fmt.Sprintf("testdata/sofafeed-%s-v1.json", ft))
	if err != nil {
		panic(fmt.Sprintf("Failed to read '%s' test data: %s", ft, err.Error()))
	}
	var feed feeds.Feed
	feed, err = sofafeed.Parse(ft, j)
	if err != nil {
		panic(fmt.Sprintf("Failed to parse '%s' test data: %s", ft, err.Error()))
	}
	testData[ft] = feed
}
