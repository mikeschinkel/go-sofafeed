package sofafeed

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/mikeschinkel/go-sofafeed/feeds/v1feed"
)

// Parse parses JSON feed data from bytes
func Parse(data []byte) (*v1feed.Feed, error) {
	var feed v1feed.Feed
	err := json.Unmarshal(data, &feed)
	if err != nil {
		err = errors.Join(ErrParseFeed, err)
		goto end
	}
end:
	return &feed, err
}

// ParseString parses JSON feed data from a string
func ParseString(data string) (*v1feed.Feed, error) {
	return Parse([]byte(data))
}

// FetchAndParse retrieves and parses the latest feed
func FetchAndParse(ctx context.Context, client *http.Client) (feed *v1feed.Feed, err error) {
	var data []byte
	data, err = Fetch(ctx, client)
	if err != nil {
		goto end
	}

	feed, err = Parse(data)
	if err != nil {
		goto end
	}

end:
	return feed, err
}

// Fetch retrieves the latest JSON from sofafeed.Endpoint. If client is nil, a
// default client with a 30-second timeout will be used.
func Fetch(ctx context.Context, client *http.Client) (body []byte, err error) {
	var req *http.Request
	var resp *http.Response
	// Create a default client if none provided
	if client == nil {
		client = &http.Client{
			Timeout: Timeout,
		}
	}

	// Create the request with the provided context
	req, err = http.NewRequestWithContext(ctx, http.MethodGet, Endpoint, nil)
	if err != nil {
		err = errors.Join(ErrFetchRequest, err)
		goto end
	}

	// Set headers for better HTTP citizenship
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "go-sofafeed/1.0")

	// Perform the request
	resp, err = client.Do(req)
	if err != nil {
		err = errors.Join(ErrPerformRequest, err)
		goto end
	}
	defer mustClose(resp.Body)

	// Check for non-200 status codes
	if resp.StatusCode != http.StatusOK {
		err = errors.Join(
			ErrUnexpectedStatusCode,
			errors.New(fmt.Sprintf("status_code=%d", resp.StatusCode)),
		)
		goto end
	}

	// Read the entire response body
	body, err = io.ReadAll(resp.Body)
	if err != nil {
		err = errors.Join(ErrReadResponseBody, err)
		goto end
	}
end:
	return body, err
}
