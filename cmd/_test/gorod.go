package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/proto"
	"github.com/go-rod/rod/lib/utils"
	"github.com/ysmood/gson"
)

//func main() {
//	url := "https://www.mobile.de/ro/automobil/mercedes-benz-clasa-gle/vhc:car,srt:price,sro:asc,ms1:17200_-58_,frn:2019,ful:diesel,mlx:125000"
//	GetData(url)
//}

func mainIII() {
	// Devtools opens the tab in each new tab opened automatically
	l := launcher.New().
		Headless(false).
		Devtools(true)

	defer l.Cleanup()

	url := l.MustLaunch()

	// Trace shows verbose debug information for each action executed
	// SlowMotion is a debug related function that waits 2 seconds between
	// each action, making it easier to inspect what your code is doing.
	browser := rod.New().
		ControlURL(url).
		Trace(true).
		SlowMotion(2 * time.Second).
		MustConnect()

	// ServeMonitor plays screenshots of each tab. This feature is extremely
	// useful when debugging with headless mode.
	// You can also enable it with flag "-rod=monitor"
	launcher.Open(browser.ServeMonitor(""))

	defer browser.MustClose()

	page := browser.MustPage("https://www.autotrack.nl/aanbod?minimumbouwjaar=2019&maximumkilometerstand=125000&brandstofsoorten=DIESEL&modelIds=e43d83ec-00d4-4cfe-915d-231d267f6d02&merkIds=1a67a3d8-178b-43ee-9071-9ae7f19b316a&paginanummer=1&paginagrootte=30&sortering=PRIJS_OPLOPEND")
	log.Println("Wait for request.idle")
	wait := page.WaitRequestIdle(time.Second*2, nil, nil, nil)
	wait()

	log.Println("Sleeeping")
	time.Sleep(10 * time.Second)
	log.Println("Waking")

	buttons, err := page.Elements("button")
	if err != nil {
		panic(err)
	}

	for _, button := range buttons {
		log.Println(button)
	}

	html, err := page.HTML()
	if err != nil {
		panic(err)
	}

	log.Println(html)
	elem, err := page.Element("#pg-shadow-host")
	if err != nil {
		panic(err)
	}

	log.Println(elem.MustElement("#shadow-root"))

	//butt, err := page.Element("#pg-accept-btn")
	//if err != nil {
	//	panic(err)
	//}
	//log.Println(butt.Text())

	//page.MustElement("input").MustInput("git").MustType(input.Enter)
	//
	//text := page.MustElement(".codesearch-results p").MustText()

	//fmt.Println(text)

	utils.Pause() // pause goroutine
}

