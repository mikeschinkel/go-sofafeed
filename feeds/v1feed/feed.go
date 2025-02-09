package v1feed

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/mikeschinkel/go-sofafeed/feeds"
)

// Feed represents the top-level structure for both macOS and iOS update feeds
type Feed struct {
	//-------------------------------------------------------
	// These two are common to both macOS and iOS feeds in v1
	//-------------------------------------------------------

	// UpdateHash is a unique identifier for the current feed state
	UpdateHash string `json:"UpdateHash"`
	// OSVersions contains information about different major OS versions, e.g. 15, 14, etc.
	OSVersions []OSVersion `json:"OSVersions"`

	//--------------------------------------------
	// These four are specific to only macOS in v1
	//--------------------------------------------

	// XProtectPayloads contains macOS-specific XProtect framework update information
	XProtectPayloads *XProtectPayloads `json:"XProtectPayloads,omitempty"`
	// XProtectPlistConfigData contains macOS-specific XProtect plist configuration data
	XProtectPlistConfigData *XProtectPlistConfigData `json:"XProtectPlistConfigData,omitempty"`
	// Models maps device identifiers to their capabilities (macOS-specific)
	Models *Models `json:"Models,omitempty"`
	// InstallationApps contains macOS-specific installer application information
	InstallationApps *InstallationApps `json:"InstallationApps,omitempty"`
}

func NewFeed() *Feed {
	return &Feed{}
}

// Parse parses JSON feed data from bytes
func (f *Feed) Parse(data []byte) (pr feeds.ParseResult, err error) {
	var feed Feed
	err = json.Unmarshal(data, &feed)
	if err != nil {
		err = errors.Join(feeds.ErrParseFeed, err)
		goto end
	}
	pr = &parseResult{result: &feed}
end:
	return pr, err
}

// ParseString parses JSON feed data from a string
func (f *Feed) ParseString(data string) (feeds.ParseResult, error) {
	return f.Parse([]byte(data))
}

// FetchAndParse retrieves and parses the latest feed
func (f *Feed) FetchAndParse(ctx context.Context, args *feeds.FetchArgs) (pr feeds.ParseResult, err error) {
	var data []byte
	data, err = f.Fetch(ctx, args)
	if err != nil {
		goto end
	}

	pr, err = f.Parse(data)
	if err != nil {
		goto end
	}

end:
	return pr, err
}

// Fetch retrieves the latest JSON from sofafeed.Endpoint. If client is nil, a
// default client with a 30-second timeout will be used.
func (f *Feed) Fetch(ctx context.Context, args *feeds.FetchArgs) (body []byte, err error) {
	var req *http.Request
	var resp *http.Response
	// Create a default client if none provided
	if args.Client == nil {
		args.Client = &http.Client{
			Timeout: feeds.Timeout,
		}
	}

	// Create the request with the provided context
	req, err = http.NewRequestWithContext(ctx, http.MethodGet, args.FeedURL, nil)
	if err != nil {
		err = errors.Join(feeds.ErrFetchRequest, err)
		goto end
	}

	// Set headers for better HTTP citizenship
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", feeds.UserAgent)

	// Perform the request
	resp, err = args.Client.Do(req)
	if err != nil {
		err = errors.Join(feeds.ErrPerformRequest, err)
		goto end
	}
	defer mustClose(resp.Body)

	// Check for non-200 status codes
	if resp.StatusCode != http.StatusOK {
		err = errors.Join(
			feeds.ErrUnexpectedStatusCode,
			errors.New(fmt.Sprintf("status_code=%d", resp.StatusCode)),
		)
		goto end
	}

	// Read the entire response body
	body, err = io.ReadAll(resp.Body)
	if err != nil {
		err = errors.Join(feeds.ErrReadResponseBody, err)
		goto end
	}
end:
	return body, err
}
