package autovit

import (
	"carscraper/pkg/jobs"
	"io"
	"log"
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

func (r Request) GetPage(pageNumber int) ([]byte, string, error) {

	url := r.urlBuilder.GetPageURL(pageNumber)

	httpMethod := "GET"
	httpClient := &http.Client{}
	httpRequest, err := http.NewRequest(httpMethod, url, nil)

	if err != nil {
		return nil, url, err
		panic(err)
	}

	response, err := httpClient.Do(httpRequest)
	log.Println("Status code : ", response.StatusCode)
	if err != nil {
		log.Printf("got response with error: %+v", err)
		return nil, url, err
	}
	defer response.Body.Close()
	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {

		return nil, url, err
	}

	return bodyBytes, url, nil
}

func (r Request) DoRequest(url string) ([]byte, error) {

	httpMethod := "GET"
	httpClient := &http.Client{}
	httpRequest, err := http.NewRequest(httpMethod, url, nil)
	httpRequest.Header.Add("Host", "mydomain.com")
	httpRequest.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/128.0.0.0 Safari/537.36")

	if err != nil {
		return nil, err
	}

	response, err := httpClient.Do(httpRequest)
	log.Println("Status code : ", response.StatusCode)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return bodyBytes, nil
}
