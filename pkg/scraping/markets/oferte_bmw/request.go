package oferte_bmw

import (
	http2 "carscraper/pkg/http"
	"encoding/json"
)

type Request struct {
	request *http2.Request
}

func NewRequest() *Request {
	r := http2.NewRequest()
	return &Request{request: r}
}

func (r Request) DoPOSTRequest(url string, body RequestBody) (*OferteBMWResponse, error) {
	bodyBytesArr, err := json.Marshal(&body)
	if err != nil {
		return nil, err
	}
	responseBytes, err := r.request.Post(url, bodyBytesArr)
	if err != nil {
		return nil, err
	}
	var res OferteBMWResponse
	err = json.Unmarshal(responseBytes, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}
