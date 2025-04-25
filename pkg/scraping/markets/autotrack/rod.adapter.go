package autotrack

import (
	"carscraper/pkg/jobs"
	"carscraper/pkg/scraping/icollector"
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
	"github.com/go-rod/rod/lib/utils"
	"github.com/ysmood/gson"
)

type AutoTrackNLRodAdapter struct{}

func NewAutoTrackNLRodAdapter() *AutoTrackNLRodAdapter { return &AutoTrackNLRodAdapter{} }

func ScreenShot(page *rod.Page) {

	// capture entire browser viewport, returning jpg with quality=90
	img, err := page.ScrollScreenshot(&rod.ScrollScreenshotOptions{
		Format:  proto.PageCaptureScreenshotFormatJpeg,
		Quality: gson.Int(90),
	})
	if err != nil {
		panic(err)
	}

	_ = utils.OutputFile("rod_scraping.jpg", img)
}

func accessDenied(page *rod.Page) bool {
	ctx, cancel := context.WithCancel(context.Background())
	pageWithCancel := page.Context(ctx)

	go func() {
		time.Sleep(2 * time.Second)
		cancel()
	}()

	elem, err := pageWithCancel.Element("h1")
	if err != nil {
		return false
	}
	//text := elem.MustText()
	text, err := elem.Text()
	if err != nil {
		log.Println("access denied ")
		return false
	}
	if strings.Contains(text, "Denied") {
		return true
	}
	return false

}

func checkForCookiesAndAccept(page *rod.Page) {
	modalContainer, err := page.Element("#pg-host-shadow-root")
	err = page.WaitDOMStable(time.Second*5, 0)
	if err != nil {
		panic(err)
	}
	modal, err := modalContainer.ShadowRoot()
	if err != nil {
		panic(err)
	}

	acceptCookiesBtn, err := modal.Element("#pg-accept-btn")
	if err != nil {
		log.Println(err, "AUTOTRACK.NL accept cookie button")
		return
	}
	log.Println("ACCEPT COOKIE BUTTON FOUND")

	log.Println(acceptCookiesBtn.MustText())
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

	if accessDenied(page) {
		//ScreenShot(page)
		return &icollector.AdsResults{
			Ads:        nil,
			IsLastPage: true,
			Error:      errors.New("AUTOTRACK PAGE ACCESS DENIED"),
		}
	}

	checkForCookiesAndAccept(page)

	artilcesSelector := "body > div.min-h-screen.flex.flex-col.justify-between.bg-base-white > div > div.px-3.xl\\:px-10.sm\\:px-4.flex.flex-col.items-center.w-full.relative > div > section > div.w-full > div.grid.gap-6 > div.grid.grid-cols-1.gap-3.mt-6.sm\\:mt-0.md\\:mt-6"

	articlesContainer, err := page.MustWaitLoad().MustWaitDOMStable().Element(artilcesSelector)

	if err != nil {
		return &icollector.AdsResults{
			Ads:        nil,
			IsLastPage: true,
			Error:      nil,
		}
	}

	var articles []*rod.Element

	for i := 0; i < 20; i++ {
		// Wait until all network requests have finished
		wait := page.Timeout(time.Second*4).WaitRequestIdle(1000*time.Millisecond, nil, nil, nil)
		wait()
		page.Mouse.MustScroll(0, 700)

	}

	wait := page.WaitRequestIdle(1000*time.Millisecond, nil, nil, nil)
	wait()

	arts, err := articlesContainer.Elements("div.w-full > div > a")
	if err != nil {
		return &icollector.AdsResults{
			Ads:        nil,
			IsLastPage: true,
			Error:      nil,
		}
	}

	if !arts.Empty() {
		counter := 1
		for _, art := range arts {
			articles = append(articles, art)
			counter++
		}
	}

	if len(articles) > 0 && len(articles) < 90 {
		return &icollector.AdsResults{
			Ads:        processElements(articles),
			IsLastPage: true,
			Error:      nil,
		}
	}

	return &icollector.AdsResults{
		Ads:        processElements(articles),
		IsLastPage: false,
		Error:      nil,
	}

}

