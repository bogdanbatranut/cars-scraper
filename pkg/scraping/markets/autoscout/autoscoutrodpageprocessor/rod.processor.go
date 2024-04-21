package autoscoutrodpageprocessor

import (
	"carscraper/pkg/jobs"
	"carscraper/pkg/scraping/icollector"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
	gorodutils "github.com/go-rod/rod/lib/utils"
	"github.com/ysmood/gson"
)

type AutoscoutRodProcessor struct {
}

func (p AutoscoutRodProcessor) ProcessPage(page *rod.Page) icollector.AdsResults {
	log.Println("a4")
	accCookiesElem, err := page.Sleeper(rod.NotFoundSleeper).Element("#as24-cmp-popup > div > div._acceptance-buttons_1fb0r_82 > button._consent-accept_1fb0r_111")
	cookieBtnFound := true
	if errors.Is(err, &rod.ErrElementNotFound{}) {
		//fmt.Println("cookie button not found")
		cookieBtnFound = false
	} else if err != nil {
		panic(err)
	}

	if cookieBtnFound {
		accCookiesElem.MustClick()
	}

	log.Println("getting articles")
	articles := page.MustElements(".list-page-item")

	var ads []jobs.Ad
	log.Println("Found ", len(articles), " articles on page")
	if len(articles) == 0 {
		log.Println(" 0 articles")
		//ts := time.Now().UnixMilli()
		tsStr := time.Now().Format("2006-01-02_15:04:05")
		//page.MustScreenshot(fmt.Sprintf("my_%s.png", strconv.FormatInt(ts, 10)))
		page.MustScreenshot(fmt.Sprintf("my_%s.png", tsStr))

		// customization version
		img, _ := page.Screenshot(true, &proto.PageCaptureScreenshot{
			Format:  proto.PageCaptureScreenshotFormatJpeg,
			Quality: gson.Int(90),
			Clip: &proto.PageViewport{
				X:      0,
				Y:      0,
				Width:  2300,
				Height: 3200,
				Scale:  1,
			},
			FromSurface: true,
		})
		//_ = gorodutils.OutputFile(fmt.Sprintf("my_%s.jpg", strconv.FormatInt(ts, 10)), img)
		_ = gorodutils.OutputFile(fmt.Sprintf("my_%s.jpg", tsStr), img)
		return icollector.AdsResults{
			Ads:        nil,
			IsLastPage: true,
			Error:      nil,
		}
	} else {
		log.Println("More articles--> ", len(articles))
		for _, article := range articles {

			page.Mouse.MustScroll(0, 700)
			log.Println("s")
			//article.MustWaitLoad()
			//log.Println("s1")
			wait := page.Timeout(time.Second * 2).MustWaitRequestIdle()
			log.Println("s2")
			wait()
			log.Println("s3")
			carModel := article.MustAttribute("data-model")
			log.Println("s3")
			adThumb := getThumbNailFromArticle(article)
			log.Println("s4")
			sellerName := getSellerNameFromArticle(article)
			log.Println("vvvvvvvvv")
			ad := jobs.Ad{
				Brand:              getBrandFromArticle(article),
				Model:              *carModel,
				Year:               getYearFromArticle(article),
				Km:                 getMileageFromArticle(article),
				Fuel:               getFuelFromArticle(article),
				Price:              getGrossPriceFromArticle(article),
				AdID:               getAdIDFromArticle(article),
				Ad_url:             getHREFFromArticle(article),
				SellerType:         getSellerTypeFromArticle(article),
				SellerName:         &sellerName,
				SellerNameInMarket: &sellerName,
				SellerOwnURL:       &sellerName,
				SellerMarketURL:    &sellerName,
				Thumbnail:          adThumb,
			}
			ads = append(ads, ad)
		}
		isLastPage := isLastPage(page)
		log.Println("returning page results")
		return icollector.AdsResults{
			Ads:        &ads,
			IsLastPage: isLastPage,
			Error:      nil,
		}
	}

}
func getTotalResultsFromPage(page *rod.Page) float64 {
	elem := page.MustElement("#__next > div > div > div.ListPage_wrapper__vFmTi > div.ListPage_container__Optya > main > header > div.ListHeader_top__N6YWA > h1 > span > span:nth-child(1)")
	elemText, err := elem.Text()
	if err != nil {
		panic(err)
	}
	totalResultsStr := strings.Replace(elemText, ".", "", -1)
	totalResults_, err := strconv.Atoi(totalResultsStr)
	if err != nil {
		return -1
	}
	totalResults := float64(totalResults_)
	return totalResults
}
func getSellerTypeFromArticle(article *rod.Element) string {
	sellerType := "dealer"

	sellerFAttr, err := article.Attribute("data-seller-type")
	if err != nil {
		panic(err)
	}
	if *sellerFAttr != "d" {
		sellerType = "privat"
	}
	return sellerType
}
func getAdIDFromArticle(article *rod.Element) string {
	adId, err := article.Attribute("id")
	if err != nil {
		panic(err)
	}
	return *adId
}
func getHREFFromArticle(article *rod.Element) string {
	adHref, err := article.MustElement("div > div.ListItem_header__J6xlG.ListItem_header_new_design__Rvyv_ > a").Attribute("href")
	if err != nil {
		panic(err)
	}
	return *adHref
}
func getBrandFromArticle(article *rod.Element) string {
	log.Println("getBrand")
	brandStr, err := article.Attribute("data-make")
	if err != nil {
		panic(err)
	}
	log.Println("returning brand")
	return *brandStr
}
func getYearFromArticle(article *rod.Element) int {
	log.Println("get year")
	yearStr, err := article.Attribute("data-first-registration")
	if err != nil {
		panic(err)
	}
	if *yearStr == "new" {
		return time.Now().Year()
	}
	year, err := strconv.Atoi((*yearStr)[3:])
	if err != nil {
		panic(err)
	}
	log.Println("got year")
	return year
}
func getMileageFromArticle(article *rod.Element) int {
	log.Println("get mileage")
	kmStr, err := article.Attribute("data-mileage")
	if err != nil {
		panic(err)
	}
	km, err := strconv.Atoi(*kmStr)
	if err != nil {
		log.Println(err)
		log.Println("got mileage")
		return -1
	}
	log.Println("got mileage")
	return km
}
func getGrossPriceFromArticle(article *rod.Element) int {
	log.Println("get price")
	priceStr, err := article.Attribute("data-price")
	if err != nil {
		panic(err)
	}
	price, err := strconv.Atoi(*priceStr)
	if err != nil {
		panic(err)
	}
	log.Println("got price")
	return price
}
func getThumbNailFromArticle(article *rod.Element) *string {
	log.Println("get thumb")
	if article == nil {
		log.Println("Error getting thumbnail")
		return nil
	}
	imgElement, err := article.Element("div.ListItem_wrapper__TxHWu > div.Gallery_wrapper__iqp3u > section > div:nth-child(1) > picture > img")

	if err != nil {
		log.Println("-> ", article.MustText())
		log.Println(err.Error())
		if err.Error() == "cannot find element" {
			return nil
		}
	}
	src, err := imgElement.Attribute("src")
	if err != nil {
		panic(err)
	}
	log.Println("got thumb")
	return src
}
func getSellerNameFromArticle(article *rod.Element) string {
	log.Println("get seller name")
	if getSellerTypeFromArticle(article) == "privat" {
		sellerElement := article.MustElement("span.SellerInfo_private__THzvQ")
		log.Println("got seller name")
		return sellerElement.MustText()
	}
	if getSellerTypeFromArticle(article) == "dealer" {
		sellerElement := article.MustElement("div.SellerInfo_wrapper__XttVo > span.SellerInfo_name__nR9JH")
		log.Println("got seller name")
		return sellerElement.MustText()
	}
	log.Println("got seller name")
	return "No name retreived..."
}
func getFuelFromArticle(article *rod.Element) string {
	log.Println("get fuel")
	adFuelStr, err := article.Attribute("data-fuel-type")
	if err != nil {
		panic(err)
	}
	if *adFuelStr == "d" {
		return "diesel"
	}
	if *adFuelStr == "b" {
		return "petrol"
	}
	log.Println("got fuel")
	return "not found"
}
func isLastPage(page *rod.Page) bool {
	log.Println("get islastpage")
	elem := page.MustElement("#__next > div > div > div.ListPage_wrapper__vFmTi > div.ListPage_container__Optya > main > div.ListPage_pagination__4Vw9q > nav > ul > li:last-child > button")
	disabledAttr := elem.MustAttribute("disabled")
	if disabledAttr != nil {
		log.Println("got isLastPage")
		return true
	}
	log.Println("got isLastPage")
	return false
}
func getSellerURLInMarket(article *rod.Element) string {
	return ""
}
