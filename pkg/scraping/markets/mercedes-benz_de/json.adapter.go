package mercedes_benz_de

import (
	"carscraper/pkg/jobs"
	"carscraper/pkg/scraping/icollector"
	"fmt"
	"log"
	"math"
	"strconv"
)

type MercedesBenzDEAdapter struct {
	request        *Request
	requestBuilder *RequestBuilder
	urlBuilder     *MercedesBenzRoURLBuilder
}

func NewMercedesBenzDEAdapter() *MercedesBenzDEAdapter {
	r := NewRequest()
	rb := NewRequestBuilder()
	b := NewMercedesBenzRoURLBuilder()
	return &MercedesBenzDEAdapter{
		request:        r,
		requestBuilder: rb,
		urlBuilder:     b,
	}
}

func (a MercedesBenzDEAdapter) GetAds(job jobs.SessionJob) icollector.AdsResults {
	url := a.urlBuilder.GetURL(job)

	requestBody := a.requestBuilder.GetRequestBody(job)
	response, err := a.request.MakeRequest(url, requestBody)
	if err != nil {
		return icollector.AdsResults{
			Ads:        nil,
			IsLastPage: true,
			Error:      err,
		}
	}

	if response.StatusInfo.Status.Code == "400" && response.StatusInfo.Status.Text == "pageIndex too high" {
		return icollector.AdsResults{
			Ads:        nil,
			IsLastPage: true,
			Error:      nil,
		}
	}

	return a.processResponse(*response, job.Criteria)
}

func (a MercedesBenzDEAdapter) processResponse(response Response, criteria jobs.Criteria) icollector.AdsResults {
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
				Ad_url:             fmt.Sprintf("https://www.mercedes-benz.de/passengercars/buy/used-car-search.html/u/used-vehicles/d/details?id=%s", vehicle.Id),
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
