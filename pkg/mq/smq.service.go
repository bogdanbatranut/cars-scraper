package mq

import (
	"carscraper/pkg/amconfig"
	"carscraper/pkg/jobs"
	"carscraper/pkg/repos"
	"encoding/json"
	"fmt"
	"log"
)

type MessageQueueService struct {
	messagesQueue    repos.IMessageQueue
	resultsTopicName string
	jobsTopicName    string
	adsChannel       chan jobs.AdsPageJobResult
	cfg              amconfig.IConfig
}

type MessageQueueServiceConfiguration func(mqs *MessageQueueService)

func NewMessageQueueService(cfg amconfig.IConfig, cfgs ...MessageQueueServiceConfiguration) *MessageQueueService {

	mqs := &MessageQueueService{
		cfg:        cfg,
		adsChannel: make(chan jobs.AdsPageJobResult),
	}
	for _, cfg := range cfgs {
		cfg(mqs)
	}
	return mqs
}

func (service MessageQueueService) PublishResults(adsResults jobs.AdsPageJobResult) {
	service.adsChannel <- adsResults
}

func WithMessageQueueRepository(mqr repos.IMessageQueue) MessageQueueServiceConfiguration {
	return func(mqs *MessageQueueService) {
		mqs.messagesQueue = mqr
	}
}

func WithLocalMessageQueue() MessageQueueServiceConfiguration {
	smqr := repos.NewSimpleMessageQueueRepository(
		"http://127.0.0.1:3333")
	return WithMessageQueueRepository(smqr)
}

func WithProdMessageQueue() MessageQueueServiceConfiguration {
	return func(service *MessageQueueService) {
		smqHost := service.cfg.GetString(amconfig.SMQURL)
		smqPort := service.cfg.GetString(amconfig.SMQHTTPPort)
		resultsTopicName := service.cfg.GetString(amconfig.SMQResultsTopicName)
		jobsTopicName := service.cfg.GetString(amconfig.SMQJobsTopicName)
		service.jobsTopicName = jobsTopicName
		service.resultsTopicName = resultsTopicName
		//smqr := repos.NewSimpleMessageQueueRepository(fmt.Sprintf("http://%s:%s", smqHost, smqPort))
		//return WithMessageQueueRepository(smqr)

		service.messagesQueue = repos.NewQueueRepository(fmt.Sprintf("http://%s:%s", smqHost, smqPort))
	}
}

func (service MessageQueueService) Start() {
	go func() {
		for {
			service.handleResults()
		}
	}()
}

func (service MessageQueueService) handleResults() {
	res := <-service.adsChannel
	resBytes, err := json.Marshal(&res)
	if err != nil {
		panic(err)
	}
	if service.messagesQueue == nil {
		log.Println("No Message queue")
		return
	}
	service.messagesQueue.PutMessage(service.resultsTopicName, resBytes)
}
