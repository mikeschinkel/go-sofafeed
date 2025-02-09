# go-sofafeed

A Go library for parsing and interacting with Apple software update feeds from sofa.macadmins.io. This package provides a structured way to access information about macOS and iOS updates, including security releases, CVE information, and device compatibility.

## Features

The go-sofafeed library offers comprehensive support for both macOS and iOS update feeds with the following capabilities:

- Parse update feed data from JSON
- Fetch feed data directly from sofa.macadmins.io
- Access structured information about:
  - OS versions and builds
  - Security releases and CVEs
  - Device compatibility
  - XProtect updates (macOS)
  - Universal macOS Applications (macOS)

## Installation

To install go-sofafeed, use `go get`:

```bash
go get github.com/mikeschinkel/go-sofafeed
```

## Usage

Here are some common usage examples:

### Fetching and Parsing Feed Data

```go
package main

import (
    "context"
    "fmt"
    "github.com/mikeschinkel/go-sofafeed"
    "github.com/mikeschinkel/go-sofafeed/feeds"
)

func main() {
    // Create context and fetch args
    ctx := context.Background()
    args := &feeds.FetchArgs{}

    // Fetch and parse macOS feed
    macFeed, err := sofafeed.FetchAndParseMacOSFeed(ctx, args)
    if err != nil {
        panic(err)
    }

    // Fetch and parse iOS feed
    iosFeed, err := sofafeed.FetchAndParseIOSFeed(ctx, args)
    if err != nil {
        panic(err)
    }

    // Access feed data
    fmt.Printf("Latest macOS: %s\n", macFeed.OSVersions[0].Latest.ProductVersion)
    fmt.Printf("Latest iOS: %s\n", iosFeed.OSVersions[0].Latest.ProductVersion)
}
```
### Parsing Local Feed Data

```go
package main

import (
    "fmt"
    "os"
    "github.com/mikeschinkel/go-sofafeed"
)

func main() {
    // Read local JSON file
    data, err := os.ReadFile("macos_feed.json")
    if err != nil {
        panic(err)
    }

    // Parse macOS feed
    feed, err := sofafeed.ParseMacOSFeed(data)
    if err != nil {
        panic(err)
    }

    // Access feed information
    for _, osv := range feed.OSVersions {
        fmt.Printf("OS Version: %s\n", osv.OSVersion)
        fmt.Printf("Latest Build: %s\n", osv.Latest.Build)
        fmt.Printf("CVE Count: %d\n", osv.Latest.UniqueCVEsCount)
    }
}
```
### Working with Security Information

```go
package main

import (
    "context"
    "fmt"
    "github.com/mikeschinkel/go-sofafeed"
    "github.com/mikeschinkel/go-sofafeed/feeds"
)

func main() {
    ctx := context.Background()
    args := &feeds.FetchArgs{}

    feed, err := sofafeed.FetchAndParseMacOSFeed(ctx, args)
    if err != nil {
        panic(err)
    }

    // Examine security releases
    for _, osv := range feed.OSVersions {
        for _, release := range osv.SecurityReleases {
            fmt.Printf("Update: %s\n", release.UpdateName)
            fmt.Printf("CVEs: %d\n", release.UniqueCVEsCount)
            
            if len(release.ActivelyExploitedCVEs) > 0 {
                fmt.Printf("Active Exploits: %v\n", release.ActivelyExploitedCVEs)
            }
        }
    }
}
```
### Customizing the HTTP Client

You can provide your own `http.Client` to control timeouts, transport settings, or add custom behaviors:

```go
package main

import (
    "context"
    "crypto/tls"
    "net/http"
    "time"
    "github.com/mikeschinkel/go-sofafeed"
    "github.com/mikeschinkel/go-sofafeed/feeds"
)

func main() {
    // Create a custom HTTP client with specific settings
    client := &http.Client{
        Timeout: 45 * time.Second,
        Transport: &http.Transport{
            TLSClientConfig: &tls.Config{
                MinVersion: tls.VersionTLS12,
            },
            MaxIdleConns:        100,
            MaxIdleConnsPerHost: 100,
            IdleConnTimeout:     90 * time.Second,
        },
    }

    // Pass the client via FetchArgs
    args := &feeds.FetchArgs{
        Client: client,
    }

    // Use the custom client for fetching
    feed, err := sofafeed.FetchAndParseMacOSFeed(context.Background(), args)
    if err != nil {
        panic(err)
    }
   fmt.Printf("Latest macOS: %s\n", feed.OSVersions[0].Latest.ProductVersion)
}
```
### Using Custom Feed URLs

For enterprise environments, you might want to use locally cached or security-approved feed URLs. You can set custom URLs when creating feed instances:

```go
package main

import (
   "context"
   "fmt"
   "github.com/mikeschinkel/go-sofafeed"
   "github.com/mikeschinkel/go-sofafeed/feeds"
)

func main() {
   // Create context and fetch arguments with custom URL
   ctx := context.Background()
   args := &feeds.FetchArgs{
     FeedURL: "https://internal-cache.company.com/approved-feeds/macos_data_feed.json",
   }
   
   // Fetch macOS feed using custom URL
   macFeed, err := sofafeed.FetchAndParseMacOSFeed(ctx, args)
   if err != nil {
     panic(err)
   }
   
   // Similarly for iOS feed
   args.FeedURL = "https://internal-cache.company.com/approved-feeds/ios_data_feed.json"
   iosFeed, err := sofafeed.FetchAndParseIOSFeed(ctx, args)
   if err != nil {
     panic(err)
   }
   
   // Access feed data as needed
   fmt.Printf("Latest macOS: %s\n", macFeed.OSVersions[0].Latest.ProductVersion)
   fmt.Printf("Latest iOS: %s\n", iosFeed.OSVersions[0].Latest.ProductVersion)
}
```

This approach is particularly useful when:
- You need to cache feeds locally to reduce external requests
- Your security policy requires feeds to be reviewed before use
- You're operating in an air-gapped environment
- You want to add rate limiting or monitoring at your proxy

## Feed Structure

### Common Fields

Both macOS and iOS feeds share these fields:

- `UpdateHash`: Unique identifier for the current feed state
- `OSVersions`: Array of OS version information including:
  - Version details
  - Security releases
  - Device compatibility
  - CVE information

### macOS-Specific Fields

The macOS feed includes additional fields:

- `XProtectPayloads`: XProtect framework update information
- `XProtectPlistConfigData`: XProtect configuration data
- `Models`: Device model compatibility information
- `InstallationApps`: Universal macOS Application information

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request. For major changes, please open an issue first to discuss what you would like to change.

Please make sure to update tests as appropriate.

## Testing

To run the tests:

```bash
go test ./...
```

## License

- [**MIT**](https://mit-license.org/) - see LICENSE file for details

## Credits

This project utilizes the [sofa.macadmins.io](https://sofa.macadmins.io) service for Apple software update feeds.