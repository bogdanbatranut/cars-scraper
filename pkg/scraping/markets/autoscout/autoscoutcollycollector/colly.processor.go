package autoscoutcollycollector

import (
	"carscraper/pkg/jobs"
	"carscraper/pkg/scraping/icollector"
	"fmt"
	"log"
	"math"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
)

type AutoScoutCollyProcessor struct {
}

func NewAutoscoutCollyProcessor() *AutoScoutCollyProcessor {
	return &AutoScoutCollyProcessor{}
}

func (collector AutoScoutCollyProcessor) GetAds(url string, pageNumber int, criteria jobs.Criteria) icollector.AdsResults {

	foundAds := []jobs.Ad{}

	c := colly.NewCollector()
	isLastPage := false

	var executionErr error

	var totalResults float64

	c.OnHTML("#__next > div > div > div.ListPage_wrapper__vFmTi > div.ListPage_container__Optya > main > header > div.ListHeader_top__N6YWA > h1 > span > span:nth-child(1)", func(e *colly.HTMLElement) {
		totalResultsStr := strings.Replace(e.Text, ".", "", -1)
		totalResults_, err := strconv.Atoi(totalResultsStr)
		if err != nil {
			executionErr = err
			return
		}
		totalResults = float64(totalResults_)

		numberOfTotalPages := math.Ceil(totalResults / 20)

		//log.Printf("Number fo total pages %.4f current page %d from total results %d", numberOfTotalPages, pageNumber, totalResults_)

		if float64(pageNumber) == numberOfTotalPages || totalResults == 0 {
			isLastPage = true
		}

	})

	c.OnHTML("#__next > div > div > div.ListPage_wrapper__vFmTi > div.ListPage_container__Optya > main > div.ListPage_pagination__4Vw9q > nav > ul > li:last-child", func(element *colly.HTMLElement) {
		_, disabled := element.DOM.Find("button").Attr("disabled")
		isLastPage = disabled
	})

	if executionErr != nil {
		return icollector.AdsResults{
			Ads:        nil,
			IsLastPage: false,
			Error:      executionErr,
		} // nil, false, executionErr
	}

	c.OnHTML("article", func(e *colly.HTMLElement) {
		sellerType := "dealer"

		sellerFAttr := e.Attr("data-seller-type")
		if sellerFAttr != "d" {
			sellerType = "privat"
		}

		adId := e.Attr("id")

		adHref, exists := e.DOM.Find("div > div.ListItem_header__J6xlG.ListItem_header_new_design__Rvyv_ > a").Attr("href")
		if !exists {
			adHref = "NOT FOUND!!"
		}

		thumbnail, exists := e.DOM.Find("div.ListItem_wrapper__TxHWu > div.Gallery_wrapper__iqp3u > section > div:nth-child(1) > picture > img").Attr("src")
		if !strings.HasPrefix(thumbnail, "https://prod.pictures.autoscout24.net/listing-images/") || !exists {
			thumbnail, exists = e.DOM.Find("div.ListItem_wrapper__TxHWu > div.Gallery_wrapper__iqp3u > section > div:nth-child(1) > picture > source:nth-child(1)").Attr("srcset")
		}

		yearStr := e.Attr("data-first-registration")
		year, err := strconv.Atoi(yearStr[3:])
		if err != nil {
			executionErr = err
			return
		}

		kmStr := e.Attr("data-mileage")
		km, err := strconv.Atoi(kmStr)

		if err != nil {
			executionErr = err
			return
		}

		grossPriceStr := e.Attr("data-price")
		grossPrice, err := strconv.Atoi(grossPriceStr)

		if err != nil {
			executionErr = err
			return
		}

		seller := "autoscout24.de"
		carad := jobs.Ad{
			Brand:              criteria.Brand,
			Model:              criteria.CarModel,
			Year:               year,
			Km:                 km,
			Fuel:               criteria.Fuel,
			Price:              grossPrice,
			AdID:               adId,
			Ad_url:             fmt.Sprintf("https://www.autoscout24.ro%s", adHref),
			SellerType:         sellerType,
			SellerName:         &seller,
			SellerNameInMarket: &seller,
			SellerOwnURL:       &seller,
			SellerMarketURL:    &seller,
			Thumbnail:          &thumbnail,
		}
		foundAds = append(foundAds, carad)
	})

	if executionErr != nil {
		return icollector.AdsResults{
			Ads:        nil,
			IsLastPage: false,
			Error:      executionErr,
		}
		//nil, false, executionErr
	}

	c.OnRequest(func(request *colly.Request) {
		request.Headers.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36")
	})

	err := c.Visit(url)
	log.Println("AUTOSCOUT Visiting ", url)
	if err != nil {

		return icollector.AdsResults{
			Ads:        nil,
			IsLastPage: false,
			Error:      err,
		}
		//nil, false, err
	}
	c.Wait()
	if len(foundAds) == 0 {
		log.Println("WE NO RESULTS SO RETURN !!!!!")
		return icollector.AdsResults{
			Ads:        nil,
			IsLastPage: true,
			Error:      nil,
		} //nil, true, nil
	}
	log.Println("AUTOSCOUT found ads : ", len(foundAds))
	//return foundAds, isLastPage, nil
	return icollector.AdsResults{
		Ads:        &foundAds,
		IsLastPage: isLastPage,
		Error:      nil,
	}
}
