package main

import (
	"carscraper/pkg/adsdb"
	"carscraper/pkg/amconfig"
	"carscraper/pkg/errorshandler"
	"carscraper/pkg/events"
	"carscraper/pkg/notifications"
	"carscraper/pkg/repos"
	"carscraper/pkg/statistics/calculators/age"
	"carscraper/pkg/statistics/calculators/discount"
	"carscraper/pkg/valueobjects"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func main() {
	log.Println("starting BACKEND service...")

	cfg, err := amconfig.NewViperConfig()
	errorshandler.HandleErr(err)

	r := mux.NewRouter().StrictSlash(true)
	r.Use(enableCORS)

	notificationsService := notifications.NewNotificationsService(cfg)
	eventsListener := events.NewEventsListener(notificationsService)

	criteriaRepo := repos.NewSQLCriteriaRepository(cfg)
	marketsRepo := repos.NewSQLMarketsRepository(cfg)
	adsRepo := repos.NewAdsRepository(cfg, eventsListener)
	adsDB := repos.NewAdsDB(cfg)

	//chartsRepo := repos.NewChartsRepository(cfg)
	//chartsRepo.GetAdsPricesByStep(5000)

	//cleanupPrices(adsRepo)

	r.HandleFunc("/updatePrices", setCurrentPrice(adsRepo)).Methods("GET")

	r.HandleFunc("/markets", getMarkets(marketsRepo)).Methods("GET", "OPTIONS")
	r.HandleFunc("/criterias", getCriterias(criteriaRepo)).Methods("GET", "OPTIONS")
	r.HandleFunc("/adsforcriteria/{id}", getAdsForCriteria(adsRepo)).Methods("GET", "OPTIONS")
	r.HandleFunc("/adsforcriteriaPaginated/{id}", getAdsForCriteriaPaginated(adsRepo)).Methods("GET", "OPTIONS")

	r.HandleFunc("/follow", follow(adsDB)).Methods("POST")

	r.HandleFunc("/test", test()).Methods("POST")
	r.HandleFunc("/ad/{id}", getAd(adsDB)).Methods("GET", "OPTIONS")

	r.HandleFunc("/marketsAndCriterias", marketsAndCriterias(criteriaRepo)).Methods("POST")

	httpPort := cfg.GetString(amconfig.BackendServiceHTTPPort)
	log.Printf("HTTP listening on port %s\n", httpPort)
	err = http.ListenAndServe(fmt.Sprintf(":%s", httpPort), enableCORS(r))
	errorshandler.HandleErr(err)

}

func setCurrentPrice(repo *repos.AdsRepository) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var ads *[]adsdb.Ad
		ads, err := repo.GetAll()
		if err != nil {
			panic(err)
		}

		for _, ad := range *ads {
			repo.UpdateCurrentPrice(ad.ID)
		}

		w.Write([]byte("done !!!"))
	}
}

func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Set headers
		w.Header().Set("Access-Control-Allow-Headers:", "*")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "*")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		fmt.Println("ok")

		// Next
		next.ServeHTTP(w, r)
		return
	})
}

func follow(repo *gorm.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			w.WriteHeader(http.StatusOK)
			return
		}

		// Set CORS headers for the actual request
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json")

		type FollowRequest struct {
			AdID   uint `json:"adID"`
			Follow bool `json:"follow"`
		}

		var followRequest FollowRequest
		err := json.NewDecoder(r.Body).Decode(&followRequest)
		if err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		setFollow(followRequest.AdID, followRequest.Follow, repo)

		w.WriteHeader(http.StatusOK)
	}
}

func setFollow(adID uint, follow bool, db *gorm.DB) {
	var ad adsdb.Ad
	result := db.Model(&ad).Where("id = ?", adID).First(&ad)
	ad.Followed = follow
	db.Debug().Save(&ad)
	//ad.Update("follow", follow)
	if result.Error != nil {
		log.Println("error updating follow status:", result.Error)
	}
	log.Printf("Rows affected: %d", result.RowsAffected)
}

func test() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var bodyBytes []byte
		var err error
		w.Header().Add("Access-Control-Allow-Origin", "*")
		//w.Header().Add("Access-Control-Allow-Methods", "POST")
		//w.Header().Add("Content-Type", "application/json")
		//
		//if r.Method == "OPTIONS" {
		//	return
		//}

		log.Println("HERE")

		if r.Body != nil {
			bodyBytes, err = io.ReadAll(r.Body)
			if err != nil {
				fmt.Printf("Body reading error: %v", err)
				return
			}
			defer r.Body.Close()
		}

		_, err = w.Write(bodyBytes)
		if err != nil {
			panic(err)
		}
	}
}

