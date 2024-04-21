package main

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/proto"
	"github.com/go-rod/rod/lib/utils"
	"github.com/ysmood/gson"
)

func main() {
	launch := launcher.New().
		Headless(false).
		Devtools(true)

	defer launch.Cleanup()

	url := launch.MustLaunch()

	browser := rod.New().ControlURL(url).
		Trace(true).
		SlowMotion(2 * time.Second).
		MustConnect()

	launcher.Open(browser.ServeMonitor(""))

	defer browser.MustClose()
	cookiesAccepted := false

	for i := 1; i < 3; i++ {
		log.Println("Page : ", i)
		url := fmt.Sprintf("https://www.autotrack.nl/aanbod?minimumbouwjaar=2019&maximumkilometerstand=125000&brandstofsoorten=DIESEL&modelIds=e43d83ec-00d4-4cfe-915d-231d267f6d02&merkIds=1a67a3d8-178b-43ee-9071-9ae7f19b316a&paginanummer=%d&paginagrootte=30&sortering=PRIJS_OPLOPEND", i)
		page := browser.MustPage(url)
		if !cookiesAccepted {
			page = page.MustWaitLoad()
			wait := page.MustWaitRequestIdle()
			wait()
			cookiesAccepted = acceptCookies(page)
			//articles := page.MustWaitDOMStable().MustElements("li[data-testid=result-list-item]")
			articles := page.MustWaitLoad().MustWaitDOMStable().MustElements("li[data-testid=result-list-item]")
			paginationElement := page.MustElement("#__next > div > main > div.main.container.pt-3 > div > div.jsx-2926786665.col-lg-9.filters__results-wrapper > div.jsx-2926786665.filters__results > section > div.jsx-3584031782.pagination__container > ul")
			processElements(articles)
			isLPage, err := isLastPage(paginationElement)
			if err != nil {

			}
			log.Println("IS LAST PAGE : ", isLPage)
			log.Println("Found articles 1: ", len(articles))
		} else {
			//wait := page.MustWaitRequestIdle()
			//wait()
			articles := page.MustWaitLoad().MustWaitDOMStable().MustElements("li[data-testid=result-list-item]")
			processElements(articles)
			paginationElement := page.MustElement("#__next > div > main > div.main.container.pt-3 > div > div.jsx-2926786665.col-lg-9.filters__results-wrapper > div.jsx-2926786665.filters__results > section > div.jsx-3584031782.pagination__container > ul")
			isLPage, err := isLastPage(paginationElement)
			if err != nil {
				panic(err)
			}
			log.Println("IS LAST PAGE : ", isLPage)
			log.Println("Found articles: ", len(articles))
		}
		//wait()

	}
	fmt.Println("done")
}

func processElements(elements rod.Elements) {
	for _, element := range elements {
		soldElement, err := element.Element("span.ItemTag__ItemTagStatusLabel-sc-hnf0c8-4.ItemTag__VerkochtStatus-sc-hnf0c8-5.SKbfR.ipKrWP")
		if err != nil {
			log.Println(err.Error())
		}
		if soldElement != nil {
			continue
		}
		adTitle := element.MustElement("div.StyledItemContent__StyleItemContentContainer-sc-1fqlnst-1.SpzQV > div > div.StyledItemContent__StyleItemContentContentHeader-sc-1fqlnst-3.dVpSsc > h3")
		adTitleText := adTitle.MustText()
		log.Println("Title: ", adTitleText)
		href := getHref(element)
		log.Println("HREF: ", href)
		log.Println("Ad ID: ", getAdID(href))
		log.Println("Dealer: ", getDealer(element))
		price, err := getPrice(element)
		if err != nil {
			// invalid price so continue
			continue
		}
		log.Println("Price: ", price)

		year, err := getYear(element)
		if err != nil {
			panic(err)
		}
		log.Println("Year : ", year)

		km, err := getKm(element)
		if err != nil {
			panic(err)
		}
		log.Println("KM : ", km)

		thumbNail := getThumbnail(element)
		log.Println("Tumbnail: ", thumbNail)
	}
}

func getThumbnail(element *rod.Element) string {
	thumbNailElement := element.MustElement("a > article > div.StyledItemContent__StyleItemContentContainer-sc-1fqlnst-1.SpzQV > figure > div > div > div > div.CarouselWrapperComponent__Slides-sc-1o127aw-6.kroIPK > div:nth-child(1) > img")
	thumbNailValue := thumbNailElement.MustAttribute("src")
	return *thumbNailValue
}

