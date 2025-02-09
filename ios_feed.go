package sofafeed

import (
	"context"

	"github.com/mikeschinkel/go-sofafeed/feeds"
	"github.com/mikeschinkel/go-sofafeed/feeds/v1feed"
)

type IOSFeed struct {
	v1feed.Feed
	feeds.FeedCommon
}

func NewIOSFeed() *IOSFeed {
	common := feeds.FeedCommon{}
	common.SetURL("https://sofafeed.iadmins.io/v1/ios_data_feed.json")
	return &IOSFeed{
		FeedCommon: common,
	}
}

func ParseIOSFeed(data []byte) (feed *IOSFeed, err error) {
	return parseIOSFeed(func() (feeds.Feed, error) {
		return Parse(IOS, data)
	})
}

func FetchAndParseIOSFeed(ctx context.Context, args *feeds.FetchArgs) (feed *IOSFeed, err error) {
	return parseIOSFeed(func() (feeds.Feed, error) {
		return FetchAndParse(ctx, IOS, args)
	})
}

func parseIOSFeed(pf func() (feeds.Feed, error)) (feed *IOSFeed, err error) {
	var f feeds.Feed
	f, err = pf()
	if err != nil {
		goto end
	}
	feed = f.ParseResult().Result().(*IOSFeed)
end:
	return feed, err
}
