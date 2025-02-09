package sofafeed

import (
	"context"
	"errors"
	"fmt"

	"github.com/mikeschinkel/go-sofafeed/feeds"
)

// Parse parses JSON feed data from bytes
func Parse(ft FeedType, data []byte) (feeds.Feed, error) {
	var result feeds.ParseResult
	feed, err := getFeed(ft)
	if err != nil {
		goto end
	}
	result, err = feed.Parse(data)
	if err != nil {
		goto end
	}
	feed.SetParseResult(result)
end:
	return feed, err
}

// ParseString parses JSON feed data from a string
func ParseString(ft FeedType, data string) (feeds.Feed, error) {
	return Parse(ft, []byte(data))
}

// FetchAndParse retrieves and parses the latest feed
func FetchAndParse(ctx context.Context, ft FeedType, args *feeds.FetchArgs) (feed feeds.Feed, err error) {
	var data []byte
	var result feeds.ParseResult

	feed, err = getFeed(ft)
	if err != nil {
		goto end
	}
	data, err = feed.Fetch(ctx, args)
	if err != nil {
		goto end
	}
	result, err = feed.Parse(data)
	if err != nil {
		goto end
	}
	feed.SetParseResult(result)
end:
	return feed, err
}

// Fetch retrieves the latest JSON from sofafeed.Endpoint. If client is nil, a
// default client with a 30-second timeout will be used.
func Fetch(ctx context.Context, ft FeedType, args *feeds.FetchArgs) (body []byte, err error) {
	feed, err := getFeed(ft)
	if err != nil {
		goto end
	}
	args.FeedURL = feed.URL()
	body, err = feed.Fetch(ctx, args)
end:
	return body, err
}

func getFeed(ft FeedType) (feed feeds.Feed, err error) {
	switch ft {
	case IOS:
		feed = NewIOSFeed()
	case MacOS:
		feed = NewMacOSFeed()
	default:
		err = errors.Join(
			feeds.ErrCheckFeedType,
			fmt.Errorf("%s=%s", feeds.FeedTypeErrArg, ft),
		)
	}
	return feed, err
}
