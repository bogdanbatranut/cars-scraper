package repos

import (
	"carscraper/pkg/adsdb"
	"carscraper/pkg/amconfig"
	"carscraper/pkg/events"
	"fmt"
	"log"
	"math"
	"strings"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Pagination struct {
	Limit      int   `json:"limit,omitempty;query:limit"`
	Page       int   `json:"page,omitempty;query:page"`
	TotalRows  int64 `json:"total_rows"`
	TotalPages int   `json:"total_pages"`
	//Rows       interface{} `json:"rows"`
}

type AdsRepository struct {
	db             *gorm.DB
	eventsListener *events.EventsListener
}

func NewAdsDB(cfg amconfig.IConfig) *gorm.DB {
	databaseName := cfg.GetString(amconfig.AppDBName)
	databaseHost := cfg.GetString(amconfig.AppDBHost)
	dbUser := cfg.GetString(amconfig.AppDBUser)
	dbPass := cfg.GetString(amconfig.AppDBPass)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbUser, dbPass, databaseHost, databaseName)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	return db
}

func NewAdsRepository(cfg amconfig.IConfig, eventsListener *events.EventsListener) *AdsRepository {
	databaseName := cfg.GetString(amconfig.AppDBName)
	databaseHost := cfg.GetString(amconfig.AppDBHost)
	dbUser := cfg.GetString(amconfig.AppDBUser)
	dbPass := cfg.GetString(amconfig.AppDBPass)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbUser, dbPass, databaseHost, databaseName)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	return &AdsRepository{
		db:             db,
		eventsListener: eventsListener,
	}
}

func (r AdsRepository) GetDB() *gorm.DB {
	return r.db
}

func (r AdsRepository) GetSellerAds(dealerID uint) *[]adsdb.Ad {
	var ads []adsdb.Ad
	tx := r.db.Unscoped().Preload("Seller").Preload("Prices").Where("seller_id = ?", dealerID).Find(&ads)
	tx.Find(&ads)
	return &ads
}

func (r AdsRepository) GetAll() (*[]adsdb.Ad, error) {
	var ads []adsdb.Ad
	tx := r.db.Find(&ads)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return &ads, nil
}

