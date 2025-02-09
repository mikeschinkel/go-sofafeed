package feeds

import (
	"context"
	"net/http"
)

type FetchArgs struct {
	Client  *http.Client
	FeedURL string
}
type ParseResult interface {
	Result() any
}

type Feed interface {
	SetURL(string)
	URL() string
	Fetch(context.Context, *FetchArgs) ([]byte, error)
	Parse([]byte) (ParseResult, error)
	ParseString(string) (ParseResult, error)
	FetchAndParse(context.Context, *FetchArgs) (ParseResult, error)
	SetParseResult(ParseResult)
	ParseResult() ParseResult
}

type FeedCommon struct {
	url    string
	parsed ParseResult
}

func (f *FeedCommon) URL() string {
	return f.url
}

func (f *FeedCommon) SetURL(url string) {
	f.url = url
}

func (f *FeedCommon) SetParseResult(parsed ParseResult) {
	f.parsed = parsed
}

func (f *FeedCommon) ParseResult() ParseResult {
	return f.parsed
}
