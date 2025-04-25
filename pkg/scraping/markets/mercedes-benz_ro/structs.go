package mercedes_benz_ro

import "time"

type Paging struct {
	PageIndex int `json:"pageIndex"`
	Quantity  int `json:"quantity"`
}

type Searchterm struct {
	FindCompleteTermOnly bool `json:"findCompleteTermOnly"`
}

type Sort struct {
	Field string `json:"field"`
	Order string `json:"order"`
}

type SearchInfo struct {
	Paging     Paging     `json:"paging"`
	Searchterm Searchterm `json:"searchterm"`
	Sort       []Sort     `json:"sort"`
}

type CodesTextEntry struct {
	Codes []string `json:"codes"`
	Text  string   `json:"text"`
}

type FirstRegistration struct {
	Max int `json:"max"`
	Min int `json:"min"`
}

type Criteria struct {
	EngineType        []CodesTextEntry  `json:"engineType"`
	FirstRegistration FirstRegistration `json:"firstRegistrationDate"`
	SalesClass        []CodesTextEntry  `json:"salesClass"`
	BodyGroup         []CodesTextEntry  `json:"bodyGroup"`
}

type Context struct {
	ProcessId string        `json:"processId"`
	Locale    string        `json:"locale"`
	OutletIds []interface{} `json:"outletIds"`
	UiId      string        `json:"uiId"`
}

type VehicleSearchRequest struct {
	SearchInfo SearchInfo `json:"searchInfo"`
	Facets     []string   `json:"facets"`
	Criteria   Criteria   `json:"criteria"`
	Context    Context    `json:"context"`
}

type Search struct {
	VehicleSearchRequest VehicleSearchRequest `json:"vehicleSearchRequest"`
}