func (r AdsRepository) Upsert(ads []adsdb.Ad) (*[]uint, error) {
	adsIds := []uint{}
	for _, foundAd := range ads {
		foundAdPrice := foundAd.Prices[0].Price
		foundAdKm := foundAd.Km

		var dbAd adsdb.Ad
		tx := r.db.Raw("SELECT * FROM `ads` WHERE  `ads`.`market_uuid` = ? LIMIT 1 ", foundAd.MarketUUID).Scan(&dbAd)
		if tx.Error != nil {
			log.Println(tx.Error)
			return nil, tx.Error
		}

		if dbAd.Model == nil {
			// we have a new ad
			tx = r.db.Create(&foundAd)
			if tx.Error != nil {
				log.Println(tx.Error)
				return nil, tx.Error
			}
			dbAd = foundAd

		} else {
			// we have the ad in the db
			// if the ad is inactive we must activate
			if dbAd.Model.DeletedAt.Valid {
				tx = r.db.Unscoped().Model(&adsdb.Ad{}).Where("id", dbAd.ID).Update("deleted_at", nil)
				if tx.Error != nil {
					log.Println(tx.Error)
					return nil, tx.Error
				}
			}
			r.db.First(&dbAd, dbAd.ID)

			if foundAd.Title != nil {
				r.db.Model(&dbAd).Update("title", *foundAd.Title)
			}
			// set new values if they exist
			if foundAd.Km != 0 && foundAd.Km != dbAd.Km {
				r.db.Model(&dbAd).Update("km", foundAdKm)
			}
			if foundAd.Title != nil {
				r.db.Model(&dbAd).Update("title", *foundAd.Title)
			}

			if foundAd.Ad_url != "" && foundAd.Ad_url != dbAd.Ad_url {
				adURL := foundAd.Ad_url
				if foundAd.MarketID == 11 {
					if !strings.Contains(foundAd.Ad_url, "www.mobile.de") {
						adURL = fmt.Sprintf("https://www.mobile.de%s", adURL)
					}
				}
				if foundAd.MarketID == 12 {
					if !strings.Contains(foundAd.Ad_url, "www.autoscout24.ro") {
						adURL = fmt.Sprintf("https://www.autoscout24.ro%s", adURL)
					}
				}
				r.db.Model(&dbAd).Update("ad_url", adURL)

			}
			if foundAd.CurrentPrice != nil && foundAd.CurrentPrice != dbAd.CurrentPrice {
				r.db.Model(&dbAd).Update("current_price", *foundAd.CurrentPrice)
			}
			if foundAd.Thumbnail != nil {
				if dbAd.Thumbnail != nil {
					if *foundAd.Thumbnail != *dbAd.Thumbnail {
						r.db.Model(&dbAd).Update("thumbnail", *foundAd.Thumbnail)
					}
				} else {
					r.db.Model(&dbAd).Update("thumbnail", *foundAd.Thumbnail)
				}
			}

		}

		// get last price id db
		var lastExistingPrice adsdb.Price
		tx = r.db.Last(&lastExistingPrice, adsdb.Price{AdID: dbAd.ID})
		if tx.Error != nil {
			log.Println(tx.Error)
			return nil, tx.Error
		}

		if lastExistingPrice.Price != foundAdPrice {
			// insert new price
			tx = r.db.Create(&adsdb.Price{Price: foundAdPrice, AdID: dbAd.ID, MarketID: dbAd.MarketID})
			if tx.Error != nil {
				log.Println(tx.Error)
				return nil, tx.Error
			}
			var prices []adsdb.Price
			tx = r.db.Model(&adsdb.Price{}).Where("ad_id = ?", dbAd.ID).Find(&prices)
			if tx.Error != nil {
				// no need to return ... just will not fire the event
				log.Println(tx.Error)
				log.Println("notification will no be fired")
			}
			dbAd.Prices = prices
			if dbAd.Followed {
				r.eventsListener.Fire(events.UpdatePriceEvent{Ad: dbAd})
			}
		}

		adsIds = append(adsIds, dbAd.ID)
	}
	return &adsIds, nil
}

func (r AdsRepository) DeleteAd(adID uint) {
	var ad adsdb.Ad
	tx := r.db.Find(&ad, adID)
	if tx.Error != nil {
		log.Println(tx.Error)
		return
	}
	if ad.Followed {
		r.eventsListener.Fire(events.DeleteEvent{Ad: ad})
	}
	r.db.Model(adsdb.Ad{}).Delete(&ad, adID)
}

func (r AdsRepository) DeletePrice(priceID uint) {
	model := gorm.Model{
		ID: priceID,
	}
	price := adsdb.Price{
		Model: &model,
	}
	r.db.Model(adsdb.Price{}).Unscoped().Delete(&price)
}

func (r AdsRepository) GetAllAdsIDs(marketID uint, criteriaID uint) *[]uint {
	var ads []adsdb.Ad
	var adsIDs []uint
	r.db.Model(adsdb.Ad{}).Select("id").Where("market_id = ? AND criteria_id = ?", marketID, criteriaID).Find(&ads)
	for _, ad := range ads {
		adsIDs = append(adsIDs, ad.ID)
	}
	return &adsIDs
}

func (r AdsRepository) GetAdsForCriteria(criteriaID uint, markets []string, minKm *int, maxKm *int, minPrice *int, maxPrice *int, years *[]string) *[]adsdb.Ad {
	var ads []adsdb.Ad
	//r.db.Preload("Prices").Preload("Market").Where("criteria_id = ?", criteriaID).Where("market_id", markets).Where("current_price <= ?", maxPrice).Where("current_price >= ? ", minPrice).Find(&ads)
	tx := r.db.Debug().Preload("Prices").Preload("Market").Preload("Seller").Where("criteria_id = ?", criteriaID).Where("market_id", markets).Where("current_price <= ?", maxPrice).Where("current_price >= ? ", minPrice)
	if years != nil {
		//tx = tx.Where("Year IN (?)", years)
		tx = tx.Where("year", *years)
	}
	tx.Find(&ads)
	return &ads
}