func main() {
	l := launcher.MustNewManaged("ws://dev.auto-mall.ro:7317")

	// You can also set any flag remotely before you launch the remote browser.
	// Available flags: https://peter.sh/experiments/chromium-command-line-switches
	l.Set("disable-gpu").Delete("disable-gpu")

	// Launch with headful mode
	l.Headless(true).XVFB("--server-num=5", "--server-args=-screen 0 1600x900x16")

	browser := rod.New().Client(l.MustClient()).MustConnect()

	// You may want to start a server to watch the screenshots of the remote browser.
	launcher.Open(browser.ServeMonitor(""))
	page, err := browser.Page(proto.TargetCreateTarget{})
	if err != nil {
		panic(err)
	}

	headers := []string{
		//"Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7",
		//"Accept-Encoding", "gzip, deflate, br",
		//"Accept-Language", "en-GB,en;q=0.9",
		//"Sec-Ch-Ua", "\"Google Chrome\";v=\"119\", \"Chromium\";v=\"119\", \"Not?A_Brand\";v=\"24\"",
		//"Sec-Ch-Ua-Mobile", "?0",
		//"Sec-Ch-Ua-Platform", "\"macOS\"",
		//"Sec-Fetch-Dest", "document",
		//"Sec-Fetch-Mode", "navigate",
		//"Sec-Fetch-Site", "none",
		//"Sec-Fetch-User", "?1",
		//"Upgrade-Insecure-Requests", "1",
		//"User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36",
	}

	result, err := page.SetExtraHeaders(headers)
	if err != nil {
		panic(err)
	}

	result()
	//extraHeaders := page.MustSetExtraHeaders(
	//	"Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7",
	//	"Accept-Encoding", "gzip, deflate, br",
	//	"Accept-Language", "en-GB,en;q=0.9",
	//	"Sec-Ch-Ua", "\"Google Chrome\";v=\"119\", \"Chromium\";v=\"119\", \"Not?A_Brand\";v=\"24\"",
	//	"Sec-Ch-Ua-Mobile", "?0",
	//	"Sec-Ch-Ua-Platform", "\"macOS\"",
	//	"Sec-Fetch-Dest", "document",
	//	"Sec-Fetch-Mode", "navigate",
	//	"Sec-Fetch-Site", "none",
	//	"Sec-Fetch-User", "?1",
	//	"Upgrade-Insecure-Requests", "1",
	//	"User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36",
	//)
	//extraHeaders()

	//err = page.Navigate("https://www.mobile.de/ro/automobil/mercedes-benz-clasa-gle/vhc:car,srt:price,sro:asc,ms1:17200_-58_,frn:2019,ful:diesel,mlx:125000")
	err = page.Navigate("https://www.autoscout24.ro/lst/mercedes-benz/glc-(toate)/ft_motorina?atype=C&cy=D%2CA%2CB%2CE%2CF%2CI%2CL%2CNL&damaged_listing=exclude&desc=0&fregfrom=2019&kmto=125000&page=1&powertype=kw&regfrom=2019&sort=standard&source=detailsearch&ustate=N%2CU")

	if err != nil {
		panic(err)
	}
	accCookiesElem := page.MustElement("#as24-cmp-popup > div > div._acceptance-buttons_1fb0r_82 > button._consent-accept_1fb0r_111")
	accCookiesElem.MustClick()

	articles := page.MustElements(".list-page-item")
	log.Println("Main articles: ", len(articles))
	main := page.MustElement("main")
	elements, err := main.Elements("article")
	if err != nil {
		panic(err)
	}
	log.Println("articles: ", len(elements))
	eleme := page.MustElements("article > div.ListItem_wrapper__TxHWu > div.Gallery_wrapper__iqp3u > section > div:nth-child(1) > picture > source")
	for _, e := range eleme {
		src, err := e.Attribute("srcset")
		if err != nil {
			panic(err)
		}
		log.Println(*src)
	}
	utils.Pause()
}
func mainOLD() {
	//	// This example is to launch a browser remotely, not connect to a running browser remotely,
	//	// to connect to a running browser check the "../connect-browser" example.
	//	// Rod provides a docker image for beginners, run the below to start a launcher.Manager:
	//	//
	//	//     docker run --rm -p 7317:7317 ghcr.io/go-rod/rod
	//	//
	//	// For available CLI flags run: docker run --rm ghcr.io/go-rod/rod rod-manager -h
	//	// For more information, check the doc of launcher.Manager
	//	l := launcher.MustNewManaged("http://dev.auto-mall.ro:7317")
	//
	//	// You can also set any flag remotely before you launch the remote browser.
	//	// Available flags: https://peter.sh/experiments/chromium-command-line-switches
	//	l.Set("disable-gpu").Delete("disable-gpu")
	//
	//	// Launch with headful mode
	//	l.Headless(false).XVFB("--server-num=5", "--server-args=-screen 0 1600x900x16")
	//
	//	browser := rod.New().Client(l.MustClient()).MustConnect()
	//
	//	// You may want to start a server to watch the screenshots of the remote browser.
	//	launcher.Open(browser.ServeMonitor(""))
	//
	//	fmt.Println(
	//		browser.MustPage("https://developer.mozilla.org").MustEval("() => document.title"),
	//	)
	//
	// Launch another browser with the same docker container.
	ll := launcher.MustNewManaged("http://dev.auto-mall.ro:7317") //.Headless(true)
	//ll.Open(browser.ServeMonitor(""))
	// You can set different flags for each browser.
	//ll.Set("disable-sync").Delete("disable-sync")

	//ll.Set("--headless")
	anotherBrowser := rod.New().Client(ll.MustClient()).MustConnect()

	//router := anotherBrowser.HijackRequests()
	//router.MustAdd()
	page, err := anotherBrowser.Page(proto.TargetCreateTarget{})
	if err != nil {
		panic(err)
	}

	extraHeaders := page.MustSetExtraHeaders(
		"Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7",
		"Accept-Encoding", "gzip, deflate, br",
		"Accept-Language", "en-GB,en;q=0.9",
		"Sec-Ch-Ua", "\"Google Chrome\";v=\"119\", \"Chromium\";v=\"119\", \"Not?A_Brand\";v=\"24\"",
		"Sec-Ch-Ua-Mobile", "?0",
		"Sec-Ch-Ua-Platform", "\"macOS\"",
		"Sec-Fetch-Dest", "document",
		"Sec-Fetch-Mode", "navigate",
		"Sec-Fetch-Site", "none",
		"Sec-Fetch-User", "?1",
		"Upgrade-Insecure-Requests", "1",
		"User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36",
	)
	if err != nil {
		panic(err)
	}

	var e proto.NetworkResponseReceived
	wait := page.WaitEvent(&e)
	err = page.Navigate("https://www.mobile.de/ro/automobil/mercedes-benz-clasa-gle/vhc:car,srt:price,sro:asc,ms1:17200_-58_,frn:2019,ful:diesel,mlx:125000")
	if err != nil {
		panic(err)
	}
	wait()
	log.Println("Response status: ", e.Response.Status)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	page.MustScreenshot("my.png")

	// customization version
	img, _ := page.Screenshot(true, &proto.PageCaptureScreenshot{
		Format:  proto.PageCaptureScreenshotFormatJpeg,
		Quality: gson.Int(90),
		Clip: &proto.PageViewport{
			X:      0,
			Y:      0,
			Width:  1800,
			Height: 3200,
			Scale:  1,
		},
		FromSurface: true,
	})
	_ = utils.OutputFile("my.jpg", img)

	defer cancel()
	//result := anotherBrowser.Context(ctx).MustPage("https://www.mobile.de/ro/automobil/mercedes-benz-clasa-gle/vhc:car,srt:price,sro:asc,ms1:17200_-58_,frn:2019,ful:diesel,mlx:125000").MustEval("() => document.querySelectorAll(\"article\")")
	result := anotherBrowser.Context(ctx).MustPage("https://www.mobile.de/ro/automobil/mercedes-benz-clasa-gle/vhc:car,srt:price,sro:asc,ms1:17200_-58_,frn:2019,ful:diesel,mlx:125000").MustElement("body")
	log.Println("Result : ", result)
	extraHeaders()
	//fmt.Println(
	//	"===>",
	//	//anotherBrowser.MustPage("https://go-rod.github.io").MustEval("() => document.title"),
	//	anotherBrowser.MustPage("https://www.mobile.de/ro/automobil/mercedes-benz-clasa-gle/vhc:car,srt:price,sro:asc,ms1:17200_-58_,frn:2019,ful:diesel,mlx:125000").MustEval("() => document.querySelectorAll(\"article\")"),
	//)

	//utils.Pause()
}

