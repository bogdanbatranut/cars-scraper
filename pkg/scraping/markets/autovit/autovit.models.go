package autovit

import (
	"carscraper/pkg/jobs"
	"strconv"
)

type AutovitGraphQLResponse struct {
	Data struct {
		AdvertSearch struct {
			Url                     string      `json:"url"`
			SortedBy                interface{} `json:"sortedBy"`
			LocationCriteriaChanged bool        `json:"locationCriteriaChanged"`
			SubscriptionKey         string      `json:"subscriptionKey"`
			TotalCount              int         `json:"totalCount"`
			AppliedLocation         interface{} `json:"appliedLocation"`
			AppliedFilters          []struct {
				Name     string `json:"name"`
				Value    string `json:"value"`
				Typename string `json:"__typename"`
			} `json:"appliedFilters"`
			Breadcrumbs []struct {
				Label    string `json:"label"`
				Url      string `json:"url"`
				Typename string `json:"__typename"`
			} `json:"breadcrumbs"`
			PageInfo struct {
				PageSize      int    `json:"pageSize"`
				CurrentOffset int    `json:"currentOffset"`
				Typename      string `json:"__typename"`
			} `json:"pageInfo"`
			Facets []struct {
				Options []struct {
					Label    string `json:"label"`
					Url      string `json:"url"`
					Count    int    `json:"count"`
					Typename string `json:"__typename"`
				} `json:"options"`
				Typename string `json:"__typename"`
			} `json:"facets"`
			AlternativeLinks []struct {
				Name  string `json:"name"`
				Title string `json:"title"`
				Links []struct {
					Title    string `json:"title"`
					Url      string `json:"url"`
					Counter  int    `json:"counter"`
					Typename string `json:"__typename"`
				} `json:"links"`
				Typename string `json:"__typename"`
			} `json:"alternativeLinks"`
			LatestAdId string `json:"latestAdId"`
			Edges      []struct {
				Vas struct {
					IsHighlighted bool    `json:"isHighlighted"`
					IsPromoted    bool    `json:"isPromoted"`
					BumpDate      *string `json:"bumpDate"`
					Typename      string  `json:"__typename"`
				} `json:"vas"`
				Node     CarNode `json:"node"`
				Typename string  `json:"__typename"`
			} `json:"edges"`
			Typename    string `json:"__typename"`
			SortOptions []struct {
				SearchKey string `json:"searchKey"`
				Label     string `json:"label"`
				Typename  string `json:"__typename"`
			} `json:"sortOptions"`
		} `json:"advertSearch"`
		Typename string `json:"__typename"`
	} `json:"data"`
}

