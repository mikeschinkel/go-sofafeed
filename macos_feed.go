package sofafeed

import (
	"context"

	"github.com/mikeschinkel/go-sofafeed/feeds"
	"github.com/mikeschinkel/go-sofafeed/feeds/v1feed"
)

type MacOSFeed struct {
	v1feed.Feed
	feeds.FeedCommon
}

func NewMacOSFeed() *MacOSFeed {
	common := feeds.FeedCommon{}
	common.SetURL("https://sofafeed.macadmins.io/v1/macos_data_feed.json")
	return &MacOSFeed{
		FeedCommon: common,
	}
}

func ParseMacOSFeed(data []byte) (feed *MacOSFeed, err error) {
	return parseMacOSFeed(func() (feeds.Feed, error) {
		return Parse(MacOS, data)
	})
}

func FetchAndParseMacOSFeed(ctx context.Context, args *feeds.FetchArgs) (feed *MacOSFeed, err error) {
	return parseMacOSFeed(func() (feeds.Feed, error) {
		return FetchAndParse(ctx, MacOS, args)
	})
}

func parseMacOSFeed(pf func() (feeds.Feed, error)) (feed *MacOSFeed, err error) {
	var f feeds.Feed
	f, err = pf()
	if err != nil {
		goto end
	}
	feed = f.ParseResult().Result().(*MacOSFeed)
end:
	return feed, err
}
