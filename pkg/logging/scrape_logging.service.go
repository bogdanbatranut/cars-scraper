package logging

import (
	"carscraper/pkg/adsdb"
	"carscraper/pkg/amconfig"
	"carscraper/pkg/jobs"
	"fmt"
	"log"

	"github.com/google/uuid"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type ScrapeLoggingService struct {
	repo *LogsRepository
}

type LogsRepository struct {
	db *gorm.DB
}

func NewLogsRepository(cfg amconfig.IConfig) *LogsRepository {
	databaseName := cfg.GetString(amconfig.AppDBLogsName)
	databaseHost := cfg.GetString(amconfig.AppDBHost)
	dbUser := cfg.GetString(amconfig.AppDBUser)
	dbPass := cfg.GetString(amconfig.AppDBPass)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbUser, dbPass, databaseHost, databaseName)
	log.Println("LOGS DATABASE ", dsn)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	return &LogsRepository{
		db: db,
	}
}

func (sls ScrapeLoggingService) CreateSession(sessionID uuid.UUID) (*adsdb.SessionLog, error) {
	s := adsdb.SessionLog{
		SessionID: sessionID,
	}
	tx := sls.repo.db.Create(&s)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return &s, nil
}

func (sls ScrapeLoggingService) CreateCriteriaLog(sessionLog adsdb.SessionLog, job jobs.SessionJob) (*adsdb.CriteriaLog, error) {
	c := adsdb.CriteriaLog{
		SessionLogID: sessionLog.ID,
		SessionID:    job.SessionID,
		CriteriaID:   job.CriteriaID,
		MarketID:     job.MarketID,
		Brand:        job.Criteria.Brand,
		CarModel:     job.Criteria.CarModel,
		MarketName:   job.Market.Name,
		NumberOfAds:  0,
		Success:      false,
		Finished:     false,
	}

	tx := sls.repo.db.Create(&c)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return &c, nil
}

func (sls ScrapeLoggingService) CreatePageLog(criteriaLog *adsdb.CriteriaLog, job jobs.SessionJob, visitURL string, pageNumber int) (*adsdb.PageLog, error) {
	p := adsdb.PageLog{
		SessionLogID:  criteriaLog.SessionLogID,
		SessionID:     job.SessionID,
		SessionLog:    adsdb.SessionLog{},
		JobID:         job.JobID,
		CriteriaLogID: criteriaLog.ID,
		CriteriaLog:   adsdb.CriteriaLog{},
		VisitURL:      visitURL,
		Brand:         criteriaLog.Brand,
		CarModel:      criteriaLog.CarModel,
		MarketName:    criteriaLog.MarketName,
		MarketID:      job.MarketID,
		NumberOfAds:   0,
		PageNumber:    pageNumber,
		IsLastPage:    false,
		Error:         "",
		Scraped:       false,
		Consumed:      false,
	}
	tx := sls.repo.db.Create(&p)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return &p, nil
}

func (sls ScrapeLoggingService) PageLogSetVisitURL(pageLog *adsdb.PageLog, visitURL string) error {
	pageLog.VisitURL = visitURL
	tx := sls.repo.db.Updates(&pageLog)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

func (sls ScrapeLoggingService) PageLogSetPageScraped(pageLog *adsdb.PageLog, numberOFResults int, isLastPage bool) error {
	pageLog.Scraped = true
	pageLog.IsLastPage = isLastPage
	pageLog.NumberOfAds = numberOFResults
	tx := sls.repo.db.Updates(pageLog)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

func (sls ScrapeLoggingService) PageLogSetPageNumber(pageLog *adsdb.PageLog, pageNumber int) error {
	pageLog.PageNumber = pageNumber
	tx := sls.repo.db.Save(pageLog)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

func (sls ScrapeLoggingService) PageLogSetConsumed(pageLog *adsdb.PageLog) error {
	var existingPageLog adsdb.PageLog
	sls.repo.db.Find(&existingPageLog, pageLog.ID)
	existingPageLog.Consumed = true
	tx := sls.repo.db.Updates(&existingPageLog)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

func (sls ScrapeLoggingService) PageLogSetNumberOfAds(pageLog *adsdb.PageLog, nOfAds int) {
	pageLog.NumberOfAds = nOfAds
	sls.repo.db.Save(pageLog)
}

func (sls ScrapeLoggingService) PageLogSetError(pageLog *adsdb.PageLog, errStr string) {
	pageLog.Error = errStr
	sls.repo.db.Save(&pageLog)
}

func (sls ScrapeLoggingService) GetPageLog(sessionID uuid.UUID, jobID uuid.UUID, criteriaLogID uint, marketID uint, pageNumber int) *adsdb.PageLog {
	var cl adsdb.PageLog

	sls.repo.db.Where(&adsdb.PageLog{
		SessionID:     sessionID,
		JobID:         jobID,
		CriteriaLogID: criteriaLogID,
		MarketID:      marketID,
		PageNumber:    pageNumber,
	}).First(&cl)

	return &cl
}

func NewScrapeLoggingService(cfg amconfig.IConfig) *ScrapeLoggingService {
	return &ScrapeLoggingService{
		repo: NewLogsRepository(cfg),
	}
}

func (sls ScrapeLoggingService) GetDB() *gorm.DB {
	return sls.repo.db
}

func (sls ScrapeLoggingService) AddSession(job jobs.Session) error {
	log := adsdb.SessionLog{
		SessionID: job.SessionID,
	}
	tx := sls.repo.db.Create(&log)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

func (sls ScrapeLoggingService) CriteriaLogSetAsFinished(log adsdb.CriteriaLog) {
	log.Finished = true
	var existingLog adsdb.CriteriaLog
	sls.repo.db.First(&existingLog, log.ID)
	existingLog.Finished = true
	sls.repo.db.Updates(&existingLog)
}

func (sls ScrapeLoggingService) CriteriaLogSetSuccessful(log adsdb.CriteriaLog) {
	log.Success = true
	var existingLog adsdb.CriteriaLog
	sls.repo.db.First(&existingLog, log.ID)
	existingLog.Success = true
	sls.repo.db.Updates(&existingLog)
}

func (sls ScrapeLoggingService) CriteriaLogAddNumberOfAds(criteriaLog adsdb.CriteriaLog, nOfAds int) {
	var existingLog adsdb.CriteriaLog
	sls.repo.db.First(&existingLog, criteriaLog.ID)
	existingLog.NumberOfAds += nOfAds
	sls.repo.db.Updates(&existingLog)
}

func (sls ScrapeLoggingService) GetCriteriaLog(sessionID uuid.UUID, criteriaID uint, marketID uint) (*adsdb.CriteriaLog, error) {
	search := adsdb.CriteriaLog{
		SessionID:  sessionID,
		CriteriaID: criteriaID,
		MarketID:   marketID,
	}

	tx := sls.repo.db.Where(&search).First(&search)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return &search, nil
}
