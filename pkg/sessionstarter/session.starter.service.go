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
	"strconv"

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
	log.Println("SMQ URL : ", smqURL)
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
}

func (sss SessionStarterService) ScrapeMarket(marketID string) {
	session := sss.newSessionForMarket(marketID)
	sss.pushSessionJobs(session.Jobs)
}

func (sss SessionStarterService) ScrapeMarketCriteria(marketID string, criteriaID string) {
	session := sss.newSessionForMarketCriteria(marketID, criteriaID)
	sss.pushSessionJobs(session.Jobs)
}

func (sss SessionStarterService) ScrapeMarketsCriterias(markets []uint, criterias []uint) {
	sessionID := uuid.New()
	session := sss.newSessionForMarketsCriterias(markets, criterias, sessionID)
	sss.pushSessionJobs(session.Jobs)
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

func (sss SessionStarterService) newSessionForMarket(marketID string) jobs.Session {
	log.Println("creating session for market ", marketID)
	sessionID := uuid.New()
	sessionJobs, err := sss.createSessionJobsForMarket(sessionID, marketID)
	if err != nil {
		panic(err)
	}
	session := jobs.Session{
		SessionID: sessionID,
		Jobs:      sessionJobs,
	}
	return session
}

func (sss SessionStarterService) newSessionForMarketCriteria(marketID string, criteriaID string) jobs.Session {
	log.Println(fmt.Sprintf("creating session for market %s and criteria %s ", marketID, criteriaID))
	sessionID := uuid.New()
	sessionJobs, err := sss.createSessionJobsForMarketCriteria(sessionID, marketID, criteriaID)
	if err != nil {
		panic(err)
	}
	session := jobs.Session{
		SessionID: sessionID,
		Jobs:      sessionJobs,
	}
	return session
}

func (sss SessionStarterService) newSessionForMarketsCriterias(markets []uint, criterias []uint, sessionID uuid.UUID) jobs.Session {
	log.Println(fmt.Sprintf("creating session for market %v and criteria %v ", markets, criterias))
	sessionJobs, err := sss.createSessionJobsForMarketsCriterias(sessionID, markets, criterias)
	if err != nil {
		panic(err)
	}
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

func (sss SessionStarterService) createSessionJobsForMarketsCriterias(sessionID uuid.UUID, markets []uint, criterias []uint) ([]jobs.SessionJob, error) {
	allowedMarketAutoklassCriterias := []uint{8, 9, 24, 6, 13, 4, 1, 5, 27, 25, 28, 3, 10, 11, 19, 14}
	allowedMercedesBenzCriterias := []uint{3, 10, 11, 14, 19, 34, 33, 40, 42}
	allowedBMWDECriterias := []uint{1, 2, 5, 6, 13, 35, 36, 39, 41, 43, 44, 46}

	sessionJobs := []jobs.SessionJob{}

	createSession, err := sss.logger.CreateSession(sessionID)
	if err != nil {
		panic(err)
	}

	log.Println(" Created Session !!! ", createSession.ID)

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

			if marketID == 20 && !inArrayUINT(criteriaID, allowedBMWDECriterias) {
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
			clog, err := sss.logger.CreateCriteriaLog(*createSession, job)
			log.Println("CREATED Criteria log ", clog.ID)
			if err != nil {
				return nil, err
			}
		}
	}
	return sessionJobs, nil

}

func (sss SessionStarterService) createSessionJobsForMarketCriteria(sessionID uuid.UUID, marketIDStr string, criteriaIDStr string) ([]jobs.SessionJob, error) {
	sessionJobs := []jobs.SessionJob{}
	markID, err := strconv.Atoi(marketIDStr)
	if err != nil {
		return nil, err
	}
	marketID := uint(markID)

	critID, err := strconv.Atoi(criteriaIDStr)
	if err != nil {
		return nil, err
	}

	criteriaID := uint(critID)

	markets := []uint{marketID}
	criterias := []uint{criteriaID}

	allowedMarketAutoklassCriterias := []uint{8, 9, 24, 6, 13, 4, 1, 5, 27, 25, 28, 3, 10, 11, 19, 14}
	allowedMercedesBenzCriterias := []uint{3, 10, 11, 14, 19, 34, 33, 40, 42}
	allowedBMWDECriterias := []uint{1, 2, 5, 6, 13, 35, 36, 39, 41, 43, 44, 46}

	createSession, err := sss.logger.CreateSession(sessionID)
	if err != nil {
		panic(err)
	}

	log.Println(" Created Session !!! ", createSession.ID)

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

			if marketID == 20 && !inArrayUINT(criteriaID, allowedBMWDECriterias) {
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
			clog, err := sss.logger.CreateCriteriaLog(*createSession, job)
			log.Println("CREATED Criteria log ", clog.ID)
			if err != nil {
				return nil, err
			}
		}
	}
	return sessionJobs, nil
}

