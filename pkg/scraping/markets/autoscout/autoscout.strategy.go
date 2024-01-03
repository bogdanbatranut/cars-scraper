package autoscout

import (
	"carscraper/pkg/jobs"
	"fmt"
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
		}
		foundAds = append(foundAds, carad)
	})

	if executionErr != nil {
		return nil, false, executionErr
	}

	//t := time.Now()
	//t = t.AddDate(1, 0, 0)
	//
	//log.Println(t)
	//
	//cookie := http.Cookie{
	//	Name:       "cconsent-v2",
	//	Value:      "%7B%22purpose%22%3A%7B%22legitimateInterests%22%3A%5B25%5D%2C%22consents%22%3A%5B%5D%7D%2C%22vendor%22%3A%7B%22legitimateInterests%22%3A%5B10211%2C10218%2C10441%5D%2C%22consents%22%3A%5B%5D%7D%7D",
	//	Path:       "/",
	//	Domain:     ".autoscout24.ro",
	//	Expires:    t,
	//	RawExpires: "",
	//	MaxAge:     0,
	//	Secure:     false,
	//	HttpOnly:   false,
	//	SameSite:   http.SameSiteLaxMode,
	//	Raw:        "",
	//	Unparsed:   nil,
	//}
	//var cookies []*http.Cookie
	//cookies = append(cookies, &cookie)
	//
	//cookie = http.Cookie{
	//	Name:       "as24-cmp-signature",
	//	Value:      "ZCen3kxXtVWBetBY98GF2RUzuWQZcIoIlFWyiRN1pJ0lc3O6P3DZuWZoaiVr9KBCYW56Kc42kptzFKP%2F4D%2BsK7uKTrAvczACES8GLVQhtiNQqu%2FjNbDB%2F%2FJzL17ss1qlU3jWKNFbL0U%2F%2FgBpY5HGheKr8szZx%2B7vkqaFE1MHoOQ%3D",
	//	Path:       "/",
	//	Domain:     ".autoscout24.ro",
	//	Expires:    t,
	//	RawExpires: "",
	//	MaxAge:     0,
	//	Secure:     false,
	//	HttpOnly:   false,
	//	SameSite:   http.SameSiteLaxMode,
	//	Raw:        "",
	//	Unparsed:   nil,
	//}
	//cookies = append(cookies, &cookie)
	//
	//c.SetCookies(url, cookies)
	//
	//cookie = http.Cookie{
	//	Name:       "addtl_consent",
	//	Value:      "1~",
	//	Path:       "/",
	//	Domain:     ".autoscout24.ro",
	//	Expires:    t,
	//	RawExpires: "",
	//	MaxAge:     0,
	//	Secure:     false,
	//	HttpOnly:   false,
	//	SameSite:   http.SameSiteLaxMode,
	//	Raw:        "",
	//	Unparsed:   nil,
	//}
	//cookies = append(cookies, &cookie)
	//
	//cookie = http.Cookie{
	//	Name:       "euconsent-v2",
	//	Value:      "CP3f0kAP3f0kAGNACBROAgEsAP_gAEPgAAAAJStR5D7dbWlBcXp3aPswWY1T19DxpsQhBhaAg6AFiDOQcIwGk2AyNAygJgACEBAEghJBIQFFHAEAAQCAQAgBBAHsIgEEgAAIIABEgEMAQQNIAAgKCIAAAQAYgEAlEFAAmBiQANLkTcigAIAADgAYAAABAIABAgIBAAAYQBIAAAAAACAAAAoAAAAAAAAAAAAAAAAAQAAAIIQoBgChUQAlAQUhFoOEQCAEQVhARAIAAAASAggAACBAQYAwCEWAiAACAAAAAAAAAggABAAAJAAAAAAAAQAAAAAAIAAAAAAAIAEBAAGAAQAAAAAgKAIAAAAAAAAAEAEAAgAhQABACSUCAAAAADgAAAAABAIAAAAAAAgAIAAAAAAAAAAAAQAIB46BiAAsACoAHAAQAAvgBkAGgAPAAiABMACrAFwAXQAxABmADeAHoAP0AhgCJAEsAJoAUYAwABhgDRAHtAPwA_QCLAEdAJKASkAuYBeQDFAHUAReAkQBKgCZAFDgKPAU2AtgBcgDBgGSAMnAZZA1cDWAHFgPHJQGAAFgAcAB4AEQAJgAVQAuABigEMARIAjgBRgDAAH4AXMAxQB1AEXgJEAUeAtgBkgDJwGsAQhKQJwAFgAVAA4ACCAGQAaAA8ACIAEwAKQAVQAxABmAD9AIYAiQBRgDAAGjAPwA_QCLAEdAJKASkAuYBeQDFAHUAReAkQBQ4CmwFsALkAZIAycBlkDWANZAcEA8cCEIQAWABsAEgARwBpADnAIOATsAzQC_wGLAMhCQLwAFgAVAA4AB4AEAAL4AZABoADwAIgATAAqgBmADeAHoAPwAhIBDAESAI4ASwAmgBgADDAGWAO4Ae0A_AD9AI0ASUAlIBcwDFAGiASIAocBR4CkQFNgLYAXIAwYBkgDJwGZwNXA1kBwQDxwIQhgBIAiwBRgDnAOoAocBTYDFgGsgPHEACAASACLAGkAOcAiIChwGsgPHHADgASABHACgAOcAd0BBwEIAIiATsBf4DBAGLAMhAZUAzMiACAACAEIoANAAVACIAJAAWgBHAC2AI4Ac4A7gCDgE7AP-AwQBixCAWAHoARwAwAB3AFzAMUAdQBKgC5AGTgPHIAAgBzgMEJABgAjgDuAIOAv8BiwDxy0AUARwAwAB3AMzAeOWABADLAI4.YAAAAAAAA4CA",
	//	Path:       "/",
	//	Domain:     ".autoscout24.ro",
	//	Expires:    t,
	//	RawExpires: "",
	//	MaxAge:     0,
	//	Secure:     false,
	//	HttpOnly:   false,
	//	SameSite:   http.SameSiteLaxMode,
	//	Raw:        "",
	//	Unparsed:   nil,
	//}
	//cookies = append(cookies, &cookie)

	//c.SetCookies(url, cookies)

	c.OnRequest(func(request *colly.Request) {
		request.Headers.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36")
	})

	err := c.Visit(url)
	log.Println("AUTOSCOUT Visiting ", url)
	if err != nil {
		return nil, false, err
	}
	c.Wait()
	if len(foundAds) == 0 {
		log.Println("WE NO RESULTS SO RETURN !!!!!")
		return nil, true, nil
	}
	log.Println("AUTOSCOUT found ads : ", len(foundAds))
	return foundAds, isLastPage, nil
}
