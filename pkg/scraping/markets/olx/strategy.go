package olx

import (
	"carscraper/pkg/jobs"
	"carscraper/pkg/logging"
	"carscraper/pkg/scraping/icollector"
	"encoding/json"
	"strconv"
)

type OlxStrategy struct {
	loggingService *logging.ScrapeLoggingService
}

func NewOlxStrategy(logginService *logging.ScrapeLoggingService) OlxStrategy {
	return OlxStrategy{loggingService: logginService}
}

func (s OlxStrategy) Execute(job jobs.SessionJob) icollector.AdsResults {
	var ads *[]jobs.Ad

	request := NewRequest()
	urlBuilder := NewURLBuilder(job.Criteria)
	url := urlBuilder.GetPageURL(job.Market.PageNumber)

	if url == nil {
		//return ads, true, nil
		return icollector.AdsResults{
			Ads:        ads,
			IsLastPage: true,
			Error:      nil,
		}
	} else {
		data, err := request.GetPage(*url)
		if err != nil {
			//return nil, true, err
			return icollector.AdsResults{
				Ads:        nil,
				IsLastPage: true,
				Error:      err,
			}
		}
		response, err := s.toStruct(data)
		if err != nil {
			//return ads, true, err
			return icollector.AdsResults{
				Ads:        ads,
				IsLastPage: true,
				Error:      err,
			}
		}
		ads = s.getAds(*response, job.Criteria)
		if response.Links.Next == nil {
			//return ads, true, nil
			return icollector.AdsResults{
				Ads:        ads,
				IsLastPage: true,
				Error:      nil,
			}
		} else {
			//return ads, false, nil
			return icollector.AdsResults{
				Ads:        ads,
				IsLastPage: false,
				Error:      nil,
			}
		}
	}
	//return ads, true, nil
	return icollector.AdsResults{
		Ads:        ads,
		IsLastPage: true,
		Error:      nil,
	}
}

func (s OlxStrategy) getAds(response OlxResponse, criteria jobs.Criteria) *[]jobs.Ad {
	olxAds := response.Data
	ads := []jobs.Ad{}
	for _, olxAd := range olxAds {
		if olxAd.Partner != nil {
			continue
		}
		ad := jobs.Ad{
			Brand:              criteria.Brand,
			Model:              criteria.CarModel,
			Year:               olxAd.getYear(),
			Km:                 olxAd.getKm(),
			Fuel:               criteria.Fuel,
			Price:              olxAd.getPrice(),
			AdID:               strconv.Itoa(olxAd.Id),
			Ad_url:             olxAd.Url,
			SellerType:         olxAd.getSellerType(),
			SellerName:         &olxAd.User.Name,
			SellerNameInMarket: &olxAd.User.Name,
			SellerOwnURL:       &olxAd.User.Name,
			SellerMarketURL:    &olxAd.User.Name,
			Thumbnail:          olxAd.getThumbnailURL(),
		}
		ads = append(ads, ad)
	}
	return &ads
}

func (s OlxStrategy) toStruct(bytes []byte) (*OlxResponse, error) {
	var res OlxResponse
	err := json.Unmarshal(bytes, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}
