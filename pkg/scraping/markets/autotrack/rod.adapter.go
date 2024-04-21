package autotrack

import (
	"carscraper/pkg/jobs"
	"carscraper/pkg/scraping/icollector"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
)

type AutoTrackNLRodAdapter struct{}

func NewAutoTrackNLRodAdapter() *AutoTrackNLRodAdapter { return &AutoTrackNLRodAdapter{} }

func checkForCookiesAndAccept(page *rod.Page) {
	modalContainer, err := page.Element("#pg-shadow-host")
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("MODAL CONTAINER FOUND")
	modal, err := modalContainer.ShadowRoot()
	if err != nil {
		panic(err)
	}
	log.Println("MODAL SHADOW ROOT FOUND")
	acceptCookiesBtn, err := modal.Element("#pg-accept-btn")
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("ACCEPT COOKIE BUTTON FOUND")

	err = acceptCookiesBtn.Click(proto.InputMouseButtonLeft, 1)
	if err != nil {
		panic(err)
	}
	log.Println("ACCEPT COOKIE BUTTON CLICKED")
}

func getTotalResults(page *rod.Page) (*int, error) {
	element, err := page.Element("#__next > div > main > div.main.container.pt-3 > div > div.jsx-2926786665.col-lg-9.filters__results-wrapper > div.jsx-2926786665.filters__results > div.SaveSearchWrapper__ResultHeading-sc-117eb1i-0.hHndHb > span")
	if err != nil {
		element, err = page.Element("#__next > div > main > div.main.container.pt-3 > div > div.jsx-2926786665.col-lg-9.filters__results-wrapper > div.jsx-2926786665.filters__results > div.SaveSearchWrapper__ResultHeading-sc-117eb1i-0.hHndHb > span")
		if err != nil {
			return nil, err
		}

	}
	elemetTxt, err := element.Text()
	if err != nil {
		return nil, err
	}
	log.Println(elemetTxt)
	resultsStr := strings.Split(elemetTxt, " ")[0]
	result, err := strconv.Atoi(resultsStr)
	if err != nil {
		if strings.Contains(resultsStr, ".") {
			resultsStr = strings.Replace(resultsStr, ".", "", -1)
			result, err = strconv.Atoi(resultsStr)
			if err != nil {
				return nil, err
			}
		}
	}
	return &result, nil
}

func (a AutoTrackNLRodAdapter) GetAds(page *rod.Page) *icollector.AdsResults {
	checkForCookiesAndAccept(page)
	totalResults, err := getTotalResults(page)
	if err != nil {
		return &icollector.AdsResults{
			Ads:        nil,
			IsLastPage: true,
			Error:      nil,
		}
	}
	if totalResults == nil {
		return &icollector.AdsResults{
			Ads:        nil,
			IsLastPage: true,
			Error:      nil,
		}
	}

	articles, err := page.MustWaitLoad().MustWaitDOMStable().Elements("li[data-testid=result-list-item]")
	if err != nil {
		return &icollector.AdsResults{
			Ads:        nil,
			IsLastPage: true,
			Error:      nil,
		}
	}
	if len(articles) > 0 && len(articles) <= 30 {
		return &icollector.AdsResults{
			Ads:        processElements(articles),
			IsLastPage: true,
			Error:      nil,
		}
	}
	paginationElement, err := page.Element("#__next > div > main > div.main.container.pt-3 > div > div.jsx-2926786665.col-lg-9.filters__results-wrapper > div.jsx-2926786665.filters__results > section > div.jsx-3584031782.pagination__container > ul")
	if err != nil {
		return &icollector.AdsResults{
			Ads:        nil,
			IsLastPage: true,
			Error:      err,
		}
	}
	isLast, err := isLastPage(paginationElement)
	if err != nil {
		return &icollector.AdsResults{
			Ads:        nil,
			IsLastPage: true,
			Error:      err,
		}
	}
	return &icollector.AdsResults{
		Ads:        processElements(articles),
		IsLastPage: *isLast,
		Error:      nil,
	}
}

