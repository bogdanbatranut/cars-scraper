package autoklass

import (
	http2 "carscraper/pkg/http"
	"encoding/json"
	"log"
)

type Request struct {
	request *http2.Request
}

func NewRequest() *Request {
	r := http2.NewRequest()
	return &Request{request: r}
}

func (r Request) MakeRequest(url string) (*Response, error) {
	log.Println("---------------------")
	log.Println(url)
	responseBytes, err := r.request.GET(url)
	if err != nil {
		return nil, err
	}
	var res Response
	err = json.Unmarshal(responseBytes, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}
