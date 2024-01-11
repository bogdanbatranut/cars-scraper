package main

import (
	"carscraper/pkg/adsdb"
	"carscraper/pkg/amconfig"
	"carscraper/pkg/errorshandler"
	"carscraper/pkg/repos"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	log.Println("starting BACKEND service...")

	cfg, err := amconfig.NewViperConfig()
	errorshandler.HandleErr(err)

	r := mux.NewRouter().StrictSlash(true)

	criteriaRepo := repos.NewSQLCriteriaRepository(cfg)
	marketsRepo := repos.NewSQLMarketsRepository(cfg)
	adsRepo := repos.NewAdsRepository(cfg)
	chartsRepo := repos.NewChartsRepository(cfg)
	chartsRepo.GetAdsPricesByStep(5000)
	////cleanupPrices(adsRepo)
	//return

	r.HandleFunc("/markets", getMarkets(marketsRepo)).Methods("GET")
	r.HandleFunc("/criterias", getCriterias(criteriaRepo)).Methods("GET")
	r.HandleFunc("/adsforcriteria/{id}", getAdsForCriteria(adsRepo)).Methods("GET")
	httpPort := cfg.GetString(amconfig.BackendServiceHTTPPort)
	log.Printf("HTTP listening on port %s\n", httpPort)
	err = http.ListenAndServe(fmt.Sprintf(":%s", httpPort), r)
	errorshandler.HandleErr(err)

}

//func getPriceDistribution(repo repos.IAdsRepository) {
//	allAds, err := repo.GetAll()
//	step := 5000
//
//}

func cleanupPrices(repo repos.IAdsRepository) {
	// get all ads
	allAds, err := repo.GetAll()
	if err != nil {
		panic(err)
	}
	for _, ad := range *allAds {
		prices := repo.GetAdPrices(ad.ID)
		if len(prices) > 1 {
			log.Println("Ad has more prices")
			firstPrice := prices[0].Price
			for i, price := range prices {
				if i >= 1 {
					if price.Price == firstPrice {
						log.Println("Deleting price ...")
						repo.DeletePrice(price.ID)
					}
				}
			}
			duplicates := removeDuplicates(prices)
			if len(duplicates) > 0 {

				for _, price := range prices {
					log.Printf("Price: %d Price ID: %d - Ad ID: %d ", price.Price, price.ID, price.AdID)
				}

				for _, duplicatePriceID := range duplicates {
					repo.DeletePrice(duplicatePriceID)
				}
			}

		}
	}
}

func removeDuplicates(prices []adsdb.Price) []uint {
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
		log.Println("we have duplicates")
		log.Printf("clean : %+v", result)
		log.Printf("duplicates : %+v", duplicates)
	}
	return duplicates
}

func getMarkets(repo repos.IMarketsRepository) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		markets := repo.GetAll()
		type MarketsResponse struct {
			Data []adsdb.Market
		}
		marketsResponse := MarketsResponse{Data: *markets}
		response, err := json.Marshal(&marketsResponse)
		if err != nil {
			panic(err)
		}
		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.Write(response)
	}
}

func getCriterias(repo repos.ICriteriaRepository) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		criterias := repo.GetAll()
		response, err := json.Marshal(&criterias)
		if err != nil {
			panic(err)
		}
		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.Write(response)
	}
}

func getAdsForCriteria(repo repos.IAdsRepository) func(w http.ResponseWriter, r *http.Request) {
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

		sortOptionDirection := r.URL.Query().Get("sortDirection")
		sortOption := r.URL.Query().Get("sortOption")
		marketsStr := r.URL.Query().Get("markets")
		limitLowStr := r.URL.Query().Get("limitLow")
		limitHighStr := r.URL.Query().Get("limitHigh")
		groupingOption := r.URL.Query().Get("groupingOption")

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

		var ads []Ad
		dbAds := repo.GetAdsForCriteria(uint(id), markets)
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
					discountVal, discountPercent := computeDiscount(dbAd)
					if discountVal > 0 {
						groupedAds.Discounted = append(groupedAds.Discounted, Ad{
							Ad:              dbAd,
							Age:             computeAge(dbAd),
							DiscountValue:   discountVal,
							DiscountPercent: discountPercent,
						})
					} else {
						groupedAds.Increased = append(groupedAds.Discounted, Ad{
							Ad:              dbAd,
							Age:             computeAge(dbAd),
							DiscountValue:   discountVal,
							DiscountPercent: discountPercent,
						})
					}

				} else {
					groupedAds.Rest = append(groupedAds.Rest, Ad{
						Ad:              dbAd,
						Age:             computeAge(dbAd),
						DiscountValue:   0,
						DiscountPercent: 0,
					})
				}

			}

			if groupingOption == "none" {

				discountVal := 0
				discountPercent := float64(0)
				if len(dbAd.Prices) > 1 {
					discountVal, discountPercent = computeDiscount(dbAd)
				}
				ads = append(ads, Ad{
					Ad:              dbAd,
					Age:             computeAge(dbAd),
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
		w.Header().Add("Access-Control-Allow-Origin", "*")
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

}

func computeAge(ad adsdb.Ad) int {
	currentTime := time.Now()
	adFirstSeenTime := ad.CreatedAt
	diff := currentTime.Sub(adFirstSeenTime)
	return int(diff.Hours() / 24)
}

func computeDiscount(ad adsdb.Ad) (int, float64) {
	discVal := ad.Prices[0].Price - ad.Prices[len(ad.Prices)-1].Price
	discPercent := float64(discVal) / float64(ad.Prices[0].Price) * 100
	ro := toFixed(discPercent, 2)
	return discVal, ro
}

func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

func toFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(round(num*output)) / output
}

type Ad struct {
	adsdb.Ad
	Age             int
	DiscountValue   int
	DiscountPercent float64
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

func sortAdsByPrice(ads *[]Ad) {
	sort.Sort(ByPrice(*ads))
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
