package olx

import (
	"fmt"
	"math"
	"strconv"
	"time"
)

type AdParam struct {
	Key   string     `json:"key"`
	Name  string     `json:"name"`
	Type  string     `json:"type"`
	Value ParamValue `json:"value"`
}

type ParamValue struct {
	Value                  float64     `json:"value,omitempty"`
	Type                   string      `json:"type,omitempty"`
	Arranged               bool        `json:"arranged,omitempty"`
	Budget                 bool        `json:"budget,omitempty"`
	Currency               string      `json:"currency,omitempty"`
	Negotiable             bool        `json:"negotiable,omitempty"`
	ConvertedValue         interface{} `json:"converted_value"`
	PreviousValue          float64     `json:"previous_value,omitempty"`
	ConvertedPreviousValue interface{} `json:"converted_previous_value"`
	ConvertedCurrency      interface{} `json:"converted_currency"`
	Label                  string      `json:"label"`
	PreviousLabel          *string     `json:"previous_label,omitempty"`
	Key                    string      `json:"key,omitempty"`
}

type OlxAd struct {
	Id                int        `json:"id"`
	Url               string     `json:"url"`
	Title             string     `json:"title"`
	LastRefreshTime   time.Time  `json:"last_refresh_time"`
	CreatedTime       time.Time  `json:"created_time"`
	ValidToTime       time.Time  `json:"valid_to_time"`
	PushupTime        *time.Time `json:"pushup_time"`
	OmnibusPushupTime time.Time  `json:"omnibus_pushup_time,omitempty"`
	Description       string     `json:"description"`
	Promotion         struct {
		Highlighted   bool     `json:"highlighted"`
		Urgent        bool     `json:"urgent"`
		TopAd         bool     `json:"top_ad"`
		Options       []string `json:"options"`
		B2CAdPage     bool     `json:"b2c_ad_page"`
		PremiumAdPage bool     `json:"premium_ad_page"`
	} `json:"promotion"`
	Params    []AdParam `json:"params"`
	KeyParams []string  `json:"key_params"`
	Business  bool      `json:"business"`
	User      struct {
		Id                       int         `json:"id"`
		Created                  time.Time   `json:"created"`
		OtherAdsEnabled          bool        `json:"other_ads_enabled"`
		Name                     string      `json:"name"`
		Logo                     interface{} `json:"logo"`
		LogoAdPage               interface{} `json:"logo_ad_page"`
		SocialNetworkAccountType *string     `json:"social_network_account_type"`
		Photo                    *string     `json:"photo"`
		BannerMobile             string      `json:"banner_mobile"`
		BannerDesktop            string      `json:"banner_desktop"`
		CompanyName              string      `json:"company_name"`
		About                    string      `json:"about"`
		B2CBusinessPage          bool        `json:"b2c_business_page"`
		IsOnline                 bool        `json:"is_online"`
		LastSeen                 time.Time   `json:"last_seen"`
		SellerType               interface{} `json:"seller_type"`
		Uuid                     *string     `json:"uuid"`
	} `json:"user"`
	Status  string `json:"status"`
	Contact struct {
		Name        string `json:"name"`
		Phone       bool   `json:"phone"`
		Chat        bool   `json:"chat"`
		Negotiation bool   `json:"negotiation"`
		Courier     bool   `json:"courier"`
	} `json:"contact"`
	Map struct {
		Zoom         int     `json:"zoom"`
		Lat          float64 `json:"lat"`
		Lon          float64 `json:"lon"`
		Radius       int     `json:"radius"`
		ShowDetailed bool    `json:"show_detailed"`
	} `json:"map"`
	Location struct {
		City struct {
			Id             int    `json:"id"`
			Name           string `json:"name"`
			NormalizedName string `json:"normalized_name"`
		} `json:"city"`
		Region struct {
			Id             int    `json:"id"`
			Name           string `json:"name"`
			NormalizedName string `json:"normalized_name"`
		} `json:"region"`
		District struct {
			Id   int    `json:"id"`
			Name string `json:"name"`
		} `json:"district,omitempty"`
	} `json:"location"`
	Photos []struct {
		Id       int    `json:"id"`
		Filename string `json:"filename"`
		Rotation int    `json:"rotation"`
		Width    int    `json:"width"`
		Height   int    `json:"height"`
		Link     string `json:"link"`
	} `json:"photos"`
	Partner *struct {
		Code string `json:"code"`
	} `json:"partner"`
	Category struct {
		Id   int    `json:"id"`
		Type string `json:"type"`
	} `json:"category"`
	Delivery struct {
		Rock struct {
			OfferId interface{} `json:"offer_id"`
			Active  bool        `json:"active"`
			Mode    interface{} `json:"mode"`
		} `json:"rock"`
	} `json:"delivery"`
	Safedeal struct {
		Weight          int           `json:"weight"`
		WeightGrams     int           `json:"weight_grams"`
		Status          string        `json:"status"`
		SafedealBlocked bool          `json:"safedeal_blocked"`
		AllowedQuantity []interface{} `json:"allowed_quantity"`
	} `json:"safedeal"`
	Shop struct {
		Subdomain interface{} `json:"subdomain"`
	} `json:"shop"`
	OfferType   string `json:"offer_type"`
	ExternalUrl string `json:"external_url,omitempty"`
}

