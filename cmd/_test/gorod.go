package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math"
	"strconv"
	"strings"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
)

func main() {
	url := "https://www.mobile.de/ro/automobil/mercedes-benz-clasa-gle/vhc:car,srt:price,sro:asc,ms1:17200_-58_,frn:2019,ful:diesel,mlx:125000"
	GetData(url)
}

func GetData(url string) {
	//isLastPage := true

	page := rod.New().MustConnect().MustPage("https://www.mobile.de/ro/automobil/mercedes-benz-clasa-gle/vhc:car,srt:price,sro:asc,ms1:17200_-58_,frn:2019,ful:diesel,mlx:125000").
		MustWaitLoad()
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
	totalResultsElement := page.MustSearch("body > div.g-content > div > div.u-display-flex.u-margin-top-18 > section > section.result-block-header.g-row > div > h1")
	totalResultsElementText, err := totalResultsElement.Text()
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
