package mobile_props

import (
	"carscraper/pkg/amconfig"
	"carscraper/pkg/jobs"
	"carscraper/pkg/logging"
	"carscraper/pkg/repos"
	"carscraper/pkg/scraping/icollector"
	"carscraper/pkg/scraping/markets/mobile/mobiledecollycollector"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly"
)

type MobilePropDECollyMarketAdapter struct {
	loggingService *logging.ScrapeLoggingService
	propertiesRepo *repos.PropertiesRepository
	builder        *URLBuilder
}

func NewMobilePropDECollyMarketAdapter(logingService *logging.ScrapeLoggingService, cfg amconfig.IConfig) *MobilePropDECollyMarketAdapter {
	return &MobilePropDECollyMarketAdapter{
		loggingService: logingService,
		propertiesRepo: repos.NewPropertiesRepository(cfg),
		builder:        NewURLBuilder(cfg),
	}
}

func (a MobilePropDECollyMarketAdapter) GetAds(job jobs.SessionJob) icollector.AdsResults {
	criteriaLog, err := a.loggingService.GetCriteriaLog(job.SessionID, job.CriteriaID, job.MarketID)
	if err != nil {
		panic(err)
	}
	pageLog, err := a.loggingService.CreatePageLog(criteriaLog, job, "", job.Market.PageNumber)
	if err != nil {
		panic(err)
	}

	url := a.builder.GetPageURL(job.Criteria, job.Market.PageNumber)

	err = a.loggingService.PageLogSetVisitURL(pageLog, url)
	if err != nil {
		log.Println(err.Error())
	}

	mobileCollector := mobiledecollycollector.NewMobileDECollyCollector().GetCollyCollector(job)

	foundAds := []jobs.Ad{}

	var executionErr error

	// On every a element which has href attribute call callback
	mobileCollector.OnHTML("article.list-entry.g-row", func(e *colly.HTMLElement) {

		title := getTitle(e)
		sellerType := getSellerType(e)
		mobileAdId := getAdId(e)
		mobileAdHref := getAdHref(e)
		thumbNail := getThumbnail(e)

		year, err := getYear(e)
		if err != nil {
			executionErr = err
			return
		}

		km, err := getKm(e)
		if err != nil {
			executionErr = err
			return
		}

		grossPrice, err := getGrossPrice(e)
		if err != nil {
			executionErr = err
			return
		}

		seller := "mobile.de"
		ad := jobs.Ad{
			//Brand:              criteria.Brand,
			//Model:              criteria.CarModel,
			Title: title,
			Year:  year,
			Km:    km,
			//Fuel:               criteria.Fuel,
			Price:              grossPrice,
			AdID:               mobileAdId,
			Ad_url:             mobileAdHref,
			SellerType:         sellerType,
			SellerName:         &seller,
			SellerNameInMarket: &seller,
			SellerOwnURL:       &seller,
			SellerMarketURL:    &seller,
			Thumbnail:          &thumbNail,
		}
		foundAds = append(foundAds, ad)
	})

	mobileCollector.Visit(url)

	if executionErr != nil {
		return icollector.AdsResults{
			Ads:        nil,
			IsLastPage: true,
			Error:      executionErr,
		}
	}
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

	if len(foundAds) < 50 {
		err2 := a.loggingService.PageLogSetPageScraped(pageLog, len(foundAds), true)

		if err2 != nil {
			log.Println(err2.Error())
		}
		return icollector.AdsResults{
			Ads:        &foundAds,
			IsLastPage: true,
			Error:      nil,
		}
	}

	log.Println("MOBILE found ads : ", len(foundAds))

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

func getSellerType(e *colly.HTMLElement) string {
	sellerType := e.DOM.Find("div > div.g-row.js-ad-entry > a > div.g-col-s-12.g-col-m-8 > div:nth-child(2) > div.u-text-grey-60.g-col-s-8.g-col-m-9.u-margin-bottom-9").Text()
	if strings.ContainsAny(sellerType, "dealer") {
		sellerType = "dealer"
	} else {
		sellerType = "privat"
	}
	return sellerType
}

func getAdId(e *colly.HTMLElement) string {
	mobileAdId, exists := e.DOM.Find("div > div.g-row.js-ad-entry > a").Attr("data-vehicle-id")
	if !exists {
		mobileAdId = "NOT FOUND!!"
	}
	return mobileAdId
}

func getAdHref(e *colly.HTMLElement) string {
	mobileAdHref, exists := e.DOM.Find("div > div.g-row.js-ad-entry > a").Attr("href")
	if !exists {
		mobileAdHref = "NOT FOUND!!"
	}
	return fmt.Sprintf("https://www.mobile.de%s", mobileAdHref)
}

func getThumbnail(e *colly.HTMLElement) string {
	thumbNail := e.DOM.Find("div > div.g-row.js-ad-entry >  a > div.thumbnail > noscript").Text()
	if thumbNail != "" {
		thumbSrc := strings.Split(thumbNail, " ")[2]
		thumbNail = strings.Split(thumbSrc, "=")[1] + "=" + strings.Split(thumbSrc, "=")[2]
		rep := "\""
		thumbNail = strings.Replace(thumbNail, rep, "", -1)
	}
	return thumbNail
}

func getYear(e *colly.HTMLElement) (int, error) {
	yearAndKm := e.DOM.Find("div > div.g-row.js-ad-entry > a > div.g-col-s-12.g-col-m-8 > div.vehicle-text.g-row > div.vehicle-information.g-col-s-6.g-col-m-8 > p.u-text-bold").Text()
	year, _, _ := time.Now().Date()
	yearStr := strconv.Itoa(year)

	if len(yearAndKm) >= 7 {
		if strings.ContainsAny(yearAndKm, "/") {
			yearStr = yearAndKm[3:7]
		} else {
			if strings.ContainsAny(yearAndKm, ",") {
				yearStr = yearAndKm[:5]
				yearStr = strings.Replace(yearStr, ",", "", -1)
				yearStr = strings.Replace(yearStr, " ", "", -1)
				yearStr = strings.Replace(yearStr, "\u00a0", "", -1)
			}
		}
	}
	year, err := strconv.Atoi(yearStr)
	if err != nil {
		return 0, err
	}
	return year, nil
}

func getKm(e *colly.HTMLElement) (int, error) {
	yearAndKm := e.DOM.Find("div > div.g-row.js-ad-entry > a > div.g-col-s-12.g-col-m-8 > div.vehicle-text.g-row > div.vehicle-information.g-col-s-6.g-col-m-8 > p.u-text-bold").Text()
	kmStr := "0"
	if len(yearAndKm) >= 7 {
		if strings.Contains(yearAndKm, "/") {
			km_ := yearAndKm[9:]
			kmStr = strings.Replace(km_, ".", "", -1)
			kmStr = strings.Replace(kmStr, "\u00a0", "", -1)
			kmStr = strings.Replace(kmStr, "km", "", -1)
		} else {
			kmStr = yearAndKm[:3]
		}

	} else {
		kmStr := strings.Replace(yearAndKm, ".", "", -1)
		kmStr = strings.Replace(kmStr, "\u00a0", "", -1)
		kmStr = strings.Replace(kmStr, "km", "", -1)
	}

	km, err := strconv.Atoi(kmStr)
	if err != nil {
		return 0, err
	}
	return km, nil
}

func getGrossPrice(e *colly.HTMLElement) (int, error) {
	grossPriceStr := e.DOM.Find("div > div.g-row.js-ad-entry > a > div.g-col-s-12.g-col-m-8 > div.vehicle-text.g-row > div.g-col-s-6.g-col-m-4.u-text-right > div > p.seller-currency.u-text-bold").Text()
	grossPriceStr = strings.Replace(grossPriceStr, "\u00a0EUR (brut)", "", -1)
	grossPriceStr = strings.Replace(grossPriceStr, ".", "", -1)

	grossPrice, err := strconv.Atoi(grossPriceStr)
	if err != nil {
		return 0, err
	}
	return grossPrice, nil
}

func getTitle(e *colly.HTMLElement) *string {
	var title *string
	titleText := e.DOM.Find("div > div.g-row.js-ad-entry > a > div.g-col-s-12.g-col-m-8 > div.vehicle-text.g-row > h3").Text()
	if titleText != "" {
		title = &titleText
	}
	return title
	// body > div.g-content > div > div.u-display-flex.u-margin-top-18 > section > div.result-list-section.js-result-list-section.u-clearfix > article:nth-child(3) > div > div.g-row.js-ad-entry > a > div.g-col-s-12.g-col-m-8 > div.vehicle-text.g-row > h3
}
