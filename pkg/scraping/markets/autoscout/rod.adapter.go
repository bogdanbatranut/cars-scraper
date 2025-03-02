package autoscout

import (
	"carscraper/pkg/jobs"
	"carscraper/pkg/scraping/icollector"
	"errors"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
)

type AutoscoutRodAdapter struct {
}

func NewAutoscoutRodAdapter() *AutoscoutRodAdapter {
	return &AutoscoutRodAdapter{}
}

func (a AutoscoutRodAdapter) GetAds(page *rod.Page) *icollector.AdsResults {
	//accCookiesElem, err := page.Sleeper(rod.NotFoundSleeper).Element("#as24-cmp-popup > div > div._acceptance-buttons_1fb0r_82 > button._consent-accept_1lphq_114")
	accCookiesElem, err := page.Sleeper(rod.NotFoundSleeper).Element("#as24-cmp-popup > div > div._acceptance-buttons_1lphq_85 > button._consent-accept_1lphq_114")
	cookieBtnFound := true
	if errors.Is(err, &rod.ErrElementNotFound{}) {
		//fmt.Println("cookie button not found")
		cookieBtnFound = false
	} else if err != nil {
		panic(err)
	}

	if cookieBtnFound {
		//accCookiesElem.MustClick()
		err := accCookiesElem.Click(proto.InputMouseButtonLeft, 1)
		if err != nil {
			return &icollector.AdsResults{
				Ads:        nil,
				IsLastPage: true,
				Error:      err,
			}
		}
	}

	//articles := page.MustElements(".list-page-item")
	articles, err := page.MustWaitDOMStable().Elements(".list-page-item")
	if err != nil {
		println(err)
	}

	var ads []jobs.Ad
	if len(articles) == 0 {
		log.Println(" No articles found")
		return &icollector.AdsResults{
			Ads:        nil,
			IsLastPage: true,
			Error:      nil,
		}
	} else {
		isLastPage := isLastPage(page)
		for _, article := range articles {
			page.Mouse.MustScroll(0, 700)
			wait := page.Timeout(time.Second / 2).MustWaitRequestIdle()
			wait()
			carModel := article.MustAttribute("data-model")
			adThumb := getThumbNailFromArticle(article)
			sellerName := getSellerNameFromArticle(article)
			articleTitle, err := getAdTitle(article)
			if err != nil {
				return &icollector.AdsResults{
					Ads:        nil,
					IsLastPage: true,
					Error:      err,
				}
			}
			ad := jobs.Ad{
				Title:              articleTitle,
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

		return &icollector.AdsResults{
			Ads:        &ads,
			IsLastPage: isLastPage,
			Error:      nil,
		}
	}
}

func getTotalResults(page *rod.Page) (*int, error) {
	elem, err := page.Element("#__next > div > div > div.ListPage_wrapper__vFmTi > header > div > div.ListHeaderExperiment_title_with_sort__Gj9w7 > h1 > span > span:nth-child(1)")
	if err != nil {
		return nil, err
	}
	elemStr, err := elem.Text()
	if err != nil {
		return nil, err
	}

	elemVal, err := strconv.Atoi(elemStr)
	if err != nil {
		return nil, err
	}

	return &elemVal, nil
}

func getAdTitle(article *rod.Element) (*string, error) {
	elem, err := article.Element("div.ListItem_wrapper__TxHWu > div.ListItem_header__J6xlG.ListItem_header_new_design__Rvyv_ > a > h2 > span")
	if err != nil {
		return nil, err
	}
	articleTitle, err := elem.Text()
	if err != nil {
		return nil, err
	}
	return &articleTitle, nil
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
	return "https://www.autoscout24.ro" + *adHref
}
func getBrandFromArticle(article *rod.Element) string {
	brandStr, err := article.Attribute("data-make")
	if err != nil {
		panic(err)
	}
	return *brandStr
}
func getYearFromArticle(article *rod.Element) int {
	yearStr, err := article.Attribute("data-first-registration")
	if err != nil {
		panic(err)
	}
	if *yearStr == "new" || *yearStr == "unknown" {
		return time.Now().Year()
	}
	year, err := strconv.Atoi((*yearStr)[3:])
	if err != nil {
		initialErr := err
		articleId, err := article.Attribute("id")
		if err != nil {
			log.Println("cannot get article ID")
		}
		log.Println("cannot get year from article ")
		log.Println(*articleId)
		panic(initialErr)
	}
	return year
}
func getMileageFromArticle(article *rod.Element) int {
	kmStr, err := article.Attribute("data-mileage")
	if err != nil {
		panic(err)
	}
	km, err := strconv.Atoi(*kmStr)
	if err != nil {
		log.Println(err)
		return -1
	}
	return km
}
func getGrossPriceFromArticle(article *rod.Element) int {
	priceStr, err := article.Attribute("data-price")
	if err != nil {
		panic(err)
	}
	price, err := strconv.Atoi(*priceStr)
	if err != nil {
		panic(err)
	}
	return price
}
func getThumbNailFromArticle(article *rod.Element) *string {
	if article == nil {
		log.Println("Error getting thumbnail")
		return nil
	}
	imgElement, err := article.Element("div.ListItem_wrapper__TxHWu > div.Gallery_wrapper__iqp3u > section > div:nth-child(1) > picture > img")

	if err != nil {
		log.Println(err.Error())
		if err.Error() == "cannot find element" {
			log.Println("Could not find image ")
			return nil
		}
	}
	src, err := imgElement.Attribute("src")
	if err != nil {
		panic(err)
	}
	return src
}
func getSellerNameFromArticle(article *rod.Element) string {
	if getSellerTypeFromArticle(article) == "privat" {
		sellerElement := article.MustElement("span.SellerInfo_private__THzvQ")
		//log.Println("got seller name")
		return sellerElement.MustText()
	}
	if getSellerTypeFromArticle(article) == "dealer" {
		sellerElement := article.MustElement("div.SellerInfo_wrapper__XttVo > span.SellerInfo_name__nR9JH")
		return sellerElement.MustText()
	}
	return "No name retreived..."
}
func getFuelFromArticle(article *rod.Element) string {
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
	return "not found"
}
func isLastPage(page *rod.Page) bool {
	elem, err := page.Element("#__next > div > div > div.ListPage_wrapper__vFmTi > div.ListPage_container__Optya > main > div.ListPage_pagination__4Vw9q > nav > ul > li:last-child > button")
	if err != nil {
		panic(err)
	}
	disabledAttr := elem.MustAttribute("disabled")
	if disabledAttr != nil {
		return true
	}
	return false
}
func getSellerURLInMarket(article *rod.Element) string {
	return ""
}
