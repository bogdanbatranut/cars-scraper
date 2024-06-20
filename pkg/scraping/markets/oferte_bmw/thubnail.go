package oferte_bmw

import (
	"fmt"
	"log"
	"strings"

	"github.com/gocolly/colly"
)

func getThumbNail(adId int, calculated704pathElement string) *string {
	url := fmt.Sprintf("https://oferte.bmw.ro/rulate/cauta/detaliu/%d/", adId)
	var thumbNail string
	// Instantiate default collector
	c := colly.NewCollector()

	// On every a element which has href attribute call callback
	//#v138109-gallery-inner > div > div.swiper-slide.carousel-item-mosaic.swiper-slide-active > a:nth-child(1) > img
	c.OnHTML("div.swiper-container > div > div > a > img", func(e *colly.HTMLElement) {

		srcset := e.Attr("data-srcset")
		links := strings.Split(srcset, ",")
		//for _, link := range links {
		//	log.Println(link)
		//}
		//log.Println(" ------------------ ")
		thumbNail = strings.Split(links[0], " ")[0]
		thumbNail = fmt.Sprintf("https://oferte.bmw.ro/rulate/api/v1/ems/bmw-used-ro_RO/vehicle/704/%s/%d-0?2024-05-31T14:14:20+00:00", calculated704pathElement, adId)

	})

	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	// Start scraping on https://hackerspaces.org
	err := c.Visit(url)
	if err != nil {
		return nil
	}
	log.Println("THUMB : ", thumbNail)
	return &thumbNail
}

func get704paramForThumbnailURL() *string {
	url := "https://oferte.bmw.ro/rulate/cauta"
	c := colly.NewCollector()

	var calculatedParam string

	selector := "#search > section.dynamic > div > article"
	selector = "#search > section.dynamic > div > article > a.item-content > div.item-media > img"
	c.OnHTML(selector, func(e *colly.HTMLElement) {

		srcset := e.Attr(":data-srcset")
		linkStringElements := strings.Split(strings.Split(srcset, "+")[0], "/")
		calculatedParam = linkStringElements[10]
	})
	err := c.Visit(url)
	if err != nil {
		log.Println(err)
		return nil
	}
	return &calculatedParam
}
