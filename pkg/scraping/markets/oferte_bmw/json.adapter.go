package oferte_bmw

import (
	"carscraper/pkg/jobs"
	"carscraper/pkg/logging"
	"carscraper/pkg/scraping/icollector"
	"fmt"
	"log"
	"strconv"
)

type OferteBMWAdapter struct {
	loggingService *logging.ScrapeLoggingService
}

func NewOferteBMWAdapter(loggingService *logging.ScrapeLoggingService) *OferteBMWAdapter {
	return &OferteBMWAdapter{
		loggingService: loggingService,
	}
}

func (a OferteBMWAdapter) GetAds(job jobs.SessionJob) icollector.AdsResults {
	request := NewRequest()

	criteriaLog, err := a.loggingService.GetCriteriaLog(job.SessionID, job.CriteriaID, job.MarketID)
	if err != nil {
		panic(err)
	}
	pageLog, err := a.loggingService.CreatePageLog(criteriaLog, job, "", job.Market.PageNumber)
	if err != nil {
		panic(err)
	}

	urlBuilder := NewOferteBMWURLBuilder()
	url := urlBuilder.GetURL(job)
	//url := "https://oferte.bmw.ro/rulate/api/v1/ems/bmw-used-ro_RO/search"

	err = a.loggingService.PageLogSetVisitURL(pageLog, *url)
	if err != nil {
		log.Println(err.Error())
	}

	requestBody := createRequestBody(job)

	data, err := request.DoPOSTRequest(*url, requestBody)

	adsResults := icollector.AdsResults{
		Ads:        nil,
		IsLastPage: false,
		Error:      nil,
	}

	if err != nil {
		log.Println(err)
		adsResults.IsLastPage = true
		adsResults.Error = err
		return adsResults
	}

	ads := a.getAds(*data, job.Criteria)
	adsResults.Ads = ads
	adsResults.IsLastPage = true

	err2 := a.loggingService.PageLogSetPageScraped(pageLog, len(*adsResults.Ads), adsResults.IsLastPage)

	if err2 != nil {
		log.Println(err2.Error())
	}

	return adsResults
}

func (a OferteBMWAdapter) getAds(response OferteBMWResponse, criteria jobs.Criteria) *[]jobs.Ad {
	oferteAds := response.List
	ads := []jobs.Ad{}
	calculated704Param := get704paramForThumbnailURL()

	for _, ad := range oferteAds {
		ad := jobs.Ad{
			Title:              &ad.Title,
			Brand:              criteria.Brand,
			Model:              criteria.CarModel,
			Year:               ad.ProductionYear,
			Km:                 ad.Mileage,
			Fuel:               criteria.Fuel,
			Price:              ad.TransactionalPrice,
			AdID:               strconv.Itoa(ad.Id),
			Ad_url:             createAdURL(ad.Id),
			SellerType:         "dealer",
			SellerName:         &ad.Dealer.Name,
			SellerNameInMarket: &ad.Dealer.Name,
			SellerOwnURL:       &ad.Dealer.Name,
			SellerMarketURL:    &ad.Dealer.Name,
			Thumbnail:          getThumbNail(ad.Id, *calculated704Param),
		}
		ads = append(ads, ad)
	}
	return &ads
}

func createAdURL(adID int) string {
	return fmt.Sprintf("https://oferte.bmw.ro/rulate/cauta/detaliu/%d/", adID)
}

func createRequestBody(job jobs.SessionJob) RequestBody {
	requestPayload := NewRequestPayload()
	brandSeriesVariants := requestPayload.GetIds(job.Criteria)

	body := RequestBody{
		Match: Match{
			TransactionalPrice: TransactionalPrice{
				Min: 0,
				Max: 999999,
			},
			Brand:        brandSeriesVariants.Brand,
			Series:       brandSeriesVariants.Series,
			Variant:      brandSeriesVariants.Variant,
			Registration: Registration{Min: fmt.Sprintf("%d-01-01", *job.Criteria.YearFrom)},
			Mileage:      Mileage{Max: *job.Criteria.KmTo},
			Fuel:         requestPayload.GetFuelID(job.Criteria.Fuel),
		},
		Skip:  0,
		Limit: 2000,
		Sort:  nil,
	}
	return body
}
