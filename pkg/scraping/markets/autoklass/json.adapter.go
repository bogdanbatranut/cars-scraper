package autoklass

import (
	"carscraper/pkg/jobs"
	"carscraper/pkg/logging"
	"carscraper/pkg/scraping/icollector"
	"fmt"
	"log"
	"strconv"
	"time"
)

type AutoklassRoAdapter struct {
	request        *Request
	requestBuilder *RequestBuilder
	urlBuilder     *AutoklassRoURLBuilder
	namingMapper   *AutoklassRONamingMapper
	loggingService *logging.ScrapeLoggingService
}

func NewAutoklassRoAdapter(loggingService *logging.ScrapeLoggingService) *AutoklassRoAdapter {
	r := NewRequest()
	rb := NewRequestBuilder()
	b := NewAutoklassRoURLBuilder()
	nm := NewAutoklassRoNamingMapper()
	return &AutoklassRoAdapter{
		request:        r,
		requestBuilder: rb,
		urlBuilder:     b,
		namingMapper:   nm,
		loggingService: loggingService,
	}
}

func (a AutoklassRoAdapter) GetAds(job jobs.SessionJob) icollector.AdsResults {
	criteriaLog, err := a.loggingService.GetCriteriaLog(job.SessionID, job.CriteriaID, job.MarketID)
	if err != nil {
		panic(err)
	}
	pageLog, err := a.loggingService.CreatePageLog(criteriaLog, job, "", job.Market.PageNumber)
	if err != nil {
		panic(err)
	}
	url := a.urlBuilder.GetURL(job, *a.namingMapper)
	err = a.loggingService.PageLogSetVisitURL(pageLog, url)
	if err != nil {
		panic(err)
	}
	response, err := a.request.MakeRequest(url)
	if err != nil {
		return icollector.AdsResults{
			Ads:        nil,
			IsLastPage: true,
			Error:      err,
		}
	}
	scrapeResults := a.processResponse(*response, job.Criteria)

	err2 := a.loggingService.PageLogSetPageScraped(pageLog, len(*scrapeResults.Ads), scrapeResults.IsLastPage)
	if err2 != nil {
		log.Println(err2.Error())
	}
	return scrapeResults
}

func (a AutoklassRoAdapter) processResponse(response Response, criteria jobs.Criteria) icollector.AdsResults {
	var foundAds []jobs.Ad

	for _, ad := range response.Response {
		t := time.Unix(ad.DateManufacture, 0)
		year := t.Year()
		seller := "Autoklass"
		sellerURL := "https://www.autoklass.ro/"
		var thumnbNail string
		if len(ad.CarsGallery) > 0 {
			thumnbNail = fmt.Sprintf("https://www.autoklass.ro/%s", ad.CarsGallery[0].CarMediaURL)
		}
		marketAd := jobs.Ad{
			Title:              &ad.Title,
			Brand:              criteria.Brand,
			Model:              criteria.CarModel,
			Year:               year,
			Km:                 ad.Km,
			Fuel:               criteria.Fuel,
			Price:              ad.SalePrice,
			AdID:               strconv.Itoa(ad.Id),
			Ad_url:             fmt.Sprintf("https://www.autoklass.ro/masini/mercedes-benz/%s/%s", ad.Model.Name, ad.Slug),
			SellerType:         "dealer",
			SellerName:         &seller,
			SellerNameInMarket: &seller,
			SellerOwnURL:       &sellerURL,
			SellerMarketURL:    &sellerURL,
			Thumbnail:          &thumnbNail,
		}
		log.Println(ad.Title)
		foundAds = append(foundAds, marketAd)
	}

	log.Printf("Found: %d ads", len(foundAds))
	return icollector.AdsResults{
		Ads:        &foundAds,
		IsLastPage: true,
		Error:      nil,
	}
}