func marketsAndCriterias(repo *repos.SQLCriteriaRepository) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.Header().Add("Access-Control-Allow-Methods", "*")

		type ChangeRequest struct {
			Criterias []valueobjects.Selectable `json:criterias`
			Markets   []valueobjects.Selectable `json:markets`
		}

		decoder := json.NewDecoder(r.Body)
		var cr ChangeRequest
		err := decoder.Decode(&cr)
		if err != nil {
			panic(err)
		}

		err = repo.UpdateSelectedCriterias(cr.Criterias)
		if err != nil {
			panic(err)
		}
		err = repo.UpdateSelectedMarkets(cr.Markets)
		if err != nil {
			panic(err)
		}

		type Respose struct {
			Data string
		}
		res := Respose{Data: "success"}

		response, err := json.Marshal(&res)
		if err != nil {
			panic(err)
		}

		_, err = w.Write(response)
		if err != nil {
			panic(err)
		}

	}
}

//func getPriceDistribution(repo repos.IAdsRepository) {
//	allAds, err := repo.GetAll()
//	step := 5000
//
//}

func cleanupPrices(repo repos.AdsRepository) {
	// get all ads
	allAds, err := repo.GetAll()
	if err != nil {
		panic(err)
	}
	for _, ad := range *allAds {
		prices := repo.GetAdPrices(ad.ID)
		if len(prices) > 1 {
			//log.Println("Ad has more prices ", ad.ID)
			firstPrice := prices[0].Price
			for i, price := range prices {
				if i >= 1 {
					if price.Price == firstPrice {
						repo.DeletePrice(price.ID)
					}
				}
			}
			duplicates := removeDuplicates(prices, ad.ID)
			if len(duplicates) > 0 {

				//for _, price := range prices {
				//	log.Printf("Price: %d Price ID: %d - Ad ID: %d ", price.Price, price.ID, price.AdID)
				//}

				for _, duplicatePriceID := range duplicates {
					repo.DeletePrice(duplicatePriceID)
				}
			}

		}
	}
}

func removeDuplicates(prices []adsdb.Price, adID uint) []uint {
	seen := make(map[int]bool)
	result := []uint{}
	duplicates := []uint{}

	for _, price := range prices {
		if _, ok := seen[price.Price]; !ok {
			seen[price.Price] = true
			result = append(result, price.ID)
		} else {
			duplicates = append(duplicates, price.ID)
		}
	}
	if len(duplicates) > 0 {

		log.Println("We have duplicates for ad ID: ", adID)
		log.Printf("clean : %+v", result)
		log.Printf("duplicates : %+v", duplicates)
	}
	return duplicates
}

func getAd(db *gorm.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		idStr, ok := vars["id"]
		if !ok {
			fmt.Println("id is missing in parameters")
		}
		id, err := strconv.Atoi(idStr)
		if err != nil {
			w.Write([]byte("Invalid ID"))
			return
		}
		var ad adsdb.Ad
		tx := db.Preload("Prices").First(&ad, id)
		if tx.Error != nil {
			panic(err)
		}

		response, err := json.Marshal(&ad)
		if err != nil {
			panic(err)
		}
		w.Write(response)
	}
}

func getMarkets(repo repos.IMarketsRepository) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method == http.MethodOptions {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			w.WriteHeader(http.StatusOK)
			return
		}

		// Set CORS headers for the actual request
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json")

		markets := repo.GetAll()
		type MarketsResponse struct {
			Data []adsdb.Market
		}
		marketsResponse := MarketsResponse{Data: *markets}
		response, err := json.Marshal(&marketsResponse)
		if err != nil {
			panic(err)
		}
		//w.Header().Set("Access-Control-Allow-Methods", "POST")
		//w.Header().Add("Access-Control-Allow-Origin", "*")
		w.Write(response)
	}
}

func getCriterias(repo repos.ICriteriaRepository) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method == http.MethodOptions {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			w.WriteHeader(http.StatusOK)
			return
		}

		// Set CORS headers for the actual request
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json")

		criterias := repo.GetAll()
		response, err := json.Marshal(&criterias)
		if err != nil {
			panic(err)
		}
		w.Write(response)
	}
}

