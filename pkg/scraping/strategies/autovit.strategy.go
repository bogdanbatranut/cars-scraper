package strategies

import "time"

type AutovitStrategy struct {
}

func NewAutovitStrategy() AutovitStrategy {
	return AutovitStrategy{}
}

func (as AutovitStrategy) Execute(url string) ([]Ad, error) {
	var ads []Ad = []Ad{
		Ad{
			Brand:      "test",
			Model:      "test",
			Year:       1,
			Km:         1,
			Fuel:       "diesel",
			Price:      1000,
			AdID:       nil,
			Ad_url:     "www",
			SellerType: "dealer",
			SellerName: nil,
			SellerURL:  nil,
		},
		Ad{
			Brand:      "ttt",
			Model:      "eee",
			Year:       0,
			Km:         0,
			Fuel:       "petro",
			Price:      0,
			AdID:       nil,
			Ad_url:     "",
			SellerType: "",
			SellerName: nil,
			SellerURL:  nil,
		},
	}
	time.Sleep(2 * time.Second)
	return ads, nil
}
