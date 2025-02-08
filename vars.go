package sofafeed

import (
	"time"
)

var (
	// Endpoint is the URL for the macadmins.io JSON feed
	Endpoint = "https://sofafeed.macadmins.io/v1/macos_data_feed.json"

	// Timeout is the recommended timeout for HTTP requests
	Timeout = 30 * time.Second
)
