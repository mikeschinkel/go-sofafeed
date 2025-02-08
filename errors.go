package sofafeed

import "errors"

var (
	// ErrParseFeed represents an error that occurs during JSON unmarshaling of the feed
	ErrParseFeed = errors.New("failed to parse feed from JSON")

	// ErrFetchRequest represents an error that occurs when creating an HTTP request
	ErrFetchRequest = errors.New("failed to create fetch request")

	// ErrPerformRequest represents an error that occurs during the HTTP request execution
	ErrPerformRequest = errors.New("failed to perform HTTP request")

	// ErrUnexpectedStatusCode represents an error for non-200 HTTP status codes
	ErrUnexpectedStatusCode = errors.New("unexpected HTTP status code")

	// ErrReadResponseBody represents an error when reading the response body
	ErrReadResponseBody = errors.New("failed to read response body")
)
