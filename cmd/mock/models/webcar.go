package models

type WebCarResponse struct {
	Data []struct {
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
			Webp struct {
				Db84D0187964B84B4Ef46685F0Ba4Ba struct {
					Name             string        `json:"name"`
					FileName         string        `json:"file_name"`
					Uuid             string        `json:"uuid"`
					PreviewUrl       string        `json:"preview_url"`
					OriginalUrl      string        `json:"original_url"`
					Order            int           `json:"order"`
					CustomProperties []interface{} `json:"custom_properties"`
					Extension        string        `json:"extension"`
					Size             int           `json:"size"`
				} `json:"7db84d01-8796-4b84-b4ef-46685f0ba4ba,omitempty"`
				F9661EB5Fe45129A401Ac1059B5044 struct {
					Name             string        `json:"name"`
					FileName         string        `json:"file_name"`
					Uuid             string        `json:"uuid"`
					PreviewUrl       string        `json:"preview_url"`
					OriginalUrl      string        `json:"original_url"`
					Order            int           `json:"order"`
					CustomProperties []interface{} `json:"custom_properties"`
					Extension        string        `json:"extension"`
					Size             int           `json:"size"`
				} `json:"66f9661e-b5fe-4512-9a40-1ac1059b5044,omitempty"`
				Dc389618Fb49CaB91F917Af18Acc17 struct {
					Name             string        `json:"name"`
					FileName         string        `json:"file_name"`
					Uuid             string        `json:"uuid"`
					PreviewUrl       string        `json:"preview_url"`
					OriginalUrl      string        `json:"original_url"`
					Order            int           `json:"order"`
					CustomProperties []interface{} `json:"custom_properties"`
					Extension        string        `json:"extension"`
					Size             int           `json:"size"`
				} `json:"27dc3896-18fb-49ca-b91f-917af18acc17,omitempty"`
				Ea09820A204E5E9C8C8Fa037D4F0B7 struct {
					Name             string        `json:"name"`
					FileName         string        `json:"file_name"`
					Uuid             string        `json:"uuid"`
					PreviewUrl       string        `json:"preview_url"`
					OriginalUrl      string        `json:"original_url"`
					Order            int           `json:"order"`
					CustomProperties []interface{} `json:"custom_properties"`
					Extension        string        `json:"extension"`
					Size             int           `json:"size"`
				} `json:"87ea0982-0a20-4e5e-9c8c-8fa037d4f0b7,omitempty"`
				D0Cdd1Fc0C49Cb87F6200B3840294B struct {
					Name             string        `json:"name"`
					FileName         string        `json:"file_name"`
					Uuid             string        `json:"uuid"`
					PreviewUrl       string        `json:"preview_url"`
					OriginalUrl      string        `json:"original_url"`
					Order            int           `json:"order"`
					CustomProperties []interface{} `json:"custom_properties"`
					Extension        string        `json:"extension"`
					Size             int           `json:"size"`
				} `json:"02d0cdd1-fc0c-49cb-87f6-200b3840294b,omitempty"`
			} `json:"webp"`
			Image struct {
				Ab5799C9D3A410F81E590883E089552 struct {
					Name             string        `json:"name"`
					FileName         string        `json:"file_name"`
					Uuid             string        `json:"uuid"`
					PreviewUrl       string        `json:"preview_url"`
					OriginalUrl      string        `json:"original_url"`
					Order            int           `json:"order"`
					CustomProperties []interface{} `json:"custom_properties"`
					Extension        string        `json:"extension"`
					Size             int           `json:"size"`
				} `json:"5ab5799c-9d3a-410f-81e5-90883e089552,omitempty"`
				D6E8848DFb7B494EA4F2B7549A20Dbf6 struct {
					Name             string        `json:"name"`
					FileName         string        `json:"file_name"`
					Uuid             string        `json:"uuid"`
					PreviewUrl       string        `json:"preview_url"`
					OriginalUrl      string        `json:"original_url"`
					Order            int           `json:"order"`
					CustomProperties []interface{} `json:"custom_properties"`
					Extension        string        `json:"extension"`
					Size             int           `json:"size"`
				} `json:"d6e8848d-fb7b-494e-a4f2-b7549a20dbf6,omitempty"`
				F4Ec85BCd7A43869Edf6C76Cbba33Ac struct {
					Name             string        `json:"name"`
					FileName         string        `json:"file_name"`
					Uuid             string        `json:"uuid"`
					PreviewUrl       string        `json:"preview_url"`
					OriginalUrl      string        `json:"original_url"`
					Order            int           `json:"order"`
					CustomProperties []interface{} `json:"custom_properties"`
					Extension        string        `json:"extension"`
					Size             int           `json:"size"`
				} `json:"3f4ec85b-cd7a-4386-9edf-6c76cbba33ac,omitempty"`
				E5B008928446292AaE3Cdc0617Ac7 struct {
					Name             string        `json:"name"`
					FileName         string        `json:"file_name"`
					Uuid             string        `json:"uuid"`
					PreviewUrl       string        `json:"preview_url"`
					OriginalUrl      string        `json:"original_url"`
					Order            int           `json:"order"`
					CustomProperties []interface{} `json:"custom_properties"`
					Extension        string        `json:"extension"`
					Size             int           `json:"size"`
				} `json:"621e5b00-8928-4462-92aa-e3cdc0617ac7,omitempty"`
				B463A4DcEe5547539C7F3A2B74C695Fa struct {
					Name             string        `json:"name"`
					FileName         string        `json:"file_name"`
					Uuid             string        `json:"uuid"`
					PreviewUrl       string        `json:"preview_url"`
					OriginalUrl      string        `json:"original_url"`
					Order            int           `json:"order"`
					CustomProperties []interface{} `json:"custom_properties"`
					Extension        string        `json:"extension"`
					Size             int           `json:"size"`
				} `json:"b463a4dc-ee55-4753-9c7f-3a2b74c695fa,omitempty"`
			} `json:"image"`
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
	} `json:"data"`
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
