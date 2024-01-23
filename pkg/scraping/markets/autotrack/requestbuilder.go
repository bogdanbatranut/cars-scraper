package autotrack

import (
	"carscraper/pkg/jobs"
)

type Request struct {
	urlBuilder *URLBuilder
}

func NewRequest(criteria jobs.Criteria) *Request {
	builder := NewURLBuilder(criteria)
	return &Request{
		urlBuilder: builder,
	}
}
