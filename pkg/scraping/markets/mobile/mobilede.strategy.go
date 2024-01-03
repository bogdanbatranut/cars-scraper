package mobile

import (
	"carscraper/pkg/jobs"
	"fmt"
	"log"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly"
)

type MobileDeStrategy struct {
}

func NewMobileDeStrategy() MobileDeStrategy {
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

	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36"),
	)
	isLastPage := false

	var executionErr error

	var totalResults float64
	c.OnHTML("body > div.g-content > div > div.u-display-flex.u-margin-top-18 > section > section.result-block-header.g-row > div > h1", func(e *colly.HTMLElement) {
		spaceIndex := strings.Index(e.Text, " ")
		log.Printf("Total results : %s", e.Text)
		totalResultsStr := e.Text[:spaceIndex]
		totalResultsStr = strings.Replace(totalResultsStr, ".", "", -1)
		totalResults_, err := strconv.Atoi(totalResultsStr)
		if err != nil {
			executionErr = err
			return
		}
		totalResults = float64(totalResults_)

		numberOfTotalPages := math.Ceil(totalResults / 50)

		//log.Printf("Number fo total pages %.4f current page %d from total results %d", numberOfTotalPages, pageNumber, totalResults_)

		if float64(pageNumber) == numberOfTotalPages || totalResults == 0 {
			isLastPage = true
		}

	})

	if executionErr != nil {
		return nil, false, executionErr
	}

	// On every a element which has href attribute call callback
	c.OnHTML("article.list-entry.g-row", func(e *colly.HTMLElement) {
		sellerType := e.DOM.Find("div > div.g-row.js-ad-entry > a > div.g-col-s-12.g-col-m-8 > div:nth-child(2) > div.u-text-grey-60.g-col-s-8.g-col-m-9.u-margin-bottom-9").Text()
		if strings.ContainsAny(sellerType, "dealer") {
			sellerType = "dealer"
		} else {
			sellerType = "privat"
		}

		mobileAdId, exists := e.DOM.Find("div > div.g-row.js-ad-entry > a").Attr("data-vehicle-id")
		if !exists {
			mobileAdId = "NOT FOUND!!"
		}

		mobileAdHref, exists := e.DOM.Find("div > div.g-row.js-ad-entry > a").Attr("href")
		if !exists {
			mobileAdHref = "NOT FOUND!!"
		}

		//title := e.DOM.Find("div > div.g-row.js-ad-entry > a > div.g-col-s-12.g-col-m-8 > div.vehicle-text.g-row > h3").Text()
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
			executionErr = err
			return
		}
		kmStr := "0"
		if len(yearAndKm) >= 7 {
			if strings.Contains(yearAndKm, "/") {
				km_ := yearAndKm[9:]
				kmStr := strings.Replace(km_, ".", "", -1)
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
			executionErr = err
			return
		}

		// div > div.g-row.js-ad-entry > a > div.g-col-s-12.g-col-m-8 > div.vehicle-text.g-row > div.vehicle-information.g-col-s-6.g-col-m-8 > p.u-text-bold

		grossPriceStr := e.DOM.Find("div > div.g-row.js-ad-entry > a > div.g-col-s-12.g-col-m-8 > div.vehicle-text.g-row > div.g-col-s-6.g-col-m-4.u-text-right > div > p.seller-currency.u-text-bold").Text()
		grossPriceStr = strings.Replace(grossPriceStr, "\u00a0EUR (brut)", "", -1)
		grossPriceStr = strings.Replace(grossPriceStr, ".", "", -1)

		grossPrice, err := strconv.Atoi(grossPriceStr)
		if err != nil {
			executionErr = err
			return
		}

		//fmt.Printf("Elem found: %v\n", e)
		//e.ForEach("div > div.g-row.js-ad-entry > a > div.g-col-s-12.g-col-m-8 > div.vehicle-text.g-row > div.g-col-s-6.g-col-m-4.u-text-right > div > p.seller-currency.u-text-bold", func(i int, element *colly.HTMLElement) {
		//	log.Println(element.Text)
		//})
		seller := "mobile.de"
		ad := jobs.Ad{
			Brand:              criteria.Brand,
			Model:              criteria.CarModel,
			Year:               year,
			Km:                 km,
			Fuel:               criteria.Fuel,
			Price:              grossPrice,
			AdID:               mobileAdId,
			Ad_url:             fmt.Sprintf("https://www.mobile.de%s", mobileAdHref),
			SellerType:         sellerType,
			SellerName:         &seller,
			SellerNameInMarket: &seller,
			SellerOwnURL:       &seller,
			SellerMarketURL:    &seller,
		}
		foundAds = append(foundAds, ad)
	})

	c.OnRequest(func(req *colly.Request) {
		fmt.Println("MOBILE Visiting", req.URL.String())
		fmt.Println("applying headers")
		req.Headers.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
		req.Headers.Add("Accept-Encoding", "gzip, deflate, br")
		req.Headers.Add("Accept-Language", "en-GB,en;q=0.9")
		req.Headers.Add("Sec-Ch-Ua", "\"Google Chrome\";v=\"119\", \"Chromium\";v=\"119\", \"Not?A_Brand\";v=\"24\"")
		req.Headers.Add("Sec-Ch-Ua-Mobile", "?0")
		req.Headers.Add("Sec-Ch-Ua-Platform", "\"macOS\"")
		req.Headers.Add("Sec-Fetch-Dest", "document")
		req.Headers.Add("Sec-Fetch-Mode", "navigate")
		req.Headers.Add("Sec-Fetch-Site", "none")
		req.Headers.Add("Sec-Fetch-User", "?1")
		req.Headers.Add("Upgrade-Insecure-Requests", "1")
		req.Headers.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36")

	})
	if executionErr != nil {
		return nil, false, executionErr
	}

	err := c.Visit(url)
	if err != nil {
		return nil, false, err
	}
	c.Wait()
	if len(foundAds) == 0 {
		log.Println("NO MORE RESULTS -> SO RETURN !!!!!")
		return nil, true, nil
	}
	log.Println("MOBILE found ads : ", len(foundAds))
	return foundAds, isLastPage, nil
}
