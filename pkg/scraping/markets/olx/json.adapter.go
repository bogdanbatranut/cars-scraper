package olx

import (
	"carscraper/pkg/jobs"
	"carscraper/pkg/logging"
	"carscraper/pkg/scraping/icollector"
	"encoding/json"
	"log"
	"strconv"
)

type OLXJSONAdapter struct {
	loggingService *logging.ScrapeLoggingService
}

func NewOLXJSONAdapter(loggingService *logging.ScrapeLoggingService) *OLXJSONAdapter {
	return &OLXJSONAdapter{
		loggingService: loggingService,
	}
}

func (a OLXJSONAdapter) GetAds(job jobs.SessionJob) icollector.AdsResults {
	var ads *[]jobs.Ad

	criteriaLog, err := a.loggingService.GetCriteriaLog(job.SessionID, job.CriteriaID, job.MarketID)
	if err != nil {
		panic(err)
	}
	pageLog, err := a.loggingService.CreatePageLog(criteriaLog, job, "", job.Market.PageNumber)
	if err != nil {
		panic(err)
	}

	request := NewRequest()
	urlBuilder := NewURLBuilder(job.Criteria)
	url := urlBuilder.GetPageURL(job.Market.PageNumber)

	err = a.loggingService.PageLogSetVisitURL(pageLog, *url)
	if err != nil {
		log.Println(err.Error())
	}

	adsResults := icollector.AdsResults{
		Ads:        nil,
		IsLastPage: false,
		Error:      nil,
	}

	if url == nil {
		adsResults.Ads = ads
		//return icollector.AdsResults{
		//	Ads:        ads,
		//	IsLastPage: true,
		//	Error:      nil,
		//}
	} else {
		data, err := request.GetPage(*url)
		if err != nil {
			//return nil, true, err
			adsResults.IsLastPage = true
			adsResults.Error = err
			//return icollector.AdsResults{
			//	Ads:        nil,
			//	IsLastPage: true,
			//	Error:      err,
			//}
		}
		response, err := a.toStruct(data)
		if err != nil {
			adsResults.Ads = ads
			adsResults.IsLastPage = true
			//return icollector.AdsResults{
			//	Ads:        ads,
			//	IsLastPage: true,
			//	Error:      err,
			//}
		}
		ads = a.getAds(*response, job.Criteria)
		if response.Links.Next == nil {
			adsResults.Ads = ads
			adsResults.IsLastPage = true

			//return icollector.AdsResults{
			//	Ads:        ads,
			//	IsLastPage: true,
			//	Error:      nil,
			//}
		} else {
			adsResults.Ads = ads
			adsResults.IsLastPage = false
			//return icollector.AdsResults{
			//	Ads:        ads,
			//	IsLastPage: false,
			//	Error:      nil,
			//}
		}
	}

	adsResults.Ads = ads
	adsResults.IsLastPage = true
	//return ads, true, nil

	err2 := a.loggingService.PageLogSetPageScraped(pageLog, len(*adsResults.Ads), adsResults.IsLastPage)

	if err2 != nil {
		log.Println(err2.Error())
	}

	return adsResults
}

func (a OLXJSONAdapter) toStruct(bytes []byte) (*OlxResponse, error) {
	var res OlxResponse
	err := json.Unmarshal(bytes, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (a OLXJSONAdapter) getAds(response OlxResponse, criteria jobs.Criteria) *[]jobs.Ad {
	olxAds := response.Data
	ads := []jobs.Ad{}
	for _, olxAd := range olxAds {
		if olxAd.Partner != nil {
			continue
		}
		ad := jobs.Ad{
			Title:              &olxAd.Title,
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
