package mercedes_benz_ro

import (
	"carscraper/pkg/jobs"
	"carscraper/pkg/logging"
	"carscraper/pkg/scraping/icollector"
	"fmt"
	"log"
	"math"
	"strconv"
)

type MercedesBenzRoAdapter struct {
	request        *Request
	requestBuilder *RequestBuilder
	urlBuilder     *MercedesBenzRoURLBuilder
	loggingService *logging.ScrapeLoggingService
}

func NewMercedesBenzRoAdapter(logingService *logging.ScrapeLoggingService) *MercedesBenzRoAdapter {
	r := NewRequest()
	rb := NewRequestBuilder()
	b := NewMercedesBenzRoURLBuilder()
	return &MercedesBenzRoAdapter{
		request:        r,
		requestBuilder: rb,
		urlBuilder:     b,
		loggingService: logingService,
	}
}

func (a MercedesBenzRoAdapter) GetAds(job jobs.SessionJob) icollector.AdsResults {

	criteriaLog, err := a.loggingService.GetCriteriaLog(job.SessionID, job.CriteriaID, job.MarketID)
	if err != nil {
		panic(err)
	}
	pageLog, err := a.loggingService.CreatePageLog(criteriaLog, job, "", job.Market.PageNumber)
	if err != nil {
		panic(err)
	}

	url := a.urlBuilder.GetURL(job)

	err = a.loggingService.PageLogSetVisitURL(pageLog, url)
	if err != nil {
		log.Println(err.Error())
	}

	requestBody := a.requestBuilder.GetRequestBody(job)
	response, err := a.request.MakeRequest(url, requestBody)
	if err != nil {
		return icollector.AdsResults{
			Ads:        nil,
			IsLastPage: true,
			Error:      err,
		}
	}

	adsResults := a.processResponse(*response, job.Criteria)

	err2 := a.loggingService.PageLogSetPageScraped(pageLog, len(*adsResults.Ads), adsResults.IsLastPage)

	if err2 != nil {
		log.Println(err2.Error())
	}

	return adsResults
}

func (a MercedesBenzRoAdapter) processResponse(response Response, criteria jobs.Criteria) icollector.AdsResults {
	var foundAds []jobs.Ad
	isLastPage := false
	for _, result := range response.Results {
		if result.Summary.Quantity < 12 {
			isLastPage = true
		}
		for _, vehicle := range result.Vehicles {
			ad := jobs.Ad{
				Title:              &vehicle.VehicleConfiguration.SalesDescription,
				Brand:              criteria.Brand,
				Model:              criteria.CarModel,
				Year:               getYear(vehicle.VehicleConfiguration.ModelYear),
				Km:                 vehicle.VehicleCondition.Mileage,
				Fuel:               criteria.Fuel,
				Price:              int(math.Round(vehicle.PriceInformation.OfferPrice.TotalPrice)),
				AdID:               vehicle.Id,
				Ad_url:             fmt.Sprintf("https://www.mercedes-benz.ro/passengercars/buy/used-car-search.html/u/used-vehicles/d/details?id=%s", vehicle.Id),
				SellerType:         "dealer",
				SellerName:         &vehicle.VehicleLocation.Formatted.Nameline1,
				SellerNameInMarket: &vehicle.VehicleLocation.Formatted.Nameline1,
				SellerOwnURL:       &vehicle.VehicleLocation.Formatted.Links.Website,
				SellerMarketURL:    &vehicle.VehicleLocation.Formatted.Links.Website,
				Thumbnail:          getThumbNailImage(vehicle.Media.Images),
			}
			log.Println("AD :")
			log.Println(*ad.Title)
			foundAds = append(foundAds, ad)
		}
		log.Printf("----------> Qty on page %d page INDEX : %d TOTAL QTY: %d", result.Summary.Quantity, result.Summary.PageIndex, result.Summary.TotalQuantity)
	}
	log.Printf("Found: %d ads", len(foundAds))
	return icollector.AdsResults{
		Ads:        &foundAds,
		IsLastPage: isLastPage,
		Error:      nil,
	}
}

func getYear(yearStr string) int {
	year := 0
	if yearStr != "" {
		yearInt, err := strconv.Atoi(yearStr)
		if err != nil {
			panic(err)
		}
		year = yearInt
	}
	return year
}

func getThumbNailImage(images []Image) *string {
	var src string
	for _, image := range images {
		if image.Format == "large" && image.Perspective != "RIM" && image.Perspective != "INTERIOR" {
			src = image.Url
			break
		}
	}
	return &src
}
