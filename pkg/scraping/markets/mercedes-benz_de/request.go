package mercedes_benz_de

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

func (r Request) MakeRequest(url string, search Search) (*Response, error) {
	bodyBytesArr, err := json.Marshal(&search)
	if err != nil {
		return nil, err
	}
	log.Println("---------------------")
	log.Println(url)
	log.Println(string(bodyBytesArr))
	responseBytes, err := r.request.Post(url, bodyBytesArr)
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
