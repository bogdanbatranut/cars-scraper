package mq

import (
	"carscraper/pkg/amconfig"
	"carscraper/pkg/repos"
	"encoding/json"
	"log"
)

type MessageQueueService struct {
	messagesQueue       repos.IMessageQueue
	criteriasRepository repos.CriteriaRepository
}

func NewMessageQueueService(cfgs ...MessageQueueServiceConfiguration) *MessageQueueService {

	mqs := &MessageQueueService{}
	for _, cfg := range cfgs {
		cfg(mqs)
	}
	return mqs
}

type MessageQueueServiceConfiguration func(mqs *MessageQueueService)

func WithMessageQueueRepository(mqr repos.IMessageQueue) MessageQueueServiceConfiguration {
	return func(mqs *MessageQueueService) {
		mqs.messagesQueue = mqr
	}
}

func WithSimpleMessageQueueRepository() MessageQueueServiceConfiguration {
	smqr := repos.NewSimpleMessageQueueRepository(
		"http://127.0.0.1:3333")
	return WithMessageQueueRepository(smqr)
}

func WithCriteriaSQLRepository(cfg amconfig.IConfig) MessageQueueServiceConfiguration {
	return func(mqs *MessageQueueService) {
		mqs.criteriasRepository = repos.NewSQLCriteriaRepository(cfg)
	}
}

func (mqs MessageQueueService) Start() {
	criterias := mqs.criteriasRepository.GetAll()
	for _, criteria := range *criterias {
		bytes, err := json.Marshal(&criteria)
		if err != nil {
			panic(err)
		}
		mqs.messagesQueue.PutMessage("test", bytes)
	}

	for {

		resp := mqs.messagesQueue.GetMessage("test")
		log.Printf(string(*resp))
		if len(*resp) == 0 {
			break
		}
	}
}