func withLimitAndOffset(pagination *Pagination, db *gorm.DB) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Offset(pagination.GetOffset()).Limit(pagination.GetLimit())
	}
}

func WithFilterConditions(db *gorm.DB, criteriaID uint, markets []string, minKm *int, maxKm *int, minPrice *int, maxPrice *int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Preload("Prices").Preload("Market").Where("criteria_id = ?", criteriaID).Where("market_id", markets).Where("current_price <= ?", maxPrice).Where("current_price >= ? ", minPrice)
	}
}

//func WithFilterConditions(db *gorm.DB, criteriaID uint, markets []string, minKm *int, maxKm *int, minPrice *int, maxPrice *int) *gorm.DB {
//	return db.Preload("Prices").Preload("Market").Where("criteria_id = ?", criteriaID).Where("market_id", markets).Where("current_price <= ?", maxPrice).Where("current_price >= ? ", minPrice)
//}

func (r AdsRepository) GetAdsForCriteriaPaginated(pagination *Pagination, criteriaID uint, markets []string, minKm *int, maxKm *int, minPrice *int, maxPrice *int) (*[]adsdb.Ad, *Pagination) {
	var ads []adsdb.Ad
	//tx := r.db.Scopes(paginate(ads, pagination, r.db)).Preload("Prices").Preload("Market").Where("criteria_id = ?", criteriaID).Where("market_id", markets).Where("current_price <= ?", maxPrice).Where("current_price >= ? ", minPrice).Find(&ads)
	//tx := r.db.Scopes(withLimitAndOffset(pagination, r.db)).Preload("Prices").Preload("Market").Where("criteria_id = ?", criteriaID).Where("market_id", markets).Where("current_price <= ?", maxPrice).Where("current_price >= ? ", minPrice).Find(&ads)
	tx := r.db.Debug().Scopes(WithFilterConditions(r.db, criteriaID, markets, minKm, maxKm, minPrice, maxPrice)).Find(&ads)
	if tx.Error != nil {
		panic(tx.Error)
	}
	var totalRows int64
	tx = r.db.Debug().Model(&adsdb.Ad{}).Scopes(WithFilterConditions(r.db, criteriaID, markets, minKm, maxKm, minPrice, maxPrice)).Count(&totalRows)
	if tx.Error != nil {
		panic(tx.Error)
	}
	pagination.TotalRows = totalRows
	totalPages := int(math.Ceil(float64(totalRows) / float64(pagination.Limit)))
	pagination.TotalPages = totalPages

	return &ads, pagination
}

func (r AdsRepository) GetAdPrices(adID uint) []adsdb.Price {
	var prices []adsdb.Price
	price := adsdb.Price{AdID: adID}
	tx := r.db.Find(&prices, price)
	if tx.Error != nil {
		panic(tx.Error)
	}
	return prices
}

func (r AdsRepository) UpdateCurrentPrice(adID uint) {
	var ad adsdb.Ad
	tx := r.db.Preload("Prices").Find(&ad, adID)
	if tx.Error != nil {
		panic(tx.Error)
	}
	lastPrice := ad.Prices[len(ad.Prices)-1].Price
	ad.CurrentPrice = &lastPrice
	tx = r.db.Save(&ad)
	if tx.Error != nil {
		panic(tx.Error)
	}
}

func (p *Pagination) GetOffset() int {
	return (p.GetPage() - 1) * p.GetLimit()
}

func (p *Pagination) GetLimit() int {
	if p.Limit == 0 {
		p.Limit = 10
	}
	return p.Limit
}
func (p *Pagination) GetPage() int {
	if p.Page == 0 {
		p.Page = 1
	}
	return p.Page
}

func AmountGreaterThan1000(db *gorm.DB, maxPrice int) *gorm.DB {
	return db.Where("current_price < ?", maxPrice)
}

