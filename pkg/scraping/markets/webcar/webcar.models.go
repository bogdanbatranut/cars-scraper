package webcar

import (
	"carscraper/pkg/jobs"
	"fmt"
	"math"
	"strconv"
	"strings"
)

type WebCarResponse struct {
	Data []CarData `json:"data"`
	Meta struct {
		Count  int `json:"count"`
		Filter struct {
			Fuel    []string `json:"fuel"`
			Gearbox []string `json:"gearbox"`
			Make    struct {
				Field1 string `json:"1"`
			} `json:"make"`
			Model struct {
				Field1 string `json:"1"`
			} `json:"model"`
			RegisteredOn struct {
				From string `json:"from"`
			} `json:"registered_on"`
			Mileage struct {
				To string `json:"to"`
			} `json:"mileage"`
		} `json:"filter"`
		NextUrl string `json:"next_url"`
		Page    int    `json:"page"`
	} `json:"meta"`
}

type CarData struct {
	Description         string   `json:"description"`
	DisplayVisitedBadge bool     `json:"displayVisitedBadge"`
	Fuel                string   `json:"fuel"`
	FuelIcon            string   `json:"fuelIcon"`
	GrossPrice          *float64 `json:"grossPrice"`
	Id                  int      `json:"id"`
	IsEligible          bool     `json:"isEligible"`
	IsDeleted           bool     `json:"isDeleted"`
	IsVatReclaimable    bool     `json:"isVatReclaimable"`
	IsVisited           bool     `json:"isVisited"`
	LeasingValue        *float64 `json:"leasingValue"`
	Meta                struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		Link        string `json:"link"`
	} `json:"meta"`
	Mileage           string  `json:"mileage"`
	Price             float64 `json:"price"`
	DisplayPriceValue float64 `json:"displayPriceValue"`
	RegisteredOn      string  `json:"registeredOn"`
	Thumbnails        struct {
		Small    string `json:"small"`
		Medium   string `json:"medium"`
		Large    string `json:"large"`
		Original struct {
			Id        int    `json:"id"`
			VehicleId int    `json:"vehicle_id"`
			Large     string `json:"large"`
			Medium    string `json:"medium"`
			Small     string `json:"small"`
			Original  string `json:"original"`
			CreatedAt string `json:"created_at"`
			UpdatedAt string `json:"updated_at"`
			SortOrder int    `json:"sort_order"`
		} `json:"original"`
	} `json:"thumbnails"`
	Title                    string `json:"title"`
	UniqueId                 string `json:"uniqueId"`
	UserCanDelete            bool   `json:"userCanDelete"`
	UserCanEdit              bool   `json:"userCanEdit"`
	UserCanEditOwn           bool   `json:"userCanEditOwn"`
	UserCanViewHiddenDetails bool   `json:"userCanViewHiddenDetails"`
	Vendor                   string `json:"vendor"`
	Salesman                 struct {
		Name  string `json:"name"`
		Title string `json:"title"`
		Phone struct {
			Main       string `json:"main"`
			Additional string `json:"additional"`
		} `json:"phone"`
	} `json:"salesman"`
	SellerDetails struct {
		CompanyName   interface{} `json:"companyName"`
		ContactName   interface{} `json:"contactName"`
		ProfileName   interface{} `json:"profileName"`
		IsAgreed      string      `json:"isAgreed"`
		IsIndependent string      `json:"isIndependent"`
		DealerType    struct {
			Value string `json:"value"`
			Label string `json:"label"`
			Badge string `json:"badge"`
		} `json:"dealerType"`
	} `json:"sellerDetails"`
}

func (c CarData) ToAd() *jobs.Ad {
	carAd := jobs.Ad{
		Brand:              "",
		Model:              "",
		Year:               c.getYear(),
		Km:                 c.getMileage(),
		Fuel:               c.getFuel(),
		Price:              c.getPrice(),
		AdID:               c.UniqueId,
		Ad_url:             "",
		SellerType:         "",
		SellerName:         &c.Vendor,
		SellerNameInMarket: &c.Vendor,
		SellerOwnURL:       c.getSellerOwnURL(),
		SellerMarketURL:    c.getSellerMarketURL(),
	}
	return &carAd
}

func (c CarData) getYear() int {
	yStr := c.RegisteredOn[len(c.RegisteredOn)-4:]
	y, err := strconv.Atoi(yStr)
	if err != nil {
		return -1
	}
	return y
}

func (c CarData) getMileage() int {
	// ": "8.807 km",
	mileageStr := c.Mileage
	mileageStr = strings.Trim(mileageStr, " km")
	mileageStr = strings.Replace(mileageStr, ".", "", -1)

	m, err := strconv.Atoi(mileageStr)
	if err != nil {
		return -1
	}
	return m
}

func (c CarData) getFuel() string {

	return strings.ToLower(c.Fuel)
}

func (c CarData) getPrice() int {
	return int(math.Round(c.Price))
}

func (c CarData) getSellerOwnURL() *string {
	u := fmt.Sprintf("www.%s.com", c.Vendor)
	return &u
}

func (c CarData) getSellerMarketURL() *string {
	u := fmt.Sprintf("www.market.%s.com", c.Vendor)
	return &u
}
