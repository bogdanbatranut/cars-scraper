package bmw_de

import (
	"carscraper/pkg/jobs"
	"carscraper/pkg/logging"
	"carscraper/pkg/scraping/icollector"
	"errors"
	"log"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
)

type BMWDECollyMarketAdapter struct {
	loggingService *logging.ScrapeLoggingService
}

func NewBMWDECollyMarketAdapter(logingService *logging.ScrapeLoggingService) *BMWDECollyMarketAdapter {
	return &BMWDECollyMarketAdapter{
		loggingService: logingService,
	}
}

func (a BMWDECollyMarketAdapter) GetAds(job jobs.SessionJob) icollector.AdsResults {
	criteriaLog, err := a.loggingService.GetCriteriaLog(job.SessionID, job.CriteriaID, job.MarketID)
	if err != nil {
		panic(err)
	}
	pageLog, err := a.loggingService.CreatePageLog(criteriaLog, job, "", job.Market.PageNumber)
	if err != nil {
		panic(err)
	}

	builder := NewURLBuilder()
	url := builder.GetPageURL(job)

	err = a.loggingService.PageLogSetVisitURL(pageLog, url)
	if err != nil {
		log.Println(err.Error())
	}

	foundAds := []jobs.Ad{}

	collector := colly.NewCollector()

	selector := "#main-content > div:nth-child(12) > div.row > div.col-sm-12.col-lg-12.col-md-12.col-xs-12.results-container > div.product__list--wrapper > div:nth-child(2) > div > div.product__listing.product__grid.row > div:nth-child(1) > div"
	selector = "div.product-list-item"
	counter := 0

	collector.OnHTML(selector, func(e *colly.HTMLElement) {
		ad := jobs.Ad{
			Title:              getTitle(e),
			Brand:              "",
			Model:              "",
			Year:               getYear(e),
			Km:                 getKm(e, url),
			Fuel:               "",
			Price:              getGrossPrice(e),
			AdID:               getAdID(e),
			Ad_url:             getAdHREF(e),
			SellerType:         "dealer",
			SellerName:         getSeller(e),
			SellerNameInMarket: getSeller(e),
			SellerOwnURL:       getSeller(e),
			SellerMarketURL:    getSeller(e),
			Thumbnail:          getThumbnail(e),
		}
		foundAds = append(foundAds, ad)
		counter++
	})

	err = collector.Visit(url)
	if err != nil {
		panic(err)
	}
	collector.Wait()

	if len(foundAds) == 0 {
		log.Println("NO MORE RESULTS -> SO RETURN !!!!!")
		err2 := a.loggingService.PageLogSetPageScraped(pageLog, len(foundAds), true)

		if err2 != nil {
			log.Println(err2.Error())
		}
		return icollector.AdsResults{
			Ads:        nil,
			IsLastPage: true,
			Error:      nil,
		}
	}

	err2 := a.loggingService.PageLogSetPageScraped(pageLog, len(foundAds), false)

	if err2 != nil {
		log.Println(err2.Error())
	}
	return icollector.AdsResults{
		Ads:        &foundAds,
		IsLastPage: false,
		Error:      nil,
	}
}

func getThumbnail(e *colly.HTMLElement) *string {
	thumbURL := ""
	selector := "div.image-content > a > picture > img"
	selection := e.DOM.Find(selector)
	src, exists := selection.Attr("src")
	if exists {
		thumbURL = src
	}
	return &thumbURL
}

// https://gebrauchtwagen.bmw.de/nsc/search?q=%3Arelevance%3Acondition-firstRegistrationYear%3A2019%3Acondition-mileageRange%3A%253C%2520150%2C000%2520km%3A%3Aprice-asc%3Aseries%3AX%3A%3Aprice-asc%3Amodel%3AX6%3A%3Aprice-asc%3Aenvironment-fuelType%3ADiesel&searchId=
// https://gebrauchtwagen.bmw.de/nsc/search?q=:relevance:condition-firstRegistrationYear:2019:condition-mileageRange:%3C%20150,000%20km::price-asc:series:X::price-asc:model:X6::price-asc:environment-fuelType:Diesel&page=1

func getTitle(e *colly.HTMLElement) *string {
	text := ""
	selector := "div.details > div.product-list-product-header > a"
	selection := e.DOM.Find(selector)
	text = selection.Text()
	text = strings.TrimSpace(text)
	return &text
}

func getGrossPrice(e *colly.HTMLElement) int {
	priceStr := ""
	selector := "div.details > div.price-de > div.price-inner > span.bigger"
	selection := e.DOM.Find(selector)
	priceStr = selection.Text()
	priceStr = strings.TrimSpace(priceStr)
	priceStr = strings.TrimSuffix(priceStr, " €")
	priceStr = strings.Replace(priceStr, ".", "", -1)

	if strings.Contains(priceStr, ",") {
		priceStr = priceStr[:len(priceStr)-3]
	}

	price := 0
	price, err := strconv.Atoi(priceStr)
	if err != nil {
		panic(err)
	}
	return price
}

func getYear(e *colly.HTMLElement) int {
	// year := 0
	selector := "div.details > div.nav-tabs-content > div.tab-content > div.tab-pane.active > ul.product-feature-list > li:nth-child(1) > div.product-feature-list-value > span.product-feature-value > span"
	selection := e.DOM.Find(selector)
	yearMonthStr := selection.Text()
	if len(yearMonthStr) > 4 {
		yearMonthStr = yearMonthStr[3:]
	}
	year, err := strconv.Atoi(yearMonthStr)
	if err != nil {
		panic(err)
	}
	return year
}

func getKm(e *colly.HTMLElement, pageurl string) int {
	//km := 0
	selector := "div.details > div.nav-tabs-content > div.tab-content > div.tab-pane.active > ul.product-feature-list > li:nth-child(3) > div.product-feature-list-value > span.product-feature-value > span"
	selection := e.DOM.Find(selector)
	kmStr := selection.Text()
	kmStr = strings.TrimSpace(kmStr)
	kmStr = strings.TrimSuffix(kmStr, " km")

	if kmStr == "Diesel" {
		selector = "div.details > div.nav-tabs-content > div.tab-content > div.tab-pane.active > ul.product-feature-list > li:nth-child(2) > div.product-feature-list-value > span.product-feature-value > span"
		selection = e.DOM.Find(selector)
		kmStr = selection.Text()
		kmStr = strings.TrimSpace(kmStr)
		kmStr = strings.TrimSuffix(kmStr, " km")

	}
	kmStr = strings.Replace(kmStr, ".", "", -1)
	km, err := strconv.Atoi(kmStr)

	if err != nil {
		panic(err)
	}
	return km
}

func getSeller(e *colly.HTMLElement) *string {
	dealer := "BMW.de"
	selector := "div.dealer-infos > p.dealer-name"
	selection := e.DOM.Find(selector)
	dealer = strings.TrimSpace(selection.Text())
	return &dealer
}

func getAdID(e *colly.HTMLElement) string {
	adID := ""
	selector := "div.card-top-items > div.carpark-content > a"
	selection := e.DOM.Find(selector)
	id, exists := selection.Attr("data-productcode")
	if !exists {
		panic(errors.New("NO AD ID FOUND"))
	}
	adID = id
	return adID
}

func getAdHREF(e *colly.HTMLElement) string {
	href := "https://gebrauchtwagen.bmw.de/"
	selector := "div.image-content > a"
	selection := e.DOM.Find(selector)
	hrefStr, exists := selection.Attr("href")
	if exists {
		href += hrefStr
	}
	return href

}
