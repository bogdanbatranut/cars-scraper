package webcar

import (
	"carscraper/pkg/jobs"
	"encoding/json"
	"fmt"
	"os"
)

type WebCarStrategy struct {
}

func NewWebCarStrategy() WebCarStrategy {
	return WebCarStrategy{}
}

func (ws WebCarStrategy) Execute(url string) ([]jobs.Ad, error) {
	var ads []jobs.Ad
	webCarResults := readResults(url)

	for _, carData := range webCarResults.Data {
		ad := carData.ToAd()
		ads = append(ads, *ad)
	}

	return ads, nil
	//return nil, nil
}

func readResults(fileNumber string) WebCarResponse {
	//path, err := os.Getwd()
	//if err != nil {
	//	log.Println(err)
	//}
	//fmt.Println(path)

	fileName := fmt.Sprintf("pkg/scraping/markets/webcar_%s.txt", fileNumber)
	data, err := os.ReadFile(fileName)
	check(err)
	//fmt.Print(string(data))

	wcr := WebCarResponse{}

	err = json.Unmarshal(data, &wcr)
	if err != nil {
		panic(err)
	}

	return wcr
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