func GetData(url string) {
	//isLastPage := true
	l := launcher.MustNewManaged("http://dev.auto-mall.ro:7317")
	page := rod.New().Client(l.MustClient()).MustConnect().MustPage("https://www.mobile.de/ro/automobil/mercedes-benz-clasa-gle/vhc:car,srt:price,sro:asc,ms1:17200_-58_,frn:2019,ful:diesel,mlx:125000").
		MustWaitLoad()
	//page := rod.New().MustConnect().MustPage("https://www.mobile.de/ro/automobil/mercedes-benz-clasa-gle/vhc:car,srt:price,sro:asc,ms1:17200_-58_,frn:2019,ful:diesel,mlx:125000").
	//	MustWaitLoad()
	// LoadState detects whether the network domain is enabled or not.
	fmt.Println("******")
	fmt.Println(page.LoadState(&proto.PageStopLoading{}))
	// We use Go's standard way to check error types, no magic.
	check := func(err error) {
		var evalErr *rod.ErrEval
		if errors.Is(err, context.DeadlineExceeded) { // timeout error
			fmt.Println("timeout err")
		} else if errors.As(err, &evalErr) { // eval error
			fmt.Println(evalErr.LineNumber)
		} else if err != nil {
			fmt.Println("can't handle", err)
		}
	}

	var totalResults float64
	totalResultsElement, err := page.Search("body > div.g-content > div > div.u-display-flex.u-margin-top-18 > section > section.result-block-header.g-row > div > h1")
	if err != nil {
		panic(err)
	}
	totalResultsElementText, err := totalResultsElement.First.Text()
	if err != nil {
		panic(err)
	}
	spaceIndex := strings.Index(totalResultsElementText, " ")
	log.Printf("Total results : %s", totalResultsElementText)
	totalResultsStr := totalResultsElementText[:spaceIndex]
	totalResultsStr = strings.Replace(totalResultsStr, ".", "", -1)
	totalResults_, err := strconv.Atoi(totalResultsStr)
	if err != nil {
		panic(err)
		return
	}
	totalResults = float64(totalResults_)

	numberOfTotalPages := math.Ceil(totalResults / 50)

	log.Println("Total pages: ", numberOfTotalPages)
	//if float64(pageNumber) == numberOfTotalPages || totalResults == 0 {
	//	isLastPage = true
	//}

	// The two code blocks below are doing the same thing in two styles:

	// The block below is better for debugging or quick scripting. We use panic to short-circuit logics.
	// So that we can take advantage of fluent interface (https://en.wikipedia.org/wiki/Fluent_interface)
	// and fail-fast (https://en.wikipedia.org/wiki/Fail-fast).
	// This style will reduce code, but it may also catch extra errors (less consistent and precise).
	{
		page.MustSearch("button.mde-consent-accept-btn").MustClick()

		err := rod.Try(func() {
			elems := page.MustElements("article")
			for _, elem := range elems {
				log.Println(getSellerType(elem))
				log.Println(getAdIDandHREF(elem))
				log.Println(getAdThumbnail(elem))
				log.Println(getYearAndKM(elem))
				priceElem := elem.MustElement("p.seller-currency")
				log.Println(priceElem.Text())
			}

			//fmt.Println(page.MustElement("article").MustHTML()) // use "Must" prefixed functions

		})
		check(err)
	}

	// The block below is better for production code. It's the standard way to handle errors.
	// Usually, this style is more consistent and precise.
	{
		el, err := page.Element("a")
		if err != nil {
			check(err)
			return
		}
		html, err := el.HTML()
		if err != nil {
			check(err)
			return
		}
		fmt.Println(html)
	}
}

