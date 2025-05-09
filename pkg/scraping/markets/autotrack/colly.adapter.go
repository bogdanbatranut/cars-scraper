package autotrack

import (
	"carscraper/pkg/jobs"
	"carscraper/pkg/scraping/icollector"
	"carscraper/pkg/scraping/markets/autotrack/autotrackcollycollector"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
)

type AutotrackCollyAdapter struct {
}

func NewAutotrackCollyAdataper() *AutotrackCollyAdapter {
	return &AutotrackCollyAdapter{}
}

//func getAdHref(element *colly.HTMLElement)

func (a AutotrackCollyAdapter) GetAds(job jobs.SessionJob) icollector.AdsResults {

	foundAds := []jobs.Ad{}

	isLastPage := false

	var executionErr error

	builder := NewURLBuilder()
	url := builder.GetURL(job)
	if url == nil {
		return icollector.AdsResults{
			Ads:        nil,
			IsLastPage: true,
			Error:      nil,
		}
		//return nil, true, nil
	}

	TestGETRequest(*url)

	autotrackCollyCollector := autotrackcollycollector.NewAutotrackCollyCollector().GetCollyCollector(job)

	autotrackCollyCollector.OnHTML("li.pagination__next.pagination__disabled", func(element *colly.HTMLElement) {
		log.Println("TESTING FOR LAST PAGE ")
		if element != nil {
			isLastPage = true
		}
	})

	autotrackCollyCollector.OnHTML("li[data-testid=result-list-item]", func(element *colly.HTMLElement) {
		soldElement := element.DOM.Find("article").Find("span.ItemTag__ItemTagStatusLabel-sc-hnf0c8-4.ItemTag__VerkochtStatus-sc-hnf0c8-5.SKbfR.ipKrWP")
		if soldElement.Nodes != nil {
			return
		}

		href, _ := element.DOM.Find("a[data-testid=result-list-item-link]").Attr("href")
		//log.Println(exists, href)

		res := strings.Split(href, "-")

		adID := res[len(res)-1]
		//log.Println(adID)
		href = fmt.Sprintf("https://autotrack.nl%s", href)

		dealer := element.DOM.Find("article").Find("div.StyledItemFooter__StyledItemFooterCtaSellerRight-sc-lcxh44-10.iucllX > div > div > div > div > div:nth-child(1) > strong").Text()
		//log.Println(dealer)

		priceStr := element.DOM.Find("article").Find("div.StyledItemContent__StyleItemContentContainer-sc-1fqlnst-1.SpzQV > div > div.StyledItemContent__FlexDiv-sc-1fqlnst-0.StyledItemContent__StyleItemContentContentDetails-sc-1fqlnst-5.kPOPyZ.jCmSiq > data > span").Text()
		leng := len(priceStr)
		if leng == 0 {
			log.Println("FOUND ADS SO Far: ", len(foundAds))
			return
		}
		//log.Println(priceStr)
		priceStr = priceStr[4:leng]
		priceStr = strings.Replace(priceStr, ".", "", -1)
		price, err := strconv.Atoi(priceStr)
		if err != nil {
			executionErr = err
		}

		//log.Println(price)

		yearStr := element.DOM.Find("article").Find("div.StyledItemContent__StyleItemContentContainer-sc-1fqlnst-1.SpzQV > div > div.StyledItemContent__FlexDiv-sc-1fqlnst-0.StyledItemContent__StyleItemContentContentDetails-sc-1fqlnst-5.kPOPyZ.jCmSiq > span.StyledItemContent__StyleItemContentContentDetailsDate-sc-1fqlnst-8.bZnrzR").Text()
		//log.Println(yearStr)
		year, err := strconv.Atoi(yearStr)
		if err != nil {
			executionErr = err
		}

		kmStr := element.DOM.Find("article").Find("div.StyledItemContent__StyleItemContentContainer-sc-1fqlnst-1.SpzQV > div > div.StyledItemContent__FlexDiv-sc-1fqlnst-0.StyledItemContent__StyleItemContentContentDetails-sc-1fqlnst-5.kPOPyZ.jCmSiq > span.StyledItemContent__StyleItemContentContentDetailsMileage-sc-1fqlnst-10.kqpAaS").Text()
		//log.Println(kmStr)
		kmStr = kmStr[0 : len(kmStr)-3]
		kmStr = strings.Replace(kmStr, ".", "", -1)
		km, err := strconv.Atoi(kmStr)
		if err != nil {
			executionErr = err
		}

		thumbnailElement := element.DOM.Find("div.StyledItemContent__StyleItemContentContainer-sc-1fqlnst-1.SpzQV > figure > div > div > div > div.CarouselWrapperComponent__Slides-sc-1o127aw-6.kroIPK > div:nth-child(1) > img")
		thumbNailValue, _ := thumbnailElement.Attr("src")
		//log.Println(thumbNailValue, exists)
		//log.Println("----- ")

		carad := jobs.Ad{
			Brand:              job.Criteria.Brand,
			Model:              job.Criteria.CarModel,
			Year:               year,
			Km:                 km,
			Fuel:               job.Criteria.Fuel,
			Price:              price,
			AdID:               adID,
			Ad_url:             href,
			SellerType:         "dealer",
			SellerName:         &dealer,
			SellerNameInMarket: &dealer,
			SellerOwnURL:       &dealer,
			SellerMarketURL:    &dealer,
			Thumbnail:          &thumbNailValue,
		}
		foundAds = append(foundAds, carad)
	})

	if executionErr != nil {
		return icollector.AdsResults{
			Ads:        nil,
			IsLastPage: true,
			Error:      executionErr,
		}
	}

	err := autotrackCollyCollector.Visit(*url)
	if err != nil {
		return icollector.AdsResults{
			Ads:        nil,
			IsLastPage: true,
			Error:      err,
		}
	}
	log.Println("AUTOTRACK Visiting ", *url)

	autotrackCollyCollector.Wait()
	if len(foundAds) == 0 {
		log.Println("NO RESULTS SO RETURN !!!!!")
		return icollector.AdsResults{
			Ads:        nil,
			IsLastPage: true,
			Error:      nil,
		}
	}

	log.Println("AUTOTRACK found ads : ", len(foundAds))
	return icollector.AdsResults{
		Ads:        &foundAds,
		IsLastPage: isLastPage,
		Error:      nil,
	}

}

func TestGETRequest(url string) {

	httpMethod := "GET"
	httpClient := &http.Client{}
	httpRequest, err := http.NewRequest(httpMethod, url, nil)

	if err != nil {
		panic(err)
	}

	response, err := httpClient.Do(httpRequest)
	log.Println("Status code : ", response.StatusCode)
	if err != nil {
		log.Printf("got response with error: %+v", err)
	}
	defer response.Body.Close()
	bodyBytes, err := io.ReadAll(response.Body)
	fmt.Println(string(bodyBytes))

	//return bodyBytes, url, nil
}