func processElements(elements rod.Elements) *[]jobs.Ad {
	var ads []jobs.Ad
	for _, element := range elements {
		soldElement, err := element.Element("span.ItemTag__ItemTagStatusLabel-sc-hnf0c8-4.ItemTag__VerkochtStatus-sc-hnf0c8-5.SKbfR.ipKrWP")
		if err != nil {
			log.Println("First sold element")
			log.Println(err.Error())
		}
		if soldElement != nil {
			continue
		}
		soldElement, err = element.Element("div.StyledItemContent__StyleItemContentContainer-sc-1fqlnst-1.hzyHqR > div > div > div.StyledItemTags-sc-46f8tm-0.hRDWcg > span")
		if err != nil {
			log.Println("Second sold element")
			log.Println(err.Error())
		}
		if soldElement != nil {
			continue
		}
		adTitle, err := element.Element("div.StyledItemContent__StyleItemContentContainer-sc-1fqlnst-1.SpzQV > div > div.StyledItemContent__StyleItemContentContentHeader-sc-1fqlnst-3.dVpSsc > h3")
		if err != nil {
			adTitle, err = element.Element("div.StyledItemContent__StyleItemContentContainer-sc-1fqlnst-1.hzyHqR > div > div > div.StyledItemContent__StyleItemContentContentHeader-sc-1fqlnst-3.dVpSsc > h3")

			if err != nil {
				adTitle, err = element.Element("div.StyledItemContent__StyleItemContentContainer-sc-1fqlnst-1.fbbMe > div > div.StyledItemContent__StyleItemContentContentContainer-sc-1fqlnst-2.guFJQU > div.StyledItemContent__StyleItemContentContentHeader-sc-1fqlnst-3.dVpSsc > h3")
				if err != nil {
					log.Println(element.MustText())
					panic(err)
				}
			}
		}
		adTitleText := adTitle.MustText()
		log.Println("Title: ", adTitleText)
		href, err := getHref(element)
		if err != nil {
			panic(err)
		}
		//log.Println("HREF: ", href)
		adID := getAdID(*href)
		//log.Println("Ad ID: ", adID)
		dealer := getDealer(element)
		//log.Println("Dealer: ", dealer)
		price, err := getPrice(element)
		if err != nil {
			// invalid price so continue
			continue
		}
		//log.Println("Price: ", price)

		year, err := getYear(element)
		if err != nil {
			panic(err)
		}
		//log.Println("Year : ", year)

		km, err := getKm(element)
		if err != nil {
			panic(err)
		}
		//log.Println("KM : ", km)

		thumbNail, err := getThumbnail(element)
		if err != nil {
			panic(err)
		}
		//log.Println("Tumbnail: ", thumbNail)

		ad := jobs.Ad{
			Brand:              "",
			Model:              "",
			Year:               *year,
			Km:                 *km,
			Fuel:               "",
			Price:              *price,
			AdID:               adID,
			Ad_url:             *href,
			SellerType:         "dealer",
			SellerName:         &dealer,
			SellerNameInMarket: &dealer,
			SellerOwnURL:       &dealer,
			SellerMarketURL:    &dealer,
			Thumbnail:          thumbNail,
		}

		ads = append(ads, ad)
	}
	return &ads
}

func getThumbnail(element *rod.Element) (*string, error) {
	thumbNailElement, err := element.Element("a > article > div.StyledItemContent__StyleItemContentContainer-sc-1fqlnst-1.SpzQV > figure > div > div > div > div.CarouselWrapperComponent__Slides-sc-1o127aw-6.kroIPK > div:nth-child(1) > img")
	if err != nil {
		thumbNailElement, err = element.Element("a > article > div.StyledItemContent__StyleItemContentContainer-sc-1fqlnst-1.hzyHqR > figure > div > div > div > div.CarouselWrapperComponent__Slides-sc-1o127aw-6.kroIPK > div > img")
		if err != nil {
			log.Println("Get THUMBNAIL ", err)
			return nil, err
		}
	}
	thumbNailValue, err := thumbNailElement.Attribute("src")
	if err != nil {
		log.Println("Get THUMBN attribute src ", err)
		return nil, err
	}
	return thumbNailValue, nil
}