func getAdsForCriteria(repo *repos.AdsRepository) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method == http.MethodOptions {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			w.WriteHeader(http.StatusOK)
			return
		}

		// Set CORS headers for the actual request
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json")

		vars := mux.Vars(r)
		idStr, ok := vars["id"]
		if !ok {
			fmt.Println("id is missing in parameters")
		}
		id, err := strconv.Atoi(idStr)
		if err != nil {
			w.Write([]byte("Invalid ID"))
			return
		}

		sortOptionDirection := r.URL.Query().Get("sortDirection")
		sortOption := r.URL.Query().Get("sortOption")
		marketsStr := r.URL.Query().Get("markets")
		limitLowStr := r.URL.Query().Get("limitLow")
		limitHighStr := r.URL.Query().Get("limitHigh")
		groupingOption := r.URL.Query().Get("groupingOption")
		yearsStr := r.URL.Query().Get("years")

		//var lowLimit *int

		low, err := strconv.Atoi(limitLowStr)
		if err != nil {
			log.Println(err)
		}
		lowLimit := &low

		high, err := strconv.Atoi(limitHighStr)
		if err != nil {
			log.Println(err)
		}
		highLimit := &high

		markets := strings.Split(marketsStr, ",")

		var years []string
		if yearsStr != "" {
			yearsArr := strings.Split(yearsStr, ",")
			for _, yearStr := range yearsArr {
				//year, err := strconv.Atoi(yearStr)
				//if err != nil {
				//	panic(err)
				//}
				years = append(years, yearStr)
			}
		}

		years = strings.Split(yearsStr, ",")

		var ads []Ad
		dbAds := repo.GetAdsForCriteria(uint(id), markets, nil, nil, lowLimit, highLimit, &years)
		//var ads []Ad
		type GroupedAds struct {
			Discounted []Ad
			Rest       []Ad
			Increased  []Ad
		}
		var groupedAds GroupedAds
		type AdsResponse struct {
			Data []Ad
		}
		var res AdsResponse
		for _, dbAd := range *dbAds {
			if !inPriceRange(lowLimit, highLimit, dbAd) {
				continue
			}
			if groupingOption == "discounted" {
				if len(dbAd.Prices) > 1 {
					discountVal, discountPercent := discount.CalculateAdDiscount(dbAd)

					if discountVal > 0 {
						groupedAds.Discounted = append(groupedAds.Discounted, Ad{
							Ad:              dbAd,
							Age:             age.CalculateAdAge(dbAd),
							DiscountValue:   discountVal,
							DiscountPercent: discountPercent,
						})
					} else {
						groupedAds.Increased = append(groupedAds.Increased, Ad{
							Ad:              dbAd,
							Age:             age.CalculateAdAge(dbAd),
							DiscountValue:   discountVal,
							DiscountPercent: discountPercent,
						})
					}

				} else {
					groupedAds.Rest = append(groupedAds.Rest, Ad{
						Ad:              dbAd,
						Age:             age.CalculateAdAge(dbAd),
						DiscountValue:   0,
						DiscountPercent: 0,
					})
				}

			}

			if groupingOption == "none" {

				discountVal := 0
				discountPercent := float64(0)
				if len(dbAd.Prices) > 1 {
					discountVal, discountPercent = discount.CalculateAdDiscount(dbAd)
				}
				ads = append(ads, Ad{
					Ad:              dbAd,
					Age:             age.CalculateAdAge(dbAd),
					DiscountValue:   discountVal,
					DiscountPercent: discountPercent,
				})

			}
		}
		if groupingOption == "discounted" {
			sortAds(groupedAds.Discounted, sortOption, sortOptionDirection)
			sortAds(groupedAds.Rest, sortOption, sortOptionDirection)
			sortAds(groupedAds.Increased, sortOption, sortOptionDirection)
			res.Data = append(res.Data, groupedAds.Discounted...)
			res.Data = append(res.Data, groupedAds.Rest...)
			res.Data = append(res.Data, groupedAds.Increased...)
		}
		if groupingOption == "none" {
			sortAds(ads, sortOption, sortOptionDirection)
			res.Data = ads
		}

		response, err := json.Marshal(&res)
		if err != nil {
			panic(err)
		}
		w.Write(response)
	}
}

