package main

import (
	"carscraper/pkg/adsdb"
	"carscraper/pkg/amconfig"
	"carscraper/pkg/errorshandler"
	"carscraper/pkg/repos"
	"encoding/json"
	"fmt"
	"log"
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

	r.HandleFunc("/markets", getMarkets(marketsRepo)).Methods("GET")
	r.HandleFunc("/criterias", getCriterias(criteriaRepo)).Methods("GET")
	r.HandleFunc("/adsforcriteria/{id}", getAdsForCriteria(adsRepo)).Methods("GET")
	httpPort := cfg.GetString(amconfig.BackendServiceHTTPPort)
	log.Printf("HTTP listening on port %s\n", httpPort)
	err = http.ListenAndServe(fmt.Sprintf(":%s", httpPort), r)
	errorshandler.HandleErr(err)

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

		sortOptionDirection := r.URL.Query().Get("sortOptionDirection")
		sortOption := r.URL.Query().Get("sortOption")
		marketsStr := r.URL.Query().Get("markets")
		markets := strings.Split(marketsStr, ",")
		log.Println(markets)

		dbAds := repo.GetAdsForCriteria(uint(id), markets)
		var ads []Ad
		for _, dbAd := range *dbAds {
			ads = append(ads, Ad{
				Ad:  dbAd,
				Age: computeAge(dbAd),
			})
		}

		type AdsResponse struct {
			Data []Ad
		}

		if sortOption == "byAge" {
			if sortOptionDirection == "desc" {
				sort.Sort(ByAgeDesc(ads))
			}
			sort.Sort(ByAge(ads))
		} else {
			if sortOptionDirection == "desc" {
				sort.Sort(ByPriceDesc(ads))
			}
			sort.Sort(ByPrice(ads))
		}

		res := AdsResponse{Data: ads}

		response, err := json.Marshal(&res)
		if err != nil {
			panic(err)
		}
		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.Write(response)
	}
}

func computeAge(ad adsdb.Ad) int {
	currentTime := time.Now()
	adFirstSeenTime := ad.CreatedAt
	diff := currentTime.Sub(adFirstSeenTime)
	return int(diff.Hours() / 24)
}

type Ad struct {
	adsdb.Ad
	Age int
}

type ByPrice []Ad
type ByAge []Ad
type ByPriceDesc []Ad
type ByAgeDesc []Ad

func sortAdsByPrice(ads *[]Ad) {
	sort.Sort(ByPrice(*ads))
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
