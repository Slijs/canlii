package canlii

import (
	"net/http"
	"net/url"
)

type search struct {
	client *Client
}

type SearchOptions struct {
	Offset      int
	ResultCount int
}

type SearchResult struct {
	TotalResults int
	Cases        []CaseListItem
	Legislations []Legislation
}

func (s *search) Search(fulltext string, opts *SearchOptions) (SearchResult, *http.Response, error) {
	const maxResultCount = 100
	if opts == nil {
		opts = &SearchOptions{
			ResultCount: maxResultCount,
		}
	}
	if opts.ResultCount > maxResultCount {
		opts.ResultCount = maxResultCount
	}

	q := make(url.Values)
	setIntParam(q, "offset", opts.Offset)
	setIntParam(q, "resultCount", opts.ResultCount)
	q.Set("fulltext", fulltext)

	res := SearchResult{}

	var v struct {
		ResultCount int `json:"resultCount"`
		Results     []struct {
			Case        CaseListItem `json:"case"`
			Legislation Legislation  `json:"legislation"`
		} `json:"results"`
	}
	resp, err := s.client.Get("search", "", q, &v)
	if err != nil {
		return res, resp, err
	}

	res.TotalResults = v.ResultCount

	for _, r := range v.Results {
		switch {
		case r.Legislation.DatabaseID != "":
			res.Legislations = append(res.Legislations, r.Legislation)
		case r.Case.DatabaseID != "":
			res.Cases = append(res.Cases, r.Case)
		}
	}

	return res, resp, err
}