func getAdsForCriteriaPaginated(repo *repos.AdsRepository) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		// Handle preflight request
		if r.Method == http.MethodOptions {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			w.WriteHeader(http.StatusOK)
			return
		}

		// Set CORS headers for the actual request
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json")

		vars := mux.Vars(r)
		idStr, ok := vars["id"]
		if !ok {
			fmt.Println("id is missing in parameters")
		}
		id, err := strconv.Atoi(idStr)
		if err != nil {
			w.Write([]byte("Invalid ID"))
			return
		}

		sortOptionDirection := r.URL.Query().Get("sortDirection")
		sortOption := r.URL.Query().Get("sortOption")
		marketsStr := r.URL.Query().Get("markets")
		limitLowStr := r.URL.Query().Get("limitLow")
		limitHighStr := r.URL.Query().Get("limitHigh")
		groupingOption := r.URL.Query().Get("groupingOption")
		limitStr := r.URL.Query().Get("limit")
		pageStr := r.URL.Query().Get("page")
		yearsStr := r.URL.Query().Get("years")

		limit, err := strconv.Atoi(limitStr)
		if err != nil {
			log.Println(err)
			limit = 50
		}

		page, err := strconv.Atoi(pageStr)
		if err != nil {
			log.Println(err)
			page = 1
		}

		//var lowLimit *int

		low, err := strconv.Atoi(limitLowStr)
		if err != nil {
			log.Println(err)
		}
		lowLimit := &low

		high, err := strconv.Atoi(limitHighStr)
		if err != nil {
			log.Println(err)
		}
		highLimit := &high

		markets := strings.Split(marketsStr, ",")

		requestPagination := repos.Pagination{
			Limit: limit,
			Page:  page,
		}

		var ads []Ad

		var years []string
		if yearsStr != "" {
			yearsArr := strings.Split(yearsStr, ",")
			for _, yearStr := range yearsArr {
				//year, err := strconv.Atoi(yearStr)
				//if err != nil {
				//	panic(err)
				//}
				years = append(years, yearStr)
			}
		}
		years = strings.Split(yearsStr, ",")

		//dbAds, pagination := repo.GetAdsForCriteriaPaginated(&requestPagination, uint(id), markets, nil, nil, lowLimit, highLimit)
		dbAds := repo.GetAdsForCriteria(uint(id), markets, nil, nil, lowLimit, highLimit, &years)

		//var ads []Ad
		type GroupedAds struct {
			Discounted []Ad
			Rest       []Ad
			Increased  []Ad
		}
		var groupedAds GroupedAds

		type Paginated struct {
			Pagination repos.Pagination
			Ads        []Ad
		}

		type AdsResponse struct {
			Data Paginated
		}
		res := AdsResponse{Data: Paginated{
			Pagination: requestPagination,
			Ads:        nil,
		}}

		for _, dbAd := range *dbAds {

			if !inPriceRange(lowLimit, highLimit, dbAd) {
				continue
			}

			if groupingOption == "discounted" {
				if len(dbAd.Prices) > 1 {
					discountVal, discountPercent := discount.CalculateAdDiscount(dbAd)
					discountDailyAmount := discount.CalculateDailyDiscount(dbAd)
					if discountVal > 0 {
						groupedAds.Discounted = append(groupedAds.Discounted, Ad{
							Ad:                    dbAd,
							Age:                   age.CalculateAdAge(dbAd),
							DiscountValue:         discountVal,
							DiscountPercent:       discountPercent,
							DailyDiscountAmmount:  discountDailyAmount,
							DealerAverageDiscount: discount.GetCalculatedAverageDealerDiscountPercent(repo, dbAd.SellerID, dbAd.MarketID),
						})
					} else {
						groupedAds.Increased = append(groupedAds.Increased, Ad{
							Ad:                    dbAd,
							Age:                   age.CalculateAdAge(dbAd),
							DiscountValue:         discountVal,
							DiscountPercent:       discountPercent,
							DealerAverageDiscount: discount.GetCalculatedAverageDealerDiscountPercent(repo, dbAd.SellerID, dbAd.MarketID),
						})
					}

				} else {
					groupedAds.Rest = append(groupedAds.Rest, Ad{
						Ad:                    dbAd,
						Age:                   age.CalculateAdAge(dbAd),
						DiscountValue:         0,
						DiscountPercent:       0,
						DealerAverageDiscount: discount.GetCalculatedAverageDealerDiscountPercent(repo, dbAd.SellerID, dbAd.MarketID),
					})
				}

			}

			if groupingOption == "none" {

				discountVal := 0
				discountPercent := float64(0)
				if len(dbAd.Prices) > 1 {
					discountVal, discountPercent = discount.CalculateAdDiscount(dbAd)
				}
				ads = append(ads, Ad{
					Ad:                    dbAd,
					Age:                   age.CalculateAdAge(dbAd),
					DiscountValue:         discountVal,
					DiscountPercent:       discountPercent,
					DealerAverageDiscount: discount.GetCalculatedAverageDealerDiscountPercent(repo, dbAd.SellerID, dbAd.MarketID),
				})

			}
		}
		if groupingOption == "discounted" {
			sortAds(groupedAds.Discounted, sortOption, sortOptionDirection)
			sortAds(groupedAds.Rest, sortOption, sortOptionDirection)
			sortAds(groupedAds.Increased, sortOption, sortOptionDirection)
			res.Data.Ads = append(res.Data.Ads, groupedAds.Discounted...)
			res.Data.Ads = append(res.Data.Ads, groupedAds.Rest...)
			res.Data.Ads = append(res.Data.Ads, groupedAds.Increased...)
			lengthOfResults := len(res.Data.Ads)
			startIndex := limit * (page - 1)
			endIndex := page * limit
			res.Data.Pagination.TotalPages = int(math.Ceil(float64(lengthOfResults) / float64(limit)))
			res.Data.Pagination.TotalRows = int64(lengthOfResults)
			if endIndex > lengthOfResults {
				endIndex = lengthOfResults
			}
			res.Data.Ads = res.Data.Ads[startIndex:endIndex]

		}
		if groupingOption == "none" {
			sortAds(ads, sortOption, sortOptionDirection)
			lengthOfResults := len(ads)
			startIndex := limit * (page - 1)
			endIndex := page * limit

			if endIndex > lengthOfResults {
				endIndex = lengthOfResults
			}
			res.Data.Ads = ads[startIndex:endIndex]
			res.Data.Pagination.TotalPages = int(math.Ceil(float64(lengthOfResults) / float64(limit)))
			res.Data.Pagination.TotalRows = int64(lengthOfResults)
		}

		response, err := json.Marshal(&res)
		if err != nil {
			panic(err)
		}

		w.Write(response)
	}
}