func getKm(element *rod.Element) (*int, error) {
	kmElem, err := element.Element("a > article > div.StyledItemContent__StyleItemContentContainer-sc-1fqlnst-1.SpzQV > div > div.StyledItemContent__FlexDiv-sc-1fqlnst-0.StyledItemContent__StyleItemContentContentDetails-sc-1fqlnst-5.kPOPyZ.jCmSiq > span.StyledItemContent__StyleItemContentContentDetailsMileage-sc-1fqlnst-10.kqpAaS")
	if err != nil {
		kmElem, err = element.Element("a > article > div.StyledItemContent__StyleItemContentContainer-sc-1fqlnst-1.hzyHqR > div > div > div.StyledItemContent__FlexDiv-sc-1fqlnst-0.StyledItemContent__StyleItemContentContentDetails-sc-1fqlnst-5.kPOPyZ.jCmSiq > span.StyledItemContent__StyleItemContentContentDetailsMileage-sc-1fqlnst-10.kqpAaS")
		if err != nil {
			log.Println("Get KM invalid selector ", err)
			return nil, err
		}
	}
	kmStr, err := kmElem.Text()
	if err != nil {
		log.Println("Get KM invalid element text ", err)
		return nil, err
	}
	kmStr = kmStr[0 : len(kmStr)-3]
	kmStr = strings.Replace(kmStr, ".", "", -1)
	km, err := strconv.Atoi(kmStr)
	if err != nil {

		log.Println("Get KM, err")
		return nil, err
	}
	return &km, nil
}

func getYear(element *rod.Element) (*int, error) {
	yearElem, err := element.Element("a > article > div.StyledItemContent__StyleItemContentContainer-sc-1fqlnst-1.SpzQV > div > div.StyledItemContent__FlexDiv-sc-1fqlnst-0.StyledItemContent__StyleItemContentContentDetails-sc-1fqlnst-5.kPOPyZ.jCmSiq > span.StyledItemContent__StyleItemContentContentDetailsDate-sc-1fqlnst-8.bZnrzR")
	if err != nil {
		yearElem, err = element.Element("a > article > div.StyledItemContent__StyleItemContentContainer-sc-1fqlnst-1.hzyHqR > div > div > div.StyledItemContent__FlexDiv-sc-1fqlnst-0.StyledItemContent__StyleItemContentContentDetails-sc-1fqlnst-5.kPOPyZ.jCmSiq > span.StyledItemContent__StyleItemContentContentDetailsDate-sc-1fqlnst-8.bZnrzR")
		if err != nil {
			log.Println("Get Year ", err)
			return nil, err
		}
	}
	yearStr, err := yearElem.Text()
	if err != nil {
		log.Println("Get Year Invalid element ", err)
		return nil, err
	}
	year, err := strconv.Atoi(yearStr)
	if err != nil {
		log.Println("Get Year ", err)
		return nil, err
	}
	return &year, nil
}

