package v1feed

type parseResult struct {
	result *Feed
}

func (pr *parseResult) Result() any {
	return pr.result
}
