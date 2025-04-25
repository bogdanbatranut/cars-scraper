package oferte_bmw

import "time"

type OferteBMWResponse struct {
	List []struct {
		Id        int    `json:"id"`
		Vin17     string `json:"vin17"`
		Type      string `json:"type"`
		IsNew     bool   `json:"isNew"`
		ModelCode string `json:"modelCode"`
		Title     string `json:"title"`
		Brand     struct {
			Id    int    `json:"id"`
			Label string `json:"label"`
		} `json:"brand"`
		Series struct {
			Id    int    `json:"id"`
			Label string `json:"label"`
		} `json:"series"`
		BodyType struct {
			Id    int    `json:"id"`
			Label string `json:"label"`
		} `json:"bodyType"`
		SeriesCode     string    `json:"seriesCode"`
		ProductionYear int       `json:"productionYear"`
		Registration   time.Time `json:"registration"`
		Mileage        int       `json:"mileage"`
		Fuel           struct {
			Id    int    `json:"id"`
			Label string `json:"label"`
		} `json:"fuel"`
		ConsumptionFuel  float64 `json:"consumptionFuel"`
		Emission         int     `json:"emission"`
		EmissionStandard struct {
			Id    int    `json:"id"`
			Label string `json:"label"`
		} `json:"emissionStandard"`
		EmissionMeasurementStandard string `json:"emissionMeasurementStandard"`
		PowerKW                     int    `json:"powerKW"`
		PowerHP                     int    `json:"powerHP"`
		Capacity                    int    `json:"capacity"`
		Transmission                struct {
			Id    int    `json:"id"`
			Label string `json:"label"`
		} `json:"transmission"`
		Color struct {
			LabelEN string `json:"labelEN"`
			Code    string `json:"code"`
			Id      int    `json:"id"`
			Label   string `json:"label"`
		} `json:"color"`
		Images             int       `json:"images"`
		ImagesLastChanged  time.Time `json:"imagesLastChanged"`
		Panoramas          int       `json:"panoramas"`
		Exterior360        bool      `json:"exterior360"`
		Interior360        bool      `json:"interior360"`
		Warranty           int       `json:"warranty"`
		UsedBrand          bool      `json:"usedBrand"`
		Reservable         bool      `json:"reservable"`
		ReservationStatus  string    `json:"reservationStatus"`
		VatReclaimable     bool      `json:"vatReclaimable"`
		Leasable           bool      `json:"leasable"`
		TransactionalPrice int       `json:"transactionalPrice"`
		Dealer             struct {
			Id             int     `json:"id"`
			BunoBMW        string  `json:"bunoBMW"`
			BunoMINI       string  `json:"bunoMINI,omitempty"`
			Owner          string  `json:"owner"`
			OwnerName      string  `json:"ownerName"`
			Name           string  `json:"name"`
			LegalName      string  `json:"legalName"`
			Street         string  `json:"street"`
			Zip            string  `json:"zip"`
			City           string  `json:"city"`
			Lat            float64 `json:"lat"`
			Lng            float64 `json:"lng"`
			BunoMOTORCYCLE string  `json:"bunoMOTORCYCLE,omitempty"`
		} `json:"dealer"`
		Created          time.Time `json:"created"`
		Age              int       `json:"age"`
		NewPrice         int       `json:"newPrice"`
		OptionsPrice     int       `json:"optionsPrice"`
		AccessoriesPrice int       `json:"accessoriesPrice"`
		IsYUC            bool      `json:"isYUC"`
		Reserved         bool      `json:"reserved"`
		Extended         struct {
			Brand string `json:"brand"`
			Buno  string `json:"buno"`
		} `json:"extended"`
		TransactionalPriceUpdated time.Time `json:"transactionalPriceUpdated,omitempty"`
		Video                     string    `json:"video,omitempty"`
		PanoramasLastChanged      time.Time `json:"panoramasLastChanged,omitempty"`
	} `json:"$list"`
	Count struct {
		Total int `json:"$total"`
	} `json:"$count"`
}

type RequestBody struct {
	Match Match  `json:"$match"`
	Skip  int    `json:"$skip"`
	Limit int    `json:"$limit"`
	Sort  []Sort `json:"$sort"`
}

type TransactionalPrice struct {
	Min int `json:"$min"`
	Max int `json:"$max"`
}

type Registration struct {
	Min string `json:"$min"`
}

type Mileage struct {
	//Min int `json:"$min"`
	Max int `json:"$max"`
}

type Sort struct {
	Field string `json:"$field"`
	Order string `json:"$order"`
}

type Match struct {
	TransactionalPrice TransactionalPrice `json:"transactionalPrice"`
	Brand              int                `json:"brand"`
	Series             int                `json:"series"`
	Variant            int                `json:"variant"`
	Registration       Registration       `json:"registration"`
	Mileage            Mileage            `json:"mileage"`
	Fuel               int                `json:"fuel"`
}
