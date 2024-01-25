package olx

import (
	"io"
	"log"
	"net/http"
)

type Request struct {
}

func NewRequest() *Request {
	return &Request{}
}

func (r Request) GetPage(url string) ([]byte, error) {

	httpMethod := "GET"
	httpClient := &http.Client{}
	httpRequest, err := http.NewRequest(httpMethod, url, nil)

	if err != nil {
		panic(err)
	}

	response, err := httpClient.Do(httpRequest)
	//log.Println("Status code : ", response.StatusCode)
	if err != nil {
		log.Printf("got response with error: %+v", err)
		return nil, err
	}
	defer response.Body.Close()
	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {

		return nil, err
	}

	return bodyBytes, nil
}
