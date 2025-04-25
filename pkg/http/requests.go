package http

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

type Request struct {
}

func NewRequest() *Request {
	return &Request{}
}

func (r Request) Post(url string, bodyBytes []byte) ([]byte, error) {

	httpMethod := "POST"
	httpClient := &http.Client{}

	httpRequest, err := http.NewRequest(httpMethod, url, bytes.NewReader(bodyBytes))
	httpRequest.Header.Add("Content-Type", "application/json")

	if err != nil {
		panic(err)
	}

	response, err := httpClient.Do(httpRequest)
	//log.Println("Status code : ", response.StatusCode)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	responseBodyBytes, err := io.ReadAll(response.Body)
	if err != nil {

		return nil, err
	}

	return responseBodyBytes, nil
}

func (r Request) GET(url string) ([]byte, error) {
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return body, nil
}
