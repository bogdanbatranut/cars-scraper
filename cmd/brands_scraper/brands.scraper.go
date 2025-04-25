package main

import (
	"carscraper/pkg/initialization"
	"carscraper/pkg/vehicles"
	"fmt"
	"log"
	"strings"

	"github.com/gocolly/colly"
)

type BrandInfo struct {
	Name       string
	ModelsLink string
	Models     *[]Models
}

type Models struct {
	Name string
}

func main() {
	brands := getBrands()
	for _, brand := range brands {
		getModels(brand)
	}

	db, err := initialization.InitVehiclesDB()
	if err != nil {
		panic(err)
	}

	for _, brand := range brands {
		vehicleModels := []vehicles.Model{}
		for _, model := range *brand.Models {
			vehicleModel := vehicles.Model{
				Name: model.Name,
			}
			vehicleModels = append(vehicleModels, vehicleModel)
		}
		vehicleBrand := vehicles.Brand{
			Name:              brand.Name,
			Models:            vehicleModels,
			SupportingMarkets: nil,
		}

		tx := db.Create(&vehicleBrand)
		if tx.Error != nil {
			panic(tx.Error)
		}
	}

}

func getBrands() []*BrandInfo {
	brands := []*BrandInfo{}
	c := colly.NewCollector()

	c.OnHTML("#pagewrapper > div.container.carlist.clearfix > div.carman", func(element *colly.HTMLElement) {
		selection := element.DOM.Find("a:nth-child(1)")
		href, hasHref := selection.Attr("href")
		//log.Println(href, has)
		title, hasTitle := selection.Attr("title")
		//log.Println(title)
		if hasTitle && hasHref {
			brands = append(brands, &BrandInfo{
				Name:       title,
				ModelsLink: href,
			})
		}

	})

	err := c.Visit("https://www.autoevolution.com/cars/")
	if err != nil {
		return nil
	}
	c.Wait()
	return brands
}

func getModels(brand *BrandInfo) {
	c := colly.NewCollector(colly.Async(false))

	models := []Models{}

	c.OnHTML("#newscol2 > div > div.carmod", func(element *colly.HTMLElement) {

		selection := element.DOM.Find("a:nth-child(1) ")
		title := selection.Text()
		title = strings.Replace(title, fmt.Sprintf("%s ", brand.Name), "", -1)
		models = append(models, Models{Name: title})
		//brand.Models = append(brand.Models, Models{
		//	Name: title,
		//})
	})

	c.OnRequest(func(request *colly.Request) {
		log.Println("Visiting", request.URL.String())
	})
	err := c.Visit(brand.ModelsLink)
	if err != nil {
		panic(err)
	}
	c.Wait()

	brand.Models = &models
}
