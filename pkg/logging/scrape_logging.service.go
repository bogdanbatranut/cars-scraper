package logging

import (
	"carscraper/pkg/adsdb"
	"carscraper/pkg/amconfig"
	"carscraper/pkg/jobs"
	"fmt"

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
	databaseName := cfg.GetString(amconfig.AppDBName)
	databaseHost := cfg.GetString(amconfig.AppDBHost)
	dbUser := cfg.GetString(amconfig.AppDBUser)
	dbPass := cfg.GetString(amconfig.AppDBPass)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbUser, dbPass, databaseHost, databaseName)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	return &LogsRepository{
		db: db,
	}
}

func (r LogsRepository) AddEntry(entry adsdb.ScrapeLog) error {
	tx := r.db.Create(&entry)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

func NewScrapeLoggingService(cfg amconfig.IConfig) ScrapeLoggingService {
	return ScrapeLoggingService{
		repo: NewLogsRepository(cfg),
	}
}

func (sls ScrapeLoggingService) GetDB() *gorm.DB {
	return sls.repo.db
}
func (sls ScrapeLoggingService) AddCriteriaEntry(job jobs.SessionJob, numberOfAds int, error string, success bool) error {
	log := adsdb.CriteriaLog{
		SessionID:   job.SessionID,
		Brand:       job.Criteria.Brand,
		CarModel:    job.Criteria.CarModel,
		MarketName:  job.Market.Name,
		NumberOfAds: numberOfAds,
		Error:       error,
		Success:     success,
	}
	tx := sls.repo.db.Create(&log)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

func (sls ScrapeLoggingService) AddPageScrapeEntry(job jobs.SessionJob, numberOfAds int, pageNumber int, isLastPage bool, visitURL string, err error) error {
	errStr := ""
	if err != nil {
		errStr = err.Error()
	}

	logEntry := adsdb.ScrapeLog{
		SessionID:   job.SessionID,
		JobID:       job.JobID,
		VisitURL:    visitURL,
		Brand:       job.Criteria.Brand,
		CarModel:    job.Criteria.CarModel,
		MarketName:  job.Market.Name,
		NumberOfAds: numberOfAds,
		PageNumber:  pageNumber,
		IsLastPage:  isLastPage,
		Error:       errStr,
	}

	err_ := sls.repo.AddEntry(logEntry)
	if err_ != nil {
		return err_
	}
	return nil
}
