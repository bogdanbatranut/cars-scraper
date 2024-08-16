package sessionstarter

import (
	"carscraper/pkg/adsdb"
	"carscraper/pkg/amconfig"
	"carscraper/pkg/jobs"
	"carscraper/pkg/logging"
	"carscraper/pkg/repos"
	"carscraper/pkg/url"
	"encoding/json"
	"fmt"
	"log"

	"github.com/google/uuid"
)

type SessionStarterService struct {
	messagesQueue              repos.IMessageQueue
	criteriasRepository        repos.ICriteriaRepository
	marketsRepository          repos.IMarketsRepository
	logger                     *logging.ScrapeLoggingService
	requesteScrapingJobs       []jobs.Session
	urlComposerImplementations *url.URLComposerImplementations
	jobsTopicName              string
}

type CrawlinglInitiatorServiceConfiguration func(cisc *SessionStarterService)

func NewSessionStarterService(cfgs ...CrawlinglInitiatorServiceConfiguration) *SessionStarterService {

	mqs := &SessionStarterService{
		urlComposerImplementations: url.NewURLComposerImplementations(),
		jobsTopicName:              "jobs",
	}
	for _, cfg := range cfgs {
		cfg(mqs)
	}
	return mqs
}

func WithCriteriaSQLRepository(cfg amconfig.IConfig) CrawlinglInitiatorServiceConfiguration {
	return func(mqs *SessionStarterService) {
		mqs.criteriasRepository = repos.NewSQLCriteriaRepository(cfg)
	}
}

func WithMarketsSQLRepository(cfg amconfig.IConfig) CrawlinglInitiatorServiceConfiguration {
	return func(mqs *SessionStarterService) {
		mqs.marketsRepository = repos.NewSQLMarketsRepository(cfg)
	}
}

func WithLogging(cfg amconfig.IConfig) CrawlinglInitiatorServiceConfiguration {
	return func(mqs *SessionStarterService) {
		mqs.logger = logging.NewScrapeLoggingService(cfg)
	}
}

func WithSimpleMessageQueueRepository(cfg amconfig.IConfig) CrawlinglInitiatorServiceConfiguration {
	smqBaseURL := cfg.GetString(amconfig.SMQURL)
	smqPort := cfg.GetString(amconfig.SMQHTTPPort)
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
	log.Printf("Pushed session %+v", session.SessionID)

}

func (sss SessionStarterService) newSession() jobs.Session {
	log.Println("creating session")
	sessionID := uuid.New()
	sessionJobs := sss.createSessionJobs(sessionID)
	session := jobs.Session{
		SessionID: sessionID,
		Jobs:      sessionJobs,
	}
	return session

}

func inArrayUINT(str uint, list []uint) bool {
	for _, v := range list {
		if v == str {
			return true
		}
	}
	return false
}

func (sss SessionStarterService) createSessionJobs(sessionID uuid.UUID) []jobs.SessionJob {
	sessionJobs := []jobs.SessionJob{}

	markets := []uint{9, 11, 12, 13, 14, 15, 16, 17, 18, 19}
	//criterias := []uint{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29}
	//markets := []uint{9}
	criterias := []uint{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31}

	allowedMarketAutoklassCriterias := []uint{8, 9, 24, 6, 13, 4, 1, 5, 27, 25, 28, 3, 10, 11, 19, 14}
	allowedMercedesBenzCriterias := []uint{3, 10, 11, 14, 19}

	createSession, err := sss.logger.CreateSession(sessionID)
	if err != nil {
		panic(err)
	}

	for _, marketID := range markets {
		if marketID == 10 {
			marketID++
			continue
		}
		for _, criteriaID := range criterias {
			// criteria 7 volvo s90
			// criteria 6 bmw 7 series
			job := sss.createJob(sessionID, criteriaID, marketID)
			// do not scrape other brands for ofertebmw
			if job.Criteria.Brand != "bmw" && marketID == 15 {
				continue
			}

			if marketID == 18 && !inArrayUINT(criteriaID, allowedMarketAutoklassCriterias) {
				continue
			}

			if marketID == 17 || marketID == 19 {
				if !inArrayUINT(criteriaID, allowedMercedesBenzCriterias) {
					continue
				}
			}

			if job.Criteria.Brand != "mercedes-benz" && marketID == 17 {
				continue
			}
			log.Println(" APPENDING JOB : ", job.ToString())
			sessionJobs = append(sessionJobs, job)
			_, err := sss.logger.CreateCriteriaLog(*createSession, job)
			if err != nil {
				panic(err)
			}
		}
	}
	return sessionJobs
}

// this should be the permanenent but while in Vacation in Greece I will use the func above
//func (sss SessionStarterService) createSessionJobs(sessionID uuid.UUID) []jobs.SessionJob {
//	sessionJobs := []jobs.SessionJob{}
//	dbCriterias := *sss.criteriasRepository.GetAll()
//	for _, criteria := range dbCriterias {
//
//		if !criteria.AllowProcess {
//			continue
//		}
//		log.Println("Criteria : ", criteria.CarModel)
//		criteriaMarketsJobs := sss.newJobsFromCriteriaMarkets(criteria, sessionID)
//		if criteriaMarketsJobs != nil {
//			sessionJobs = append(sessionJobs, criteriaMarketsJobs...)
//		}
//	}
//	return sessionJobs
//}

// newJobsFromCriteriaMarkets Creates jobs for all markets for a criteria
func (sss SessionStarterService) newJobsFromCriteriaMarkets(criteria adsdb.Criteria, sessionID uuid.UUID) []jobs.SessionJob {
	var rsjs []jobs.SessionJob
	for _, market := range *criteria.Markets {
		if !market.AllowProcess {
			continue
		}
		log.Println("Market :", market.Name)
		//url := sss.urlComposerImplementations.GetComposerImplementation(market.Name).Create(criteria)
		rsj := jobs.SessionJob{
			AllowIncrementPage: true,
			SessionID:          sessionID,
			JobID:              uuid.New(),
			CriteriaID:         criteria.ID,
			MarketID:           market.ID,
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
				Name:       market.Name,
				PageNumber: 1,
			},
		}
		rsjs = append(rsjs, rsj)
	}
	return rsjs
}

func (sss SessionStarterService) createJob(sessionID uuid.UUID, criteriaID uint, marketID uint) jobs.SessionJob {

	criteria := sss.criteriasRepository.GetCriteriaByID(criteriaID)
	market := sss.marketsRepository.GetMarketByID(marketID)
	rsj := jobs.SessionJob{
		AllowIncrementPage: true,
		SessionID:          sessionID,
		JobID:              uuid.New(),
		CriteriaID:         criteriaID,
		MarketID:           marketID,
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
			Name:       market.Name,
			PageNumber: 1,
		},
	}
	return rsj
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
		log.Printf("SESSION STARTER - Pushing job %s", job.ToString())
		sss.messagesQueue.PutMessage(sss.jobsTopicName, jobBytes)
	}
}
