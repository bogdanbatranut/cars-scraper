package oferte_bmw

import (
	"carscraper/pkg/jobs"
	"carscraper/pkg/scraping/icollector"
	"fmt"
	"log"
	"strconv"
)

type OferteBMWAdapter struct {
}

func NewOferteBMWAdapter() *OferteBMWAdapter {
	return &OferteBMWAdapter{}
}

func (a OferteBMWAdapter) GetAds(job jobs.SessionJob) icollector.AdsResults {
	request := NewRequest()
	urlBuilder := NewOferteBMWURLBuilder()
	url := urlBuilder.GetURL(job)
	//url := "https://oferte.bmw.ro/rulate/api/v1/ems/bmw-used-ro_RO/search"

	requestBody := createRequestBody(job)

	data, err := request.DoPOSTRequest(*url, requestBody)
	if err != nil {
		log.Println(err)
		return icollector.AdsResults{
			Ads:        nil,
			IsLastPage: true,
			Error:      err,
		}
	}
	ads := a.getAds(*data, job.Criteria)
	return icollector.AdsResults{
		Ads:        ads,
		IsLastPage: true,
		Error:      nil,
	}
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