func sortAds(ads []Ad, sortOption string, sortOptionDirection string) {
	if sortOption == "byAge" {
		if sortOptionDirection == "desc" {
			sort.Sort(ByAgeDesc(ads))
		} else {
			sort.Sort(ByAge(ads))
		}
	}
	if sortOption == "byPrice" {
		if sortOptionDirection == "desc" {
			sort.Sort(ByPriceDesc(ads))
		} else {
			sort.Sort(ByPrice(ads))
		}
	}
	if sortOption == "byLastChanged" {
		if sortOptionDirection == "desc" {
			sort.Sort(ByLastChangeDesc(ads))
		} else {
			sort.Sort(ByLastChange(ads))
		}
	}
	if sortOption == "byDiscount" {
		if sortOptionDirection == "desc" {
			sort.Sort(ByDiscountDesc(ads))
		} else {
			sort.Sort(ByDiscount(ads))
		}
	}
	if sortOption == "byDiscountPercent" {
		if sortOptionDirection == "desc" {
			sort.Sort(ByDiscountPercentDesc(ads))
		} else {
			sort.Sort(ByDiscountPercent(ads))
		}
	}
	if sortOption == "byDailyDiscountAmount" {
		if sortOptionDirection == "desc" {
			sort.Sort(ByDailyDiscountAmountDesc(ads))
		} else {
			sort.Sort(ByDiscountPercent(ads))
		}
	}

}

type Ad struct {
	adsdb.Ad
	Age                   int
	DiscountValue         int
	DiscountPercent       float64
	DailyDiscountAmmount  float64
	DealerAverageDiscount float64
}

type ByPrice []Ad
type ByAge []Ad
type ByPriceDesc []Ad
type ByAgeDesc []Ad
type ByLastChange []Ad
type ByLastChangeDesc []Ad
type ByDiscount []Ad
type ByDiscountDesc []Ad
type ByDiscountPercent []Ad
type ByDiscountPercentDesc []Ad
type ByDailyDiscountAmountDesc []Ad
type ByDailyDiscountAmount []Ad

func sortAdsByPrice(ads *[]Ad) {
	sort.Sort(ByPrice(*ads))
}