type Href struct {
	Href string `json:"href"`
}

type OlxResponse struct {
	Data     []OlxAd `json:"data"`
	Metadata struct {
		TotalElements     int    `json:"total_elements"`
		VisibleTotalCount int    `json:"visible_total_count"`
		Promoted          []int  `json:"promoted"`
		SearchId          string `json:"search_id"`
		Adverts           struct {
			Places interface{} `json:"places"`
			Config struct {
				Targeting struct {
					Env               string        `json:"env"`
					Lang              string        `json:"lang"`
					Account           string        `json:"account"`
					DfpUserId         string        `json:"dfp_user_id"`
					UserStatus        string        `json:"user_status"`
					CatL0Id           string        `json:"cat_l0_id"`
					CatL1Id           string        `json:"cat_l1_id"`
					CatL2Id           string        `json:"cat_l2_id"`
					CatL0             string        `json:"cat_l0"`
					CatL0Path         string        `json:"cat_l0_path"`
					CatL1             string        `json:"cat_l1"`
					CatL1Path         string        `json:"cat_l1_path"`
					CatL2             string        `json:"cat_l2"`
					CatL2Path         string        `json:"cat_l2_path"`
					CatL0Name         string        `json:"cat_l0_name"`
					CatL1Name         string        `json:"cat_l1_name"`
					CatL2Name         string        `json:"cat_l2_name"`
					CatId             string        `json:"cat_id"`
					PrivateBusiness   string        `json:"private_business"`
					OfferSeek         string        `json:"offer_seek"`
					View              string        `json:"view"`
					Currency          string        `json:"currency"`
					SearchEngineInput string        `json:"search_engine_input"`
					Page              string        `json:"page"`
					Segment           []interface{} `json:"segment"`
					AppVersion        string        `json:"app_version"`
					Model             []string      `json:"model"`
					Petrol            []string      `json:"petrol"`
					RulajPana         []string      `json:"rulaj_pana"`
					Year              []string      `json:"year"`
				} `json:"targeting"`
			} `json:"config"`
		} `json:"adverts"`
		Source struct {
			Promoted []int `json:"promoted"`
			Organic  []int `json:"organic"`
		} `json:"source"`
	} `json:"metadata"`
	Links struct {
		Self  *Href `json:"self,omitempty"`
		Next  *Href `json:"next,omitempty"`
		First *Href `json:"first,omitempty"`
	} `json:"links"`
}

func (oAd OlxAd) getPriceParam() *AdParam {
	for _, param := range oAd.Params {
		if param.Key == "price" {
			return &param
		}
	}
	return nil
}

func (oAd OlxAd) getPrice() int {
	for _, param := range oAd.Params {
		if param.Key == "price" {
			return int(math.Round(param.Value.Value))
		}
	}
	return 0
}

func (oAd OlxAd) getYear() int {
	for _, param := range oAd.Params {
		if param.Key == "year" {
			year, err := strconv.Atoi(param.Value.Key)
			if err != nil {
				return 0
			}
			return year
		}
	}
	return 0
}

func (oAd OlxAd) getKm() int {
	for _, param := range oAd.Params {
		if param.Key == "rulaj_pana" {
			km, err := strconv.Atoi(param.Value.Key)
			if err != nil {
				return 0
			}
			return km
		}
	}
	return 0
}

func (oAd OlxAd) getSellerType() string {
	if oAd.Business {
		return "dealer"
	}
	return "privat"
}

func (oAd OlxAd) getThumbnailURL() *string {
	if len(oAd.Photos) == 0 {
		return nil
	}
	w := oAd.Photos[0].Width
	h := oAd.Photos[0].Height
	fileName := oAd.Photos[0].Filename
	th := fmt.Sprintf("https://frankfurt.apollo.olxcdn.com:443/v1/files/%s/image;s=%dx%d", fileName, w, h)
	return &th
}

func (oAd OlxAd) getModel() string {
	for _, param := range oAd.Params {
		if param.Key == "model" {
			return param.Value.Key
		}
	}
	return ""
}

func (oAd OlxAd) getBrand() int {
	return oAd.Category.Id
}