func (sss SessionStarterService) createSessionJobsForMarket(sessionID uuid.UUID, marketIDStr string) ([]jobs.SessionJob, error) {
	sessionJobs := []jobs.SessionJob{}
	markID, err := strconv.Atoi(marketIDStr)
	if err != nil {
		return nil, err
	}
	marketID := uint(markID)

	markets := []uint{marketID}
	criterias := []uint{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 46, 47}

	allowedMarketAutoklassCriterias := []uint{8, 9, 24, 6, 13, 4, 1, 5, 27, 25, 28, 3, 10, 11, 19, 14}
	allowedMercedesBenzCriterias := []uint{3, 10, 11, 14, 19, 34, 33, 40, 42}
	allowedBMWDECriterias := []uint{1, 2, 5, 6, 13, 35, 36, 39, 41, 43, 44, 46}

	createSession, err := sss.logger.CreateSession(sessionID)
	if err != nil {
		panic(err)
	}

	log.Println(" Created Session !!! ", createSession.ID)

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

			if marketID == 20 && !inArrayUINT(criteriaID, allowedBMWDECriterias) {
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
			clog, err := sss.logger.CreateCriteriaLog(*createSession, job)
			log.Println("CREATED Criteria log ", clog.ID)
			if err != nil {
				return nil, err
			}
		}
	}
	return sessionJobs, nil
}

func (sss SessionStarterService) createSessionJobs(sessionID uuid.UUID) []jobs.SessionJob {
	sessionJobs := []jobs.SessionJob{}

	markets := []uint{9, 11, 12, 14, 15, 16, 17, 18, 19, 20}
	//criterias := []uint{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29}
	//markets := []uint{9}
	criterias := []uint{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 46, 47, 48}

	allowedMarketAutoklassCriterias := []uint{8, 9, 24, 6, 13, 4, 1, 5, 27, 25, 28, 3, 10, 11, 19, 14}
	allowedMercedesBenzCriterias := []uint{3, 10, 11, 14, 19, 34, 33, 40, 42}
	allowedBMWDECriterias := []uint{1, 2, 5, 6, 13, 35, 36, 39, 41, 43, 44, 46}
	alowedSkodaKodiakinMarkets := []uint{9, 11, 12, 13}

	createSession, err := sss.logger.CreateSession(sessionID)
	if err != nil {
		panic(err)
	}

	log.Println(" Created Session !!! ", createSession.ID)

	for _, marketID := range markets {
		if marketID == 10 {
			marketID++
			continue
		}
		for _, criteriaID := range criterias {

			if criteriaID == 48 && !inArrayUINT(marketID, alowedSkodaKodiakinMarkets) {
				continue
			}

			// criteria 7 volvo s90
			// criteria 6 bmw 7 series
			job := sss.createJob(sessionID, criteriaID, marketID)
			// do not scrape other brands for ofertebmw
			if job.Criteria.Brand != "bmw" && marketID == 15 {
				continue
			}

			if marketID == 20 && !inArrayUINT(criteriaID, allowedBMWDECriterias) {
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
			clog, err := sss.logger.CreateCriteriaLog(*createSession, job)
			log.Println("CREATED Criteria log ", clog.ID)
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