func paginate(value interface{}, pagination *Pagination, db *gorm.DB) func(db *gorm.DB) *gorm.DB {
	var totalRows int64
	db.Debug().Model(value).Count(&totalRows)

	pagination.TotalRows = totalRows
	totalPages := int(math.Ceil(float64(totalRows) / float64(pagination.Limit)))
	pagination.TotalPages = totalPages

	return func(db *gorm.DB) *gorm.DB {
		return db.Offset(pagination.GetOffset()).Limit(pagination.GetLimit())
	}
}

func (r AdsRepository) GetTodaysLowestPriceInCriteria(criteriaID uint) (*adsdb.Ad, error) {
	var ad adsdb.Ad

	// Query to find the ad with the lowest price updated today for the given criteriaID
	tx := r.db.Model(&adsdb.Ad{}).
		Joins("JOIN prices ON prices.ad_id = ads.id").
		Where("ads.criteria_id = ?", criteriaID).
		Where("DATE(prices.created_at) = CURDATE()").
		Order("prices.price ASC").
		Preload("Prices"). // Preload prices for the ad
		First(&ad)

	if tx.Error != nil {
		if tx.Error == gorm.ErrRecordNotFound {
			return nil, nil // No ad found with a price updated today
		}
		return nil, tx.Error // Return other errors
	}

	return &ad, nil
}

func (r AdsRepository) GetNewEntryLowestPrice(criteriaID uint) (*adsdb.Ad, error) {
	var ad adsdb.Ad

	// Query to find the ad with only one price and the price created today
	tx := r.db.Model(&adsdb.Ad{}).
		Joins("JOIN prices ON prices.ad_id = ads.id").
		Where("ads.criteria_id = ?", criteriaID).
		Where("DATE(prices.created_at) = CURDATE()").
		Group("ads.id").
		Having("COUNT(prices.id) = 1").
		Order("MIN(prices.price) ASC").
		Preload("Prices"). // Preload prices for the ad
		First(&ad)

	if tx.Error != nil {
		if tx.Error == gorm.ErrRecordNotFound {
			return nil, nil // No ad found matching the criteria
		}
		return nil, tx.Error // Return other errors
	}

	return &ad, nil
}

func (r AdsRepository) GetMinPriceForCriteria(criteriaID uint) ([]adsdb.Ad, error) {
	//SELECT
	//a.criteria_id,
	//	a.id AS ad_id,
	//	a.title,
	//	a.brand,
	//	a.car_model,
	//	a.current_price AS lowest_price
	//FROM
	//automall.ads a
	//JOIN (
	//	SELECT
	//criteria_id,
	//	MIN(current_price) AS lowest_price
	//FROM
	//automall.ads
	//WHERE
	//deleted_at IS NULL
	//GROUP BY
	//criteria_id
	//) subquery ON a.criteria_id = subquery.criteria_id AND a.current_price = subquery.lowest_price
	//WHERE
	//a.deleted_at IS NULL
	//ORDER BY
	//a.criteria_id ASC;

	var ads []adsdb.Ad

	// Subquery to find the minimum current_price for each criteria_id
	subquery := r.db.Model(&adsdb.Ad{}).
		Select("criteria_id, MIN(current_price) AS lowest_price").
		Where("deleted_at IS NULL").
		Group("criteria_id")

	// Main query to join the subquery and fetch the required fields
	tx := r.db.Model(&adsdb.Ad{}).
		Select("ads.criteria_id, ads.id AS id, ads.title, ads.brand, ads.car_model, ads.current_price AS lowest_price").
		Joins("JOIN (?) AS subquery ON ads.criteria_id = subquery.criteria_id AND ads.current_price = subquery.lowest_price", subquery).
		Where("ads.deleted_at IS NULL").
		Where("ads.criteria_id = ?", criteriaID). // Pass criteria_id as a parameter
		Order("ads.criteria_id ASC").
		Find(&ads)

	if tx.Error != nil {
		panic(tx.Error)
	}
	return ads, nil
}
