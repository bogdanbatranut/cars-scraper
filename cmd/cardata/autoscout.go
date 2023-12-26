package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gocolly/colly"
)

type CarInfo struct {
	Brand  CarBrand
	Models CarModels
}

type CarBrand struct {
	ID   int
	Name string
}

type CarModels struct {
	Models struct {
		Model struct {
			Values []struct {
				Id    int    `json:"id"`
				Name  string `json:"name"`
				Label struct {
					DeDE *string `json:"de_DE"`
					DeAT *string `json:"de_AT"`
					ItIT *string `json:"it_IT"`
					FrBE *string `json:"fr_BE"`
					FrFR *string `json:"fr_FR"`
					FrLU *string `json:"fr_LU"`
					NlBE *string `json:"nl_BE"`
					NlNL *string `json:"nl_NL"`
					EsES *string `json:"es_ES"`
					EnGB *string `json:"en_GB"`
					PlPL *string `json:"pl_PL"`
					SvSE *string `json:"sv_SE"`
					BgBG *string `json:"bg_BG"`
					RoRO *string `json:"ro_RO"`
					CsCZ *string `json:"cs_CZ"`
					HuHU *string `json:"hu_HU"`
					RuRU *string `json:"ru_RU"`
					TrTR *string `json:"tr_TR"`
					UkUA *string `json:"uk_UA"`
					HrHR *string `json:"hr_HR"`
				} `json:"label"`
				MakeId        int    `json:"makeId"`
				ModelLineId   *int   `json:"modelLineId"`
				VehicleTypeId string `json:"vehicleTypeId"`
			} `json:"values"`
		} `json:"model"`
		ModelLine struct {
			Values []struct {
				Id            int    `json:"id"`
				Name          string `json:"name"`
				MakeId        int    `json:"makeId"`
				VehicleTypeId string `json:"vehicleTypeId"`
				Label         struct {
					DeDE *string `json:"de_DE"`
					DeAT *string `json:"de_AT"`
					ItIT *string `json:"it_IT"`
					FrBE *string `json:"fr_BE"`
					FrFR *string `json:"fr_FR"`
					FrLU *string `json:"fr_LU"`
					NlBE *string `json:"nl_BE"`
					NlNL *string `json:"nl_NL"`
					EsES *string `json:"es_ES"`
					EnGB *string `json:"en_GB"`
					PlPL *string `json:"pl_PL"`
					SvSE *string `json:"sv_SE"`
					BgBG *string `json:"bg_BG"`
					RoRO *string `json:"ro_RO"`
					CsCZ *string `json:"cs_CZ"`
					HuHU *string `json:"hu_HU"`
					RuRU *string `json:"ru_RU"`
					TrTR *string `json:"tr_TR"`
					UkUA *string `json:"uk_UA"`
					HrHR *string `json:"hr_HR"`
				} `json:"label"`
			} `json:"values"`
		} `json:"modelLine"`
	} `json:"models"`
}

func main() {
	var cars []CarInfo

	brands := getCarBrands()

	for _, brand := range brands {
		models := getCarModel(brand.ID)
		cars = append(cars, CarInfo{
			Brand:  brand,
			Models: models,
		})
	}

	log.Printf("%+v", cars)
}

func getCarModel(makeID int) CarModels {
	//	https://www.autoscout24.ro/as24-home/api/taxonomy/cars/makes/47/models
	log.Println("Getting model for id", makeID)
	url := fmt.Sprintf("https://www.autoscout24.ro/as24-home/api/taxonomy/cars/makes/%d/models", makeID)
	results, err := http.Get(url)
	defer results.Body.Close()
	if err != nil {
		panic(err)
	}

	var models CarModels
	err = json.NewDecoder(results.Body).Decode(&models)
	if err != nil {
		panic(err)
	}

	return models
}

func getCarBrands() []CarBrand {
	collector := colly.NewCollector()
	var brands []CarBrand

	collector.OnHTML("#make > optgroup:nth-child(2)", func(element *colly.HTMLElement) {
		element.ForEach("option", func(i int, element *colly.HTMLElement) {
			id, err := strconv.Atoi(element.Attr("value"))
			if err != nil {
				panic(err)
			}
			brand := CarBrand{
				ID:   id,
				Name: element.Text,
			}
			brands = append(brands, brand)
		})
	})

	collector.OnHTML("#make > optgroup:nth-child(3)", func(element *colly.HTMLElement) {
		element.ForEach("option", func(i int, element *colly.HTMLElement) {
			id, err := strconv.Atoi(element.Attr("value"))
			if err != nil {
				panic(err)
			}
			brand := CarBrand{
				ID:   id,
				Name: element.Text,
			}
			brands = append(brands, brand)
		})
	})

	collector.Visit("https://www.autoscout24.ro/")

	collector.Wait()

	for _, b := range brands {
		log.Printf("%+v", b)
	}
	return brands
}