type CarNode struct {
	Id               string        `json:"id"`
	Title            string        `json:"title"`
	CreatedAt        string        `json:"createdAt"`
	ShortDescription string        `json:"shortDescription"`
	Url              string        `json:"url"`
	Badges           []interface{} `json:"badges"`
	Category         struct {
		Id       string `json:"id"`
		Typename string `json:"__typename"`
	} `json:"category"`
	Location struct {
		City struct {
			Name     string `json:"name"`
			Typename string `json:"__typename"`
		} `json:"city"`
		Region struct {
			Name     string `json:"name"`
			Typename string `json:"__typename"`
		} `json:"region"`
		Typename string `json:"__typename"`
	} `json:"location"`
	Thumbnail struct {
		X1       string `json:"x1"`
		X2       string `json:"x2"`
		Typename string `json:"__typename"`
	} `json:"thumbnail"`
	Price struct {
		Amount struct {
			Units        int    `json:"units"`
			Nanos        int    `json:"nanos"`
			Value        string `json:"value"`
			CurrencyCode string `json:"currencyCode"`
			Typename     string `json:"__typename"`
		} `json:"amount"`
		Badges     []string `json:"badges"`
		GrossPrice *struct {
			Value        string `json:"value"`
			CurrencyCode string `json:"currencyCode"`
			Typename     string `json:"__typename"`
		} `json:"grossPrice"`
		NetPrice *struct {
			Value        string `json:"value"`
			CurrencyCode string `json:"currencyCode"`
			Typename     string `json:"__typename"`
		} `json:"netPrice"`
		Typename string `json:"__typename"`
	} `json:"price"`
	Parameters []Parameter `json:"parameters"`
	SellerLink struct {
		Id         string  `json:"id"`
		Name       *string `json:"name"`
		WebsiteUrl *string `json:"websiteUrl"`
		Logo       *struct {
			X1       string `json:"x1"`
			Typename string `json:"__typename"`
		} `json:"logo"`
		Typename string `json:"__typename"`
	} `json:"sellerLink"`
	BrandProgram struct {
		Logo      interface{} `json:"logo"`
		SearchUrl interface{} `json:"searchUrl"`
		Name      interface{} `json:"name"`
		Typename  string      `json:"__typename"`
	} `json:"brandProgram"`
	Dealer4ThPackage *struct {
		Package struct {
			Id       string `json:"id"`
			Name     string `json:"name"`
			Typename string `json:"__typename"`
		} `json:"package"`
		Services []struct {
			Code     string `json:"code"`
			Label    string `json:"label"`
			Typename string `json:"__typename"`
		} `json:"services"`
		Photos struct {
			Nodes []struct {
				Url      string `json:"url"`
				Typename string `json:"__typename"`
			} `json:"nodes"`
			TotalCount int    `json:"totalCount"`
			Typename   string `json:"__typename"`
		} `json:"photos"`
		Typename string `json:"__typename"`
	} `json:"dealer4thPackage"`
	PriceEvaluation struct {
		Indicator string `json:"indicator"`
		Typename  string `json:"__typename"`
	} `json:"priceEvaluation"`
	Typename string `json:"__typename"`
}

type Parameter struct {
	Key          string `json:"key"`
	DisplayValue string `json:"displayValue"`
	Label        string `json:"label"`
	Value        string `json:"value"`
	Typename     string `json:"__typename"`
}

func (carnode CarNode) getSellerName() *string {
	return carnode.SellerLink.Name
}

func (carnode CarNode) getSellerLink() *string {
	return carnode.SellerLink.WebsiteUrl
}

func (carnode CarNode) getAutovitID() *int {
	idStr := carnode.Id
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return nil
	}
	return &id
}

func (cn CarNode) ToAd() *jobs.Ad {

	carAd := jobs.Ad{
		Brand:           *getMake(cn.Parameters),
		Model:           *getModel(cn.Parameters),
		Year:            *getYear(cn.Parameters),
		Km:              *getKm(cn.Parameters),
		Fuel:            *getFuelType(cn.Parameters),
		Price:           cn.Price.Amount.Units,
		Ad_url:          cn.Url,
		SellerType:      "",
		SellerName:      cn.SellerLink.Name,
		SellerMarketURL: cn.SellerLink.WebsiteUrl,
	}

	if cn.SellerLink.Name == nil {
		pf := "privat"
		carAd.SellerNameInMarket = &pf
		carAd.SellerOwnURL = &pf
		carAd.SellerMarketURL = &pf
		carAd.SellerName = &pf
		carAd.SellerType = pf
	}

	return &carAd
}

func getMake(params []Parameter) *string {
	for _, param := range params {
		if param.Key == "make" {
			return &param.DisplayValue
		}
	}
	return nil
}

func getModel(params []Parameter) *string {
	for _, param := range params {
		if param.Key == "model" {
			if param.DisplayValue == "GLE" {
				s := "gle_classe"
				return &s
			}
			return &param.DisplayValue
		}
	}
	return nil
}

func getKm(params []Parameter) *int {
	for _, param := range params {
		if param.Key == "mileage" {
			val, err := strconv.Atoi(param.Value)
			if err != nil {
				panic(err)
			}
			return &val
		}
	}
	return nil
}

func getYear(params []Parameter) *int {
	for _, param := range params {
		if param.Key == "year" {
			val, err := strconv.Atoi(param.Value)
			if err != nil {
				panic(err)
			}
			return &val
		}
	}
	return nil
}

func getFuelType(params []Parameter) *string {
	for _, param := range params {
		if param.Key == "fuel_type" {
			return &param.Value
		}
	}
	return nil
}
