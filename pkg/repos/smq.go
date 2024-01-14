package repos

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
)

type IMessageQueue interface {
	GetMessageWithDelete(topic string) *[]byte
	GetMessage(topic string) *[]byte
	PutMessage(topic string, message []byte)
}

func NewSimpleMessageQueueRepository(baseURL string) *SimpleMessageQueueRepository {
	return &SimpleMessageQueueRepository{
		baseURL: baseURL,
	}
}

type SimpleMessageQueueRepository struct {
	baseURL string
}

func (smqr SimpleMessageQueueRepository) peekURL(topic string) string {
	return fmt.Sprintf("%s/peek/%s", smqr.baseURL, topic)
}

func (smqr SimpleMessageQueueRepository) popURL(topic string) string {
	return fmt.Sprintf("%s/pop/%s", smqr.baseURL, topic)
}

func (smqr SimpleMessageQueueRepository) pushURL(topic string) string {
	return fmt.Sprintf("%s/push/%s", smqr.baseURL, topic)
}

func (s SimpleMessageQueueRepository) GetMessageWithDelete(topic string) *[]byte {
	log.Printf("Get message from %s", s.popURL(topic))
	resp, err := http.Get(s.popURL(topic))
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("Response code: ", resp.StatusCode)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	return &body
}

func (s SimpleMessageQueueRepository) GetMessage(topic string) *[]byte {
	log.Println("Peek message")
	resp, err := http.Get(s.peekURL(topic))
	if err != nil {
		log.Fatalln(err)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	return &body
}

func (s SimpleMessageQueueRepository) PutMessage(topic string, message []byte) {
	// Create a HTTP post request
	r, err := http.NewRequest("POST", s.pushURL(topic), bytes.NewBuffer(message))
	if err != nil {
		panic(err)
	}
	client := &http.Client{}
	_, err = client.Do(r)
	if err != nil {
		panic(err)
	}
	//log.Printf("Push Message to topic: %s response code: %d", topic, res.StatusCode)
}

type MessageQueueConfig struct {
}