func getPrice(element *rod.Element) (*int, error) {
	priceElem, err := element.Element("a > article > div.StyledItemContent__StyleItemContentContainer-sc-1fqlnst-1.SpzQV > div > div.StyledItemContent__FlexDiv-sc-1fqlnst-0.StyledItemContent__StyleItemContentContentDetails-sc-1fqlnst-5.kPOPyZ.jCmSiq > data > span")
	if err != nil {
		priceElem, err = element.Element("a > article > div.StyledItemContent__StyleItemContentContainer-sc-1fqlnst-1.SpzQV > div > div.StyledItemContent__FlexDiv-sc-1fqlnst-0.StyledItemContent__StyleItemContentContentDetails-sc-1fqlnst-5.kPOPyZ.jCmSiq > span.StyledItemContent__StyleItemContentContentDetailsPriceNotAvailable-sc-1fqlnst-7.gcuxWP > i")
		if err != nil {
			priceElem, err = element.Element("a > article > div.StyledItemContent__StyleItemContentContainer-sc-1fqlnst-1.hzyHqR > div > div > div.StyledItemContent__FlexDiv-sc-1fqlnst-0.StyledItemContent__StyleItemContentContentDetails-sc-1fqlnst-5.kPOPyZ.jCmSiq > data > span")
			if err != nil {
				log.Println("Get PRice", err)
				return nil, err
			}
		}
	}
	priceStr := priceElem.MustText()
	leng := len(priceStr)
	if leng == 0 {
		return nil, errors.New(fmt.Sprintf("Something wrong with the price %s", priceStr))
	}
	//log.Println(priceStr)
	priceStr = priceStr[4:leng]
	priceStr = strings.Replace(priceStr, ".", "", -1)
	price, err := strconv.Atoi(priceStr)
	if err != nil {
		return nil, err
	}
	return &price, nil
}

func getDealer(element *rod.Element) string {
	dealerElem, err := element.Element("a > article > div.StyledItemFooter-sc-lcxh44-1.bsiiwc > div > div.StyledItemFooter__StyledItemFooterProvider-sc-lcxh44-3.dBkJUP > div > div > div > div > div:nth-child(1) > strong")
	if err != nil {
		dealerElem, err = element.Element("a > article > div.StyledItemFooter-sc-lcxh44-1.bsiiwc > div > div.StyledItemFooter__StyledItemFooterProvider-sc-lcxh44-3.dBkJUP > div > div > div > div.jsx-2031032736.seller__name-place > div:nth-child(1) > strong")
		if err != nil {
			dealerElem, err = element.Element("a > article > div.StyledItemFooter-sc-lcxh44-1.cxwLQI > div > div.jsx-2031032736.seller > div > div > div > div:nth-child(1) > strong")
			if err != nil {
				dealerElem, err = element.Element("a > article > div.StyledItemFooter-sc-lcxh44-1.bsiiwc > div > div.StyledItemFooter__StyledItemFooterProvider-sc-lcxh44-3.dBkJUP > div > div.jsx-1115023392.seller__image__details")
				if err != nil {
					log.Println("GetDealer: ", err)
					return "Dealer problem"
				}
			}
		}
	}
	return dealerElem.MustText()
}

func getAdID(href string) string {
	res := strings.Split(href, "-")
	adID := res[len(res)-1]
	return adID
}

func getHref(articleElement *rod.Element) (*string, error) {
	hrefElem, err := articleElement.Element("a[data-testid=result-list-item-link]")
	if err != nil {
		log.Println("Get HREF : ", err)
	}
	href, err := hrefElem.Attribute("href")
	if err != nil {
		log.Println("Get HREF Get attribute ", err)
		return nil, err
	}

	res := fmt.Sprintf("https://autotrack.nl%s", *href)
	return &res, nil
}

func isLastPage(paginationElement *rod.Element) (*bool, error) {
	isLast := false
	lastPageElement, err := paginationElement.Element("li:last-child")
	if err != nil {
		return nil, err
	}

	lastPageElementClass := lastPageElement.MustAttribute("class")
	if strings.Contains(*lastPageElementClass, "disabled") {
		isLast = true
		return &isLast, nil
	}

	return &isLast, nil
}

func acceptCookies(page *rod.Page) bool {
	modalContainer, err := page.Element("#pg-shadow-host")
	if err != nil {
		panic(err)
	}
	modal, err := modalContainer.ShadowRoot()
	if err != nil {
		panic(err)
	}
	acceptCookiesBtn, err := modal.Element("#pg-accept-btn")
	if err != nil {
		panic(err)
	}

	err = acceptCookiesBtn.Click(proto.InputMouseButtonLeft, 1)
	if err != nil {
		panic(err)
	}
	return true
}
