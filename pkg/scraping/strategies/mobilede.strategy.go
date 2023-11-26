package strategies

import "time"

type MobileDeStrategy struct {
}

func NewMobileDeStrategy() MobileDeStrategy {
	return MobileDeStrategy{}
}

func (as MobileDeStrategy) Execute(url string) ([]Ad, error) {
	sellserName := "MOBILE>DE"
	var ads []Ad = []Ad{
		Ad{
			Brand:      "mibiletest",
			Model:      "MobileDeStrategytest",
			Year:       1,
			Km:         1,
			Fuel:       "diesel",
			Price:      1000,
			AdID:       nil,
			Ad_url:     "www.MobileDeStrategy",
			SellerType: "dealer",
			SellerName: &sellserName,
			SellerURL:  nil,
		},
		Ad{
			Brand:      "MobileDeStrategyttt",
			Model:      "MobileDeStrategyeee",
			Year:       0,
			Km:         0,
			Fuel:       "petrol",
			Price:      0,
			AdID:       nil,
			Ad_url:     "MobileDeStrategy",
			SellerType: "",
			SellerName: &sellserName,
			SellerURL:  nil,
		},
	}
	time.Sleep(3 * time.Second)
	return ads, nil
}