type Response struct {
	Version string `json:"version"`
	Token   string `json:"token"`
	Context struct {
		ProcessId string        `json:"processId"`
		Locale    string        `json:"locale"`
		OutletIds []interface{} `json:"outletIds"`
		UiId      string        `json:"uiId"`
	} `json:"context"`
	StatusInfo struct {
		Status struct {
			Code string `json:"code"`
			Text string `json:"text"`
		} `json:"status"`
		Duration int `json:"duration"`
	} `json:"statusInfo"`
	Results []struct {
		Vehicles []struct {
			Id                    string `json:"id"`
			IdReference           string `json:"idReference"`
			TechnicalIdReference  string `json:"technicalIdReference"`
			TransactionId         string `json:"transactionId"`
			CommercialVehicleData struct {
				VehicleGroup struct {
					Code string `json:"code"`
				} `json:"vehicleGroup"`
			} `json:"commercialVehicleData"`
			MetaData struct {
				CountryCode             string    `json:"countryCode"`
				ProcessId               []string  `json:"processId"`
				CountryIsoCode          string    `json:"countryIsoCode"`
				Language                string    `json:"language"`
				VehicleDistributionType string    `json:"vehicleDistributionType"`
				UpdatedDate             time.Time `json:"updatedDate"`
				ImportedDate            time.Time `json:"importedDate"`
				VehicleCategory         struct {
					Code string `json:"code"`
				} `json:"vehicleCategory"`
				StockCategory struct {
					Code string `json:"code"`
				} `json:"stockCategory"`
				LegalStockCategory string    `json:"legalStockCategory"`
				CreatedDate        time.Time `json:"createdDate"`
			} `json:"metaData"`
			PriceInformation struct {
				OfferPrice struct {
					Currency       string  `json:"currency"`
					VatReclaimable bool    `json:"vatReclaimable"`
					BasePrice      float64 `json:"basePrice"`
					TaxInfos       []struct {
						Type struct {
							Code string `json:"code"`
							Text string `json:"text"`
						} `json:"type"`
						TaxAmount float64 `json:"taxAmount"`
						TaxRate   int     `json:"taxRate"`
					} `json:"taxInfos"`
					TotalPrice float64 `json:"totalPrice"`
				} `json:"offerPrice"`
			} `json:"priceInformation"`
			SalesInformation struct {
				SalesProgram            []interface{} `json:"salesProgram"`
				WarrantyExtendedFactory string        `json:"warrantyExtendedFactory"`
				MarketingAttributes     []struct {
					Code string `json:"code"`
					Text string `json:"text"`
				} `json:"marketingAttributes"`
			} `json:"salesInformation"`
			TechnicalData struct {
				Allwheel         bool `json:"allwheel"`
				EmissionStandard struct {
					Code string `json:"code"`
					Text string `json:"text"`
				} `json:"emissionStandard"`
				EnergyEfficiency struct {
					Code string `json:"code"`
					Text string `json:"text"`
				} `json:"energyEfficiency,omitempty"`
				EnergyEfficiencyClass string `json:"energyEfficiencyClass,omitempty"`
				EngineType            struct {
					Code string `json:"code"`
					Text string `json:"text"`
				} `json:"engineType"`
				EnginePrimary struct {
					FuelType struct {
						Code string `json:"code"`
					} `json:"fuelType"`
					PowerKW                int     `json:"powerKW"`
					PowerPS                int     `json:"powerPS"`
					Torque                 int     `json:"torque"`
					CylinderCapacity       int     `json:"cylinderCapacity"`
					CylinderNumber         int     `json:"cylinderNumber"`
					CylinderType           string  `json:"cylinderType"`
					Co2EmissionCombinedMin float64 `json:"co2EmissionCombinedMin,omitempty"`
					Co2EmissionCombinedMax float64 `json:"co2EmissionCombinedMax,omitempty"`
					Co2EmissionCityMax     float64 `json:"co2EmissionCityMax,omitempty"`
					Co2EmissionCityMin     float64 `json:"co2EmissionCityMin,omitempty"`
					Co2EmissionOverlandMax float64 `json:"co2EmissionOverlandMax,omitempty"`
					Co2EmissionOverlandMin float64 `json:"co2EmissionOverlandMin,omitempty"`
					ConsumptionCityMin     float64 `json:"consumptionCityMin"`
					ConsumptionCombinedMin float64 `json:"consumptionCombinedMin"`
					ConsumptionOverlandMin float64 `json:"consumptionOverlandMin"`
					ConsumptionCityMax     float64 `json:"consumptionCityMax"`
					ConsumptionCombinedMax float64 `json:"consumptionCombinedMax"`
					ConsumptionOverlandMax float64 `json:"consumptionOverlandMax"`
					EngineType             struct {
						Code string `json:"code"`
						Text string `json:"text"`
					} `json:"engineType"`
				} `json:"enginePrimary"`
				EngineSecondary struct {
					FuelType struct {
					} `json:"fuelType"`
				} `json:"engineSecondary"`
				FuelTypes []struct {
					Code string `json:"code,omitempty"`
				} `json:"fuelTypes"`
				NumberOfDoors int `json:"numberOfDoors"`
				NumberOfSeats int `json:"numberOfSeats"`
				Steering      struct {
					Code string `json:"code"`
					Text string `json:"text"`
				} `json:"steering"`
				NumberOfGears  int `json:"numberOfGears"`
				TopSpeed       int `json:"topSpeed"`
				WeightKerb     int `json:"weightKerb"`
				WeightMaxTotal int `json:"weightMaxTotal"`
				WeightPayload  int `json:"weightPayload"`
				Wheelbase      int `json:"wheelbase"`
				Items          []struct {
					Id    string `json:"id"`
					Value string `json:"value"`
					Unit  string `json:"unit,omitempty"`
				} `json:"items"`
				TestProcedure string `json:"testProcedure"`
				Emico         struct {
					ApprovedUsage    string    `json:"approvedUsage"`
					DriveConcept     string    `json:"driveConcept"`
					LegalContext     string    `json:"legalContext"`
					TestProcedure    string    `json:"testProcedure"`
					RuleSetVersionID string    `json:"ruleSetVersionID"`
					UpdatedAt        time.Time `json:"updatedAt"`
				} `json:"emico"`
				Acceleration float64 `json:"acceleration"`
				CwValue      float64 `json:"cwValue,omitempty"`
			} `json:"technicalData"`
			VehicleCondition struct {
				Mileage              int  `json:"mileage"`
				NumberPreviousOwners int  `json:"numberPreviousOwners"`
				HasPreviousDamage    bool `json:"hasPreviousDamage"`
				TotalCondition       struct {
					Code string `json:"code"`
					Text string `json:"text"`
				} `json:"totalCondition"`
				DamagesRepaired bool `json:"damagesRepaired"`
			} `json:"vehicleCondition"`
			VehicleConfiguration struct {
				BodyGroup struct {
					Code string `json:"code"`
					Text string `json:"text"`
				} `json:"bodyGroup"`
				CurrentSalesModel struct {
					Code string `json:"code"`
					Text string `json:"text"`
				} `json:"currentSalesModel"`
				LengthGroup struct {
					Code string `json:"code"`
					Text string `json:"text,omitempty"`
				} `json:"lengthGroup"`
				EngineConcept struct {
					Code string `json:"code"`
					Text string `json:"text"`
				} `json:"engineConcept"`
				Equipments []struct {
					CategoryCode           string `json:"categoryCode"`
					CategoryText           string `json:"categoryText,omitempty"`
					CategoryPosition       int    `json:"categoryPosition"`
					Code                   string `json:"code"`
					Text                   string `json:"text,omitempty"`
					Position               int    `json:"position"`
					Type                   string `json:"type"`
					Visibility             string `json:"visibility"`
					EntryType              string `json:"entryType"`
					AvailableInPackageOnly bool   `json:"availableInPackageOnly"`
					CodeText               struct {
						Code string `json:"code"`
						Text string `json:"text,omitempty"`
					} `json:"codeText"`
					GroupCode     string `json:"groupCode,omitempty"`
					GroupCodeText struct {
						Code string `json:"code"`
						Text string `json:"text,omitempty"`
					} `json:"groupCodeText,omitempty"`
					GroupText string `json:"groupText,omitempty"`
				} `json:"equipments"`
				Lines         []interface{} `json:"lines"`
				ModelYear     string        `json:"modelYear"`
				ModelYearCode string        `json:"modelYearCode"`
				PaintGroups   []struct {
					Code string `json:"code"`
					Text string `json:"text"`
				} `json:"paintGroups"`
				PaintPrimary struct {
					Code      string `json:"code"`
					EntryType string `json:"entryType,omitempty"`
					GroupCode string `json:"groupCode"`
					GroupText string `json:"groupText"`
					Text      string `json:"text"`
					Type      string `json:"type"`
				} `json:"paintPrimary"`
				SalesClass struct {
					Code string `json:"code"`
					Text string `json:"text"`
				} `json:"salesClass"`
				SalesDescription string `json:"salesDescription"`
				TransmissionType struct {
					Code string `json:"code"`
					Text string `json:"text"`
				} `json:"transmissionType"`
				Upholstery struct {
					Code       string `json:"code"`
					ColorGroup struct {
						Code string `json:"code"`
					} `json:"colorGroup"`
					GroupCode string `json:"groupCode"`
					Text      string `json:"text"`
					Type      struct {
						Code string `json:"code"`
					} `json:"type"`
				} `json:"upholstery"`
				IsMotorhome  bool   `json:"isMotorhome"`
				DriveConcept string `json:"driveConcept"`
				Alterable    bool   `json:"alterable"`
			} `json:"vehicleConfiguration"`
			VehicleData struct {
				FirstRegistrationDate string `json:"firstRegistrationDate"`
				Baumuster             string `json:"baumuster"`
				Brand                 struct {
					Code string `json:"code"`
					Text string `json:"text"`
				} `json:"brand"`
				Division     string `json:"division"`
				EngineNumber string `json:"engineNumber"`
				IsAmg        bool   `json:"isAmg"`
				OrderNumber  string `json:"orderNumber"`
				TypeKey      struct {
					Hsn string `json:"hsn"`
					Tsn string `json:"tsn"`
					Vsn string `json:"vsn"`
				} `json:"typeKey"`
				UcNumber            string `json:"ucNumber"`
				AllocationPeriod    string `json:"allocationPeriod"`
				FinalInspectionDate string `json:"finalInspectionDate"`
			} `json:"vehicleData"`
			VehicleLocation struct {
				GeoPoint struct {
					Lat float64 `json:"lat"`
					Lon float64 `json:"lon"`
				} `json:"geoPoint"`
				CompanyId string `json:"companyId"`
				OutletId  string `json:"outletId"`
				ZipCode   string `json:"zipCode"`
				Street    string `json:"street"`
				Formatted struct {
					AddressLine1         string `json:"addressLine1"`
					CallbackActivated    bool   `json:"callbackActivated"`
					CityZipcodeLine      string `json:"cityZipcodeLine"`
					ContextSpecificEmail string `json:"contextSpecificEmail"`
					ContextSpecificPhone string `json:"contextSpecificPhone"`
					Email                string `json:"email"`
					Links                struct {
						Website  string `json:"website"`
						Facebook string `json:"facebook,omitempty"`
					} `json:"links"`
					Nameline1          string `json:"nameline1"`
					Nameline2          string `json:"nameline2"`
					PhoneLine          string `json:"phoneLine"`
					ContextSpecificFax string `json:"contextSpecificFax,omitempty"`
					FaxLine            string `json:"faxLine,omitempty"`
				} `json:"formatted"`
				Programs []struct {
					ActivityId     string `json:"activityId"`
					BrandId        string `json:"brandId"`
					ProductGroupId string `json:"productGroupId"`
				} `json:"programs"`
			} `json:"vehicleLocation"`
			VehicleOwner struct {
				GeoPoint struct {
					Lat float64 `json:"lat"`
					Lon float64 `json:"lon"`
				} `json:"geoPoint"`
				CompanyId string `json:"companyId"`
				OutletId  string `json:"outletId"`
				ZipCode   string `json:"zipCode"`
				Street    string `json:"street"`
				Formatted struct {
				} `json:"formatted"`
			} `json:"vehicleOwner"`
			Media         Media `json:"media"`
			GenericFields struct {
				GenericInt0 int `json:"genericInt0"`
			} `json:"genericFields"`
			Status struct {
				Code string `json:"code"`
				Text string `json:"text"`
			} `json:"status"`
		} `json:"vehicles"`
		StatusInfo struct {
			Status struct {
				Code string `json:"code"`
				Text string `json:"text"`
			} `json:"status"`
			Duration int `json:"duration"`
		} `json:"statusInfo"`
		Summary struct {
			PageIndex     int `json:"pageIndex"`
			Quantity      int `json:"quantity"`
			TotalQuantity int `json:"totalQuantity"`
		} `json:"summary"`
		Facets struct {
			SalesClass []struct {
				Text  string `json:"text"`
				Count int    `json:"count"`
				Codes []struct {
					Code     string `json:"code"`
					CodeType string `json:"codeType"`
				} `json:"codes"`
				Img string `json:"img"`
			} `json:"salesClass"`
		} `json:"facets"`
	} `json:"results"`
}

type Media struct {
	Images               []Image `json:"images,omitempty"`
	HasImage             bool    `json:"hasImage"`
	NumberOfImages       int     `json:"numberOfImages,omitempty"`
	NumberOfUniqueImages int     `json:"numberOfUniqueImages,omitempty"`
}

type Image struct {
	Format      string `json:"format"`
	Perspective string `json:"perspective"`
	Source      string `json:"source"`
	Type        string `json:"type"`
	Url         string `json:"url"`
	Visibility  string `json:"visibility"`
	ImageId     int    `json:"imageId"`
	MediaIndex  string `json:"mediaIndex"`
}
