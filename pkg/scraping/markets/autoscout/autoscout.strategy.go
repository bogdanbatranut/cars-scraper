package autoscout

import (
	"carscraper/pkg/jobs"
	"log"
	"math"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
)

type MobileDeStrategy struct {
}

func NewAutoscoutStrategy() MobileDeStrategy {
	return MobileDeStrategy{}
}

func (as MobileDeStrategy) Execute(job jobs.SessionJob) ([]jobs.Ad, bool, error) {
	builder := NewURLBuilder(job.Criteria)
	url := builder.GetPageURL(job.Market.PageNumber)
	ads, isLastPage, err := getData(url, job.Market.PageNumber, job.Criteria)
	if err != nil {
		return nil, false, err
	}

	//isLastPage = true
	return ads, isLastPage, nil
}

func getData(url string, pageNumber int, criteria jobs.Criteria) ([]jobs.Ad, bool, error) {

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
		return nil, false, executionErr
	}

	c.OnHTML("article.cldt-summary-full-item.listing-impressions-tracking.list-page-item.false.ListItem_article__qyYw7", func(e *colly.HTMLElement) {
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
			Ad_url:             adHref,
			SellerType:         sellerType,
			SellerName:         &seller,
			SellerNameInMarket: &seller,
			SellerOwnURL:       &seller,
			SellerMarketURL:    &seller,
		}
		foundAds = append(foundAds, carad)
	})

	if executionErr != nil {
		return nil, false, executionErr
	}

	err := c.Visit(url)
	log.Println("Visiting ", url)
	if err != nil {
		return nil, false, err
	}
	c.Wait()
	if len(foundAds) == 0 {
		log.Println("WE NO RESULTS SO RETURN !!!!!")
		return nil, true, nil
	}

	return foundAds, isLastPage, nil
}
