package main

import (
	"carscraper/pkg/scraping/markets/webcar"
	"log"
	"strconv"
	"strings"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {

	wcs := webcar.NewWebCarStrategy()
	wcs.Execute("1")

	log.Println(getMileage("8.807 km"))

	//data, err := os.ReadFile("cmd/mock/webcar_1.txt")
	//check(err)
	////fmt.Print(string(data))
	//
	//wcr := models.WebCarResponse{}
	//
	//err = json.Unmarshal(data, &wcr)
	//if err != nil {
	//	panic(err)
	//}
	//
	//fmt.Printf("%v", wcr)
}

func getMileage(str string) int {
	// ": "8.807 km",
	mileageStr := str
	mileageStr = strings.Trim(mileageStr, " km")
	mileageStr = strings.Replace(mileageStr, ".", "", -1)

	m, err := strconv.Atoi(mileageStr)
	if err != nil {
		return -1
	}
	return m
}