func (a ByDailyDiscountAmountDesc) Len() int      { return len(a) }
func (a ByDailyDiscountAmountDesc) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByDailyDiscountAmountDesc) Less(i, j int) bool {
	return a[i].DailyDiscountAmmount > a[j].DailyDiscountAmmount
}

func (a ByDailyDiscountAmount) Len() int      { return len(a) }
func (a ByDailyDiscountAmount) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByDailyDiscountAmount) Less(i, j int) bool {
	return a[i].DailyDiscountAmmount < a[j].DailyDiscountAmmount
}

func (a ByDiscountPercent) Len() int      { return len(a) }
func (a ByDiscountPercent) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByDiscountPercent) Less(i, j int) bool {
	return a[i].DiscountPercent < a[j].DiscountPercent
}

func (a ByDiscountPercentDesc) Len() int      { return len(a) }
func (a ByDiscountPercentDesc) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByDiscountPercentDesc) Less(i, j int) bool {
	return a[i].DiscountPercent > a[j].DiscountPercent
}

func (a ByDiscount) Len() int      { return len(a) }
func (a ByDiscount) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByDiscount) Less(i, j int) bool {
	return a[i].DiscountValue < a[j].DiscountValue
}

func (a ByDiscountDesc) Len() int      { return len(a) }
func (a ByDiscountDesc) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByDiscountDesc) Less(i, j int) bool {
	return a[i].DiscountValue > a[j].DiscountValue
}

func (a ByLastChange) Len() int      { return len(a) }
func (a ByLastChange) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByLastChange) Less(i, j int) bool {
	lenOfIPrices := len(a[i].Ad.Prices)
	lenOfJPrices := len(a[j].Ad.Prices)
	date_i := a[i].Ad.Prices[lenOfIPrices-1].CreatedAt
	date_j := a[j].Ad.Prices[lenOfJPrices-1].CreatedAt
	return date_i.Before(date_j)
}

func (a ByLastChangeDesc) Len() int      { return len(a) }
func (a ByLastChangeDesc) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByLastChangeDesc) Less(i, j int) bool {
	lenOfIPrices := len(a[i].Ad.Prices)
	lenOfJPrices := len(a[j].Ad.Prices)
	date_i := a[i].Ad.Prices[lenOfIPrices-1].CreatedAt
	date_j := a[j].Ad.Prices[lenOfJPrices-1].CreatedAt
	return date_j.Before(date_i)
}

func (a ByPrice) Len() int      { return len(a) }
func (a ByPrice) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByPrice) Less(i, j int) bool {
	lenOfIPrices := len(a[i].Ad.Prices)
	lenOfJPrices := len(a[j].Ad.Prices)
	price_i := a[i].Ad.Prices[lenOfIPrices-1].Price
	price_j := a[j].Ad.Prices[lenOfJPrices-1].Price
	return price_i < price_j
}

func (a ByAge) Len() int      { return len(a) }
func (a ByAge) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByAge) Less(i, j int) bool {
	return a[i].Age < a[j].Age
}

func (a ByPriceDesc) Len() int      { return len(a) }
func (a ByPriceDesc) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByPriceDesc) Less(i, j int) bool {
	lenOfIPrices := len(a[i].Ad.Prices)
	lenOfJPrices := len(a[j].Ad.Prices)
	price_i := a[i].Ad.Prices[lenOfIPrices-1].Price
	price_j := a[j].Ad.Prices[lenOfJPrices-1].Price
	return price_i > price_j
}

func (a ByAgeDesc) Len() int      { return len(a) }
func (a ByAgeDesc) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByAgeDesc) Less(i, j int) bool {
	return a[i].Age > a[j].Age
}

func inPriceRange(low *int, high *int, ad adsdb.Ad) bool {
	lastPrice := ad.Prices[len(ad.Prices)-1].Price
	noHighLimit := high == nil || *high == 0
	noLowLimit := low == nil || *low == 0

	if noLowLimit && noHighLimit {
		return true
	}
	if low == nil && high == nil {
		return true
	}

	hasLowLimit := low != nil && *low > 0
	hasHighLimit := high != nil && *high > 0

	if !hasLowLimit && hasHighLimit {
		if lastPrice <= *high {
			return true
		} else {
			return false
		}
	}

	if hasLowLimit && !hasHighLimit {
		if lastPrice >= *low {
			return true
		} else {
			return false
		}
	}

	if lastPrice >= *low && lastPrice <= *high {
		return true
	}

	return false
}
