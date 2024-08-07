package tiriacauto

import (
	"carscraper/pkg/jobs"
	"carscraper/pkg/logging"
	"carscraper/pkg/scraping/icollector"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly"
)

type TiriacAutoCollyMarketAdapter struct {
	loggingService *logging.ScrapeLoggingService
}

func NewTiriacAutoCollyMarketAdapter(loggingService *logging.ScrapeLoggingService) *TiriacAutoCollyMarketAdapter {
	return &TiriacAutoCollyMarketAdapter{
		loggingService: loggingService,
	}
}

func (a TiriacAutoCollyMarketAdapter) GetAds(job jobs.SessionJob) icollector.AdsResults {

	criteriaLog, err := a.loggingService.GetCriteriaLog(job.SessionID, job.CriteriaID, job.MarketID)
	if err != nil {
		panic(err)
	}
	pageLog, err := a.loggingService.CreatePageLog(criteriaLog, job, "", job.Market.PageNumber)
	if err != nil {
		panic(err)
	}

	builder := NewTiriacAutoURLBuilder()
	url := builder.GetURL(job)

	err = a.loggingService.PageLogSetVisitURL(pageLog, *url)
	if err != nil {
		log.Println(err.Error())
	}

	foundAds := a.findAds(*url, job.Criteria)

	if foundAds.Error != nil {
		a.loggingService.PageLogSetError(pageLog, foundAds.Error.Error())
	} else {
		err2 := a.loggingService.PageLogSetPageScraped(pageLog, len(*foundAds.Ads), foundAds.IsLastPage)

		if err2 != nil {
			log.Println(err2.Error())
		}
	}

	return foundAds
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func (a TiriacAutoCollyMarketAdapter) findAds(url string, criteria jobs.Criteria) icollector.AdsResults {

	adsCollector := colly.NewCollector(colly.Async(false))
	adsCollector.Limit(
		&colly.LimitRule{
			DomainGlob:  "*",
			Delay:       5 * time.Second,
			RandomDelay: 60 * time.Second,
			Parallelism: 1,
		})
	adsResults := icollector.AdsResults{
		Ads:        nil,
		IsLastPage: false,
		Error:      nil,
	}

	foundAds := []jobs.Ad{}
	adsCollector.OnResponse(func(response *colly.Response) {

		log.Println("Response from tiriac auto", response.StatusCode)
		if response.StatusCode == 429 {
			log.Println("TOO MANY REQUESTS !!!")
		}
		responseBody := string(response.Body)
		if strings.Contains(responseBody, "https://www.tiriacauto.ro/deschide-cont-nou") {
			err := os.WriteFile("create.html", response.Body, 0644)
			check(err)
		}
	})
	//body > div:nth-child(4) > div > div.row > div.scrollable.col-12.col-md-8.pl-md-0.col-sm-12.tar-search-container > div > div:nth-child(4) > div:nth-child(1) > div
	selector := "body > div:nth-child(4) > div > div.row > div.scrollable.col-12.col-md-8.pl-md-0.col-sm-12.tar-search-container > div > div.row"
	adsCollector.OnHTML(selector, func(e *colly.HTMLElement) {

		e.ForEach("div.carItem", func(index int, e *colly.HTMLElement) {

			adHref, has := e.DOM.Find("a").Attr("href")
			if has {
				log.Println(adHref)
			}
			if adHref == "https://www.tiriacauto.ro/deschide-cont-nou" {
				log.Println(" Create account page !!!!")
				return
			}

			adHREFSplit := strings.Split(adHref, "-")
			leng := len(adHREFSplit) - 1
			adIdStr := adHREFSplit[leng]

			thumbNail, has := e.DOM.Find("a > img").Attr("src")
			if has {
				log.Println(thumbNail)
			}

			var title string
			title = strings.TrimSpace(e.DOM.Find("div.details > div > div").Text())
			log.Println(title)
			priceStr := e.DOM.Find("div.details > div.price > p.mb-0 > b").Text()
			priceStr = strings.TrimSpace(priceStr)
			priceStr = strings.Split(priceStr, " ")[0]
			priceStr = strings.ReplaceAll(priceStr, ".", "")
			price := 0
			price, err := strconv.Atoi(priceStr)
			if err != nil {
				log.Println(err)
			}
			log.Println(" Price ", price)

			km := 0
			kmStr := e.DOM.Find("div.details > ul.specs > li:nth-child(2)").Text()
			if kmStr != "" {
				kmStr = strings.ReplaceAll(kmStr, "Km", "")
				kmStr = strings.ReplaceAll(kmStr, ".", "")
				km, err = strconv.Atoi(kmStr)
				if err != nil {
					panic(err)
				}
			}

			log.Println(" KM ", km)

			//log.Println(e.DOM.Find("div.details > ul.specs > li:nth-child(3)").Text())
			yearStr := e.DOM.Find("div.details > ul.specs > li:nth-child(3)").Text()
			year, _, _ := time.Now().Date()
			if yearStr != "" && !strings.Contains(adHref, "auto-noi") {
				year, err = strconv.Atoi(yearStr)
				if err != nil {
					panic(err)
				}
			}

			log.Println(year)
			dealerName := "Tiriac Auto"
			dealerLink := "https://www.tiriacauto.ro/"
			ad := jobs.Ad{
				Title:              &title,
				Brand:              criteria.Brand,
				Model:              criteria.CarModel,
				Year:               year,
				Km:                 km,
				Fuel:               criteria.Fuel,
				Price:              price,
				AdID:               adIdStr,
				Ad_url:             adHref,
				SellerType:         "dealer",
				SellerName:         &dealerName,
				SellerNameInMarket: &dealerName,
				SellerOwnURL:       &dealerLink,
				SellerMarketURL:    &dealerLink,
				Thumbnail:          &thumbNail,
			}
			foundAds = append(foundAds, ad)
		})
	})

	adsCollector.OnRequest(func(request *colly.Request) {
		log.Println(" VISITING ::: ", request.URL.String())
	})

	err := adsCollector.Visit(url)
	if err != nil {
		adsResults = icollector.AdsResults{
			Ads:        nil,
			IsLastPage: true,
			Error:      err,
		}
	}
	adsCollector.Wait()

	adsResults.Ads = &foundAds
	adsResults.IsLastPage = true

	return adsResults

}
