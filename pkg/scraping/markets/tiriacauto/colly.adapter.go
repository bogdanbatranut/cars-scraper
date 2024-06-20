package tiriacauto

import (
	"carscraper/pkg/jobs"
	"carscraper/pkg/scraping/icollector"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly"
)

type TiriacAutoCollyMarketAdapter struct {
}

func NewTiriacAutoCollyMarketAdapter() *TiriacAutoCollyMarketAdapter {
	return &TiriacAutoCollyMarketAdapter{}
}

func (a TiriacAutoCollyMarketAdapter) GetAds(job jobs.SessionJob) icollector.AdsResults {
	builder := NewTiriacAutoURLBuilder()
	url := builder.GetURL(job)
	foundAds := a.findAds(*url, job.Criteria)
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
		log.Println(response.StatusCode)
		if response.StatusCode == 429 {
			log.Println("TOO MANY REQUESTS !!!")
		}
		responseBody := string(response.Body)
		if strings.Contains(responseBody, "https://www.tiriacauto.ro/deschide-cont-nou") {
			err := os.WriteFile("create.html", response.Body, 0644)
			check(err)
		}
	})

	selector := "body  > div:nth-child(4) > div > div.row > div.scrollable.col-12.col-md-8.pl-md-0.col-sm-12.tar-search-container > div > div.row"
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
		return icollector.AdsResults{
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
