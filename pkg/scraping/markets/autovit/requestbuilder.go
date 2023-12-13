package autovit

import (
	"carscraper/pkg/jobs"
	"io"
	"net/http"
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

func (r Request) GetPage(pageNumber int) []byte {

	url := r.urlBuilder.GetPageURL(pageNumber)

	httpMethod := "GET"
	httpClient := &http.Client{}
	httpRequest, err := http.NewRequest(httpMethod, url, nil)

	if err != nil {
		panic(err)
	}

	response, err := httpClient.Do(httpRequest)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()
	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}

	return bodyBytes
}
