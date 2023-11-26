package sessionstarter

import (
	"carscraper/pkg/adsdb"
	"carscraper/pkg/config"
	"carscraper/pkg/jobs"
	"carscraper/pkg/repos"
	"carscraper/pkg/url"
	"encoding/json"
	"fmt"
	"log"

	"github.com/google/uuid"
)

type SessionStarterService struct {
	messagesQueue              repos.IMessageQueue
	criteriasRepository        repos.CriteriaRepository
	requesteScrapingJobs       []jobs.Session
	urlComposerImplementations *url.URLComposerImplementations
	criteriasTopicName         string
}

type CrawlinglInitiatorServiceConfiguration func(cisc *SessionStarterService)

func NewSessionStarterService(cfgs ...CrawlinglInitiatorServiceConfiguration) *SessionStarterService {

	mqs := &SessionStarterService{
		urlComposerImplementations: url.NewURLComposerImplementations(),
		criteriasTopicName:         "jobs",
	}
	for _, cfg := range cfgs {
		cfg(mqs)
	}
	return mqs
}

func WithCriteriaSQLRepository(cfg config.IConfig) CrawlinglInitiatorServiceConfiguration {
	return func(mqs *SessionStarterService) {
		mqs.criteriasRepository = repos.NewSQLCriteriaRepository(cfg)
	}
}

func WithSimpleMessageQueueRepository(cfg config.IConfig) CrawlinglInitiatorServiceConfiguration {
	smqBaseURL := cfg.GetString(config.SMQURL)
	smqPort := cfg.GetString(config.SMQHTTPPort)
	smqURL := fmt.Sprintf("http://%s:%s", smqBaseURL, smqPort)

	smqr := repos.NewSimpleMessageQueueRepository(smqURL)

	//smqr := repos.NewSimpleMessageQueueRepository(
	//	"http://host.docker.internal:3333",
	//)

	return WithMessageQueueRepository(smqr)
}

func WithMessageQueueRepository(mqr repos.IMessageQueue) CrawlinglInitiatorServiceConfiguration {
	return func(cis *SessionStarterService) {
		cis.messagesQueue = mqr
	}
}

func (sss SessionStarterService) Start() {
	// create and push jobs to queue
	session := sss.newSession()
	sss.pushSessionJobs(session.Jobs)
	log.Printf("Puhsed session %+v", session.SessionID)

}

func (sss SessionStarterService) newSession() jobs.Session {
	sessionID := uuid.New()
	sessionJobs := sss.createSessionJobs(sessionID)
	session := jobs.Session{
		SessionID: sessionID,
		Jobs:      sessionJobs,
	}
	return session

}

func (sss SessionStarterService) createSessionJobs(sessionID uuid.UUID) []jobs.SessionJob {
	var sessionJobs []jobs.SessionJob
	dbCriterias := *sss.criteriasRepository.GetAll()
	for _, criteria := range dbCriterias {
		criteriaMarketsJobs := sss.newJobsFromCriteriaMarkets(criteria, sessionID)
		sessionJobs = append(sessionJobs, criteriaMarketsJobs...)
	}
	return sessionJobs
}

// newJobsFromCriteriaMarkets Creates jobs for all markets for a criteria
func (sss SessionStarterService) newJobsFromCriteriaMarkets(criteria adsdb.Criteria, sessionID uuid.UUID) []jobs.SessionJob {
	var rsjs []jobs.SessionJob
	for _, market := range *criteria.Markets {
		//url := sss.urlComposerImplementations.GetComposerImplementation(market.Name).Create(criteria)
		rsj := jobs.SessionJob{
			SessionID:  sessionID,
			JobID:      uuid.New(),
			CriteriaID: criteria.ID,
			MarketID:   market.ID,
			Criteria: jobs.Criteria{
				Brand:    criteria.Brand,
				CarModel: criteria.CarModel,
				YearFrom: criteria.YearFrom,
				YearTo:   criteria.YearTo,
				Fuel:     criteria.Fuel,
				KmFrom:   criteria.KmFrom,
				KmTo:     criteria.KmTo,
			},
			Market: jobs.Market{
				Name: market.Name,
				URL:  market.URL,
			},
			PageNumber: 1,
		}
		rsjs = append(rsjs, rsj)
	}
	return rsjs
}

func (sss SessionStarterService) startSession(session jobs.Session) {
	sss.pushSessionJobs(session.Jobs)
}

func (sss SessionStarterService) pushSessionJobs(jobs []jobs.SessionJob) {
	for _, job := range jobs {
		jobBytes, err := json.Marshal(&job)
		if err != nil {
			panic(err)
		}
		sss.messagesQueue.PutMessage(sss.criteriasTopicName, jobBytes)

	}
}