func processElements(elements rod.Elements) *[]jobs.Ad {
	var ads []jobs.Ad
	for _, element := range elements {
		soldElement, err := element.Element("span.ItemTag__ItemTagStatusLabel-sc-hnf0c8-4.ItemTag__VerkochtStatus-sc-hnf0c8-5.SKbfR.ipKrWP")
		if soldElement != nil {
			continue
		}
		soldElement, err = element.Element("div.StyledItemContent__StyleItemContentContainer-sc-1fqlnst-1.hzyHqR > div > div > div.StyledItemTags-sc-46f8tm-0.hRDWcg > span")
		//}
		if soldElement != nil {
			continue
		}
		adTitle, err := element.Element("div > div > div.flex.flex-col.flex-1.sm\\:flex-row.mx-auto.sm\\:mx-0.sm\\:w-auto.gap-5.sm\\:gap-0.w-full > div.flex-col.flex.justify-between.w-full.self-center.shrink.xl\\:h-full > div.mx-3.sm\\:mx-4.py-2.gap-3.grid > div")

		adTitleText := strings.Replace(adTitle.MustText(), "\n\n", " ", -1)

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
			Title:              &adTitleText,
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
	thumbNailElement, err := element.Element("div > div > div > div > div > div > div > div > img")
	if err != nil {
		log.Println("Get THUMBNAIL ", err)
		return nil, err
	}
	thumbNailValue, err := thumbNailElement.Attribute("src")
	if err != nil {
		log.Println("Get THUMBN attribute src ", err)
		return nil, err
	}
	return thumbNailValue, nil
}

func getKm(element *rod.Element) (*int, error) {
	kmElem, err := element.Element("div > div > div > div.flex-col > div.mx-3 > div.grid > div > p:nth-child(5)")
	if err != nil {
		log.Println("Get KM invalid selector ", err)
		return nil, err
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
	yearElem, err := element.Element("div > div > div > div.flex-col > div.mx-3 > div.grid > div > p:nth-child(3)")
	if err != nil {
		log.Println("Get Year ", err)
		return nil, err
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
	priceElem, err := element.Element("div > div > div > div.flex-col > div.mx-3 > div.grid > div > div")
	if err != nil {
		return nil, err
	}
	priceStr := priceElem.MustElement("div > p").MustText()
	leng := len(priceStr)
	if leng == 0 {
		return nil, errors.New(fmt.Sprintf("Something wrong with the price %s", priceStr))
	}
	//log.Println(priceStr)
	priceStr = strings.Replace(priceStr, ".", "", -1)
	priceStr = strings.Replace(priceStr, "€", "", -1)
	priceStr = strings.Replace(priceStr, " ", "", -1)

	price, err := strconv.Atoi(priceStr)
	if err != nil {
		log.Println("cannot convert pricestr ", priceStr)
		return nil, err
	}

	vatElemem, err := priceElem.Element("p.text-xxs")
	if err != nil {
		return &price, nil
	}
	if strings.Contains(vatElemem.MustText(), "excl") {
		price = price + (price * 21 / 100)
	}

	if err == nil {
		return &price, nil
	}

	priceStr = priceStr[4:leng]

	price, err = strconv.Atoi(priceStr)
	if err != nil {
		return nil, err
	}

	return &price, nil
}

func getDealer(element *rod.Element) string {
	dealerElem, err := element.Element("div > div > div > div.flex-col > div.border-t > div > div > div > div > p")
	if err != nil {
		log.Println("Could not find exact dealer ")
		return "autotrack.nl"
	}
	txt := dealerElem.MustText()
	return txt
}

func getAdID(href string) string {
	res := strings.Split(href, "-")
	adID := res[len(res)-1]
	return adID
}

func getHref(articleElement *rod.Element) (*string, error) {
	href, err := articleElement.Attribute("href")
	if err != nil {
		log.Println("Get HREF : ", err)
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