func getKm(element *rod.Element) (*int, error) {
	kmStr := element.MustElement("a > article > div.StyledItemContent__StyleItemContentContainer-sc-1fqlnst-1.SpzQV > div > div.StyledItemContent__FlexDiv-sc-1fqlnst-0.StyledItemContent__StyleItemContentContentDetails-sc-1fqlnst-5.kPOPyZ.jCmSiq > span.StyledItemContent__StyleItemContentContentDetailsMileage-sc-1fqlnst-10.kqpAaS").MustText()
	kmStr = kmStr[0 : len(kmStr)-3]
	kmStr = strings.Replace(kmStr, ".", "", -1)
	km, err := strconv.Atoi(kmStr)
	if err != nil {
		return nil, err
	}
	return &km, nil
}

func getYear(element *rod.Element) (*int, error) {
	yearStr := element.MustElement("a > article > div.StyledItemContent__StyleItemContentContainer-sc-1fqlnst-1.SpzQV > div > div.StyledItemContent__FlexDiv-sc-1fqlnst-0.StyledItemContent__StyleItemContentContentDetails-sc-1fqlnst-5.kPOPyZ.jCmSiq > span.StyledItemContent__StyleItemContentContentDetailsDate-sc-1fqlnst-8.bZnrzR").MustText()
	year, err := strconv.Atoi(yearStr)
	if err != nil {
		return nil, err
	}
	return &year, nil
}

func getPrice(element *rod.Element) (*int, error) {
	priceElem, err := element.Element("a > article > div.StyledItemContent__StyleItemContentContainer-sc-1fqlnst-1.SpzQV > div > div.StyledItemContent__FlexDiv-sc-1fqlnst-0.StyledItemContent__StyleItemContentContentDetails-sc-1fqlnst-5.kPOPyZ.jCmSiq > data > span")
	if err != nil {
		priceElem, err = element.Element("a > article > div.StyledItemContent__StyleItemContentContainer-sc-1fqlnst-1.SpzQV > div > div.StyledItemContent__FlexDiv-sc-1fqlnst-0.StyledItemContent__StyleItemContentContentDetails-sc-1fqlnst-5.kPOPyZ.jCmSiq > span.StyledItemContent__StyleItemContentContentDetailsPriceNotAvailable-sc-1fqlnst-7.gcuxWP > i")
		if err != nil {
			log.Println(err)
			return nil, err
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
		log.Println(err)
		dealerElem, err = element.Element("a > article > div.StyledItemFooter-sc-lcxh44-1.bsiiwc > div > div.StyledItemFooter__StyledItemFooterProvider-sc-lcxh44-3.dBkJUP > div > div > div > div.jsx-2031032736.seller__name-place > div:nth-child(1) > strong")
		if err != nil {
			log.Println(err)
			dealerElem, err = element.Element("a > article > div.StyledItemFooter-sc-lcxh44-1.bsiiwc > div > div.StyledItemFooter__StyledItemFooterProvider-sc-lcxh44-3.dBkJUP > div > div.jsx-1115023392.seller__image__details")
			if err != nil {
				log.Println(err)
				return "Dealer problem"
			}
		}
	}
	//dealerElem := element.MustElement("a > article > div.StyledItemFooter-sc-lcxh44-1.bsiiwc > div > div.StyledItemFooter__StyledItemFooterProvider-sc-lcxh44-3.dBkJUP > div > div > div > div > div:nth-child(1) > strong")
	return dealerElem.MustText()
}

func getAdID(href string) string {
	res := strings.Split(href, "-")
	adID := res[len(res)-1]
	return adID
}

func getHref(articleElement *rod.Element) string {
	href := articleElement.MustElement("a[data-testid=result-list-item-link]").MustAttribute("href")

	res := fmt.Sprintf("https://autotrack.nl%s", *href)
	return res
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

func ScreenShot(page *rod.Page, fileName string) {

	// capture entire browser viewport, returning jpg with quality=90
	img, err := page.ScrollScreenshot(&rod.ScrollScreenshotOptions{
		Format:  proto.PageCaptureScreenshotFormatJpeg,
		Quality: gson.Int(90),
	})
	if err != nil {
		panic(err)
	}

	_ = utils.OutputFile(fmt.Sprintf("%s.jpg", fileName), img)
}