func getSellerType(e *rod.Element) string {
	se := e.MustElement("div > div.g-row.js-ad-entry > a > div.g-col-s-12.g-col-m-8 > div:nth-child(2) > div.u-text-grey-60.g-col-s-8.g-col-m-9.u-margin-bottom-9")
	sellerTypeStr, err := se.Text()
	if err != nil {
		panic(err)
	}
	if strings.ContainsAny(sellerTypeStr, "dealer") {
		return "dealer"
	} else {
		return "privat"
	}
}

func getAdIDandHREF(e *rod.Element) (string, string) {
	se := e.MustElement("div > div.g-row.js-ad-entry > a")
	mobileAdId, err := se.Attribute("data-vehicle-id")
	if err != nil {
		panic(err)
	}
	adHREF, err := se.Attribute("href")
	if err != nil {
		panic(err)
	}
	return *mobileAdId, *adHREF
}

func getAdThumbnail(e *rod.Element) string {
	se := e.MustElement("noscript")
	thumbNail, err := se.Text()
	thumbSrc := strings.Split(thumbNail, " ")[2]
	thumbNail = strings.Split(thumbSrc, "=")[1]
	rep := "\""
	thumbNail = strings.Replace(thumbNail, rep, "", -1)
	//se := e.MustElement("div > div.g-row.js-ad-entry > a > div.thumbnail > img")
	//thumbNail, err := se.Attribute("src")
	if err != nil {
		panic(err)
	}

	return thumbNail
}

func getYearAndKM(e *rod.Element) (int, int) {
	ykmElem := e.MustElement("div > div.g-row.js-ad-entry > a > div.g-col-s-12.g-col-m-8 > div.vehicle-text.g-row > div.vehicle-information.g-col-s-6.g-col-m-8 > p.u-text-bold")
	ykmElemtext, err := ykmElem.Text()
	if err != nil {
		panic(err)
	}

	yearStr := ""

	if len(ykmElemtext) >= 7 {
		if strings.ContainsAny(ykmElemtext, "/") {
			yearStr = ykmElemtext[3:7]
		} else {
			if strings.ContainsAny(ykmElemtext, ",") {
				yearStr = ykmElemtext[:5]
				yearStr = strings.Replace(yearStr, ",", "", -1)
				yearStr = strings.Replace(yearStr, " ", "", -1)
				yearStr = strings.Replace(yearStr, "\u00a0", "", -1)
			}
		}
	}

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		//executionErr = err
		panic(err)
	}
	kmStr := "0"
	if len(ykmElemtext) >= 7 {
		if strings.Contains(ykmElemtext, "/") {
			km_ := ykmElemtext[9:]
			kmStr = strings.Replace(km_, ".", "", -1)
			kmStr = strings.Replace(kmStr, "\u00a0", "", -1)
			kmStr = strings.Replace(kmStr, "km", "", -1)
		} else {
			kmStr = ykmElemtext[:3]
		}

	} else {
		kmStr := strings.Replace(ykmElemtext, ".", "", -1)
		kmStr = strings.Replace(kmStr, "\u00a0", "", -1)
		kmStr = strings.Replace(kmStr, "km", "", -1)
	}

	km, err := strconv.Atoi(kmStr)
	if err != nil {
		//executionErr = err
		panic(err)
	}
	return year, km
}

func getGrossPrice(e *rod.Element) int {
	priceElem := e.MustElement("div > div.g-row.js-ad-entry > a > div.g-col-s-12.g-col-m-8 > div.vehicle-text.g-row > div.g-col-s-6.g-col-m-4.u-text-right > div > p.seller-currency.u-text-bold")
	grossPriceStr, err := priceElem.Text()
	if err != nil {
		panic(err)
	}

	grossPrice, err := strconv.Atoi(grossPriceStr)
	if err != nil {
		//executionErr = err
		panic(err)
	}
	return grossPrice
}
