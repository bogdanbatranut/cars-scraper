package mobilede_de

type PageData struct {
	Advertising struct {
		StickyBillboardStatus           int  `json:"stickyBillboardStatus"`
		IsSecondAdvertisingLayerVisible bool `json:"isSecondAdvertisingLayerVisible"`
	} `json:"advertising"`
	Insurance struct {
		EstimatedRates struct {
			Data struct {
			} `json:"data"`
			Status string `json:"status"`
		} `json:"estimatedRates"`
	} `json:"insurance"`
	Parkings struct {
		ParkedAds struct {
			Data struct {
				Items []interface{} `json:"items"`
			} `json:"data"`
			Status string `json:"status"`
		} `json:"parkedAds"`
	} `json:"parkings"`
	Home struct {
		Data   interface{} `json:"data"`
		Status string      `json:"status"`
	} `json:"home"`
	Shared struct {
		Cookies []interface{} `json:"cookies"`
		Filters struct {
			Data struct {
				Dam string   `json:"dam"`
				Fe  []string `json:"fe"`
				Fr  struct {
					Min string `json:"min"`
					Max string `json:"max"`
				} `json:"fr"`
				Ml struct {
					Min string `json:"min"`
					Max string `json:"max"`
				} `json:"ml"`
				Ms []struct {
					Make             string `json:"make"`
					Model            string `json:"model"`
					ModelGroup       string `json:"modelGroup"`
					ModelDescription string `json:"modelDescription"`
				} `json:"ms"`
				Od         string `json:"od"`
				Ref        string `json:"ref"`
				S          string `json:"s"`
				Sb         string `json:"sb"`
				Vc         string `json:"vc"`
				PageNumber string `json:"pageNumber"`
			} `json:"data"`
			Chips struct {
				BasicData []struct {
					FilterKey string `json:"filterKey"`
					Value     struct {
						Make             string `json:"make,omitempty"`
						Model            string `json:"model,omitempty"`
						ModelGroup       string `json:"modelGroup,omitempty"`
						ModelDescription string `json:"modelDescription,omitempty"`
						Min              string `json:"min,omitempty"`
						Max              string `json:"max,omitempty"`
					} `json:"value"`
					Label             string `json:"label"`
					Section           string `json:"section"`
					RemoveFilterValue struct {
						Ms []interface{} `json:"ms,omitempty"`
						Fr interface{}   `json:"fr"`
						Ml interface{}   `json:"ml"`
					} `json:"removeFilterValue"`
					TestId string `json:"testId"`
				} `json:"basicData"`
				TechData []struct {
					FilterKey         string `json:"filterKey"`
					Value             string `json:"value"`
					Label             string `json:"label"`
					Section           string `json:"section"`
					RemoveFilterValue struct {
						Fe struct {
							HYBRIDPLUGIN interface{} `json:"HYBRID_PLUGIN"`
						} `json:"fe"`
					} `json:"removeFilterValue"`
					TestId string `json:"testId"`
				} `json:"techData"`
				OfferDetails []struct {
					FilterKey         string `json:"filterKey"`
					Value             string `json:"value"`
					Label             string `json:"label"`
					Section           string `json:"section"`
					RemoveFilterValue struct {
						Dam interface{} `json:"dam"`
					} `json:"removeFilterValue"`
					TestId string `json:"testId"`
				} `json:"offerDetails"`
			} `json:"chips"`
		} `json:"filters"`
		FiltersConfig struct {
			Config []struct {
				Id     string `json:"id"`
				Type   string `json:"type"`
				Param  string `json:"param"`
				Values []struct {
					Value string `json:"value"`
					Label string `json:"label,omitempty"`
					Param string `json:"param,omitempty"`
				} `json:"values,omitempty"`
				Label       string `json:"label,omitempty"`
				Description struct {
					Title string `json:"title"`
					Text  string `json:"text"`
				} `json:"description,omitempty"`
				Default string `json:"default,omitempty"`
			} `json:"config"`
			VehicleCategory string `json:"vehicleCategory"`
		} `json:"filtersConfig"`
		ReferenceData struct {
			CategoryData struct {
			} `json:"categoryData"`
			ModelData struct {
			} `json:"modelData"`
		} `json:"referenceData"`
		SavedSearches struct {
			Items  []interface{} `json:"items"`
			Status string        `json:"status"`
		} `json:"savedSearches"`
		Theme string `json:"theme"`
		User  struct {
			Data struct {
				UserName       interface{} `json:"userName"`
				IsUserLoggedIn bool        `json:"isUserLoggedIn"`
			} `json:"data"`
			Status string `json:"status"`
		} `json:"user"`
	} `json:"shared"`
	Search struct {
		Srp struct {
			Data struct {
				IsSeoSrp          bool `json:"isSeoSrp"`
				AdvertisingConfig struct {
					AdsenseCsa struct {
						Page struct {
							Ads []struct {
								Id   string `json:"id"`
								Vars struct {
									Hl         string `json:"hl"`
									StyleId    string `json:"styleId"`
									Query      string `json:"query"`
									LinkTarget string `json:"linkTarget"`
									Channel    string `json:"channel"`
									PubId      string `json:"pubId"`
								} `json:"vars"`
								Units []struct {
									Container string `json:"container"`
									Number    int    `json:"number,omitempty"`
								} `json:"units"`
								IsAdblockFallback bool `json:"isAdblockFallback,omitempty"`
							} `json:"ads"`
							SiteId         string `json:"siteId"`
							SingleHeadline string `json:"singleHeadline"`
							MultiHeadline  string `json:"multiHeadline"`
						} `json:"page"`
					} `json:"adsenseCsa"`
					Resources struct {
						Async []string `json:"async"`
					} `json:"resources"`
					Links struct {
						Placements []struct {
							Options string `json:"options"`
							Href    string `json:"href"`
							Id      string `json:"id"`
						} `json:"placements"`
					} `json:"links"`
					Gpt struct {
						Page struct {
							Targeting struct {
								Ch24Mod string   `json:"ch24mod"`
								Pt      string   `json:"pt"`
								Li      string   `json:"li"`
								Ez      string   `json:"ez"`
								Ct      string   `json:"ct"`
								Cm      string   `json:"cm"`
								CmMod   string   `json:"cm_mod"`
								Ch      string   `json:"ch"`
								Ab      string   `json:"ab"`
								Regy    string   `json:"regy"`
								Gear    string   `json:"gear"`
								Ch24Mak string   `json:"ch24mak"`
								Reg     string   `json:"reg"`
								Km      string   `json:"km"`
								Geo     string   `json:"geo"`
								GCmo    string   `json:"g_cmo"`
								Tuev    string   `json:"tuev"`
								Mod     string   `json:"mod"`
								P       string   `json:"p"`
								City    string   `json:"city"`
								Ma      string   `json:"ma"`
								Ch24Reg string   `json:"ch24reg"`
								GHd     string   `json:"g_hd"`
								Restr   string   `json:"restr"`
								GCy     string   `json:"g_cy"`
								KwEcg   string   `json:"kw_ecg"`
								Typ     string   `json:"typ"`
								GPn     string   `json:"g_pn"`
								Hn      string   `json:"hn"`
								Ug      string   `json:"ug"`
								Lang    string   `json:"lang"`
								Kba     []string `json:"kba"`
								GCm     string   `json:"g_cm"`
							} `json:"targeting"`
						} `json:"page"`
						Prebid struct {
							AdUnits []struct {
								Code string `json:"code"`
								Bids []struct {
									Bidder string `json:"bidder"`
									Params struct {
										PlacementId string   `json:"placementId,omitempty"`
										NetworkId   int      `json:"networkId,omitempty"`
										Size        []int    `json:"size,omitempty"`
										SiteId      string   `json:"siteId,omitempty"`
										DelDomain   string   `json:"delDomain,omitempty"`
										Unit        string   `json:"unit,omitempty"`
										Pmzoneid    string   `json:"pmzoneid,omitempty"`
										PublisherId string   `json:"publisherId,omitempty"`
										AdSlot      string   `json:"adSlot,omitempty"`
										AccountId   string   `json:"accountId,omitempty"`
										Keywords    []string `json:"keywords,omitempty"`
										Sizes       []int    `json:"sizes,omitempty"`
										ZoneId      string   `json:"zoneId,omitempty"`
										Position    string   `json:"position,omitempty"`
										Inventory   struct {
											Marke   string `json:"Marke"`
											Modell  string `json:"Modell"`
											Baujahr string `json:"Baujahr"`
											Channel string `json:"Channel"`
										} `json:"inventory,omitempty"`
										AdslotId string `json:"adslotId,omitempty"`
										SupplyId string `json:"supplyId,omitempty"`
									} `json:"params"`
								} `json:"bids"`
								MediaTypes struct {
									Banner struct {
										Sizes [][]int `json:"sizes"`
									} `json:"banner"`
								} `json:"mediaTypes"`
							} `json:"adUnits"`
							PrebidTimeout  int    `json:"prebidTimeout"`
							IdentitiesJson string `json:"identitiesJson"`
						} `json:"prebid"`
						Placements []struct {
							Options   string        `json:"options,omitempty"`
							Id        string        `json:"id"`
							Path      string        `json:"path"`
							Size      []interface{} `json:"size,omitempty"`
							Targeting struct {
								Pos string      `json:"pos"`
								A   interface{} `json:"a"`
							} `json:"targeting"`
							Placeholder       bool              `json:"placeholder,omitempty"`
							SizeMappings      [][][]interface{} `json:"sizeMappings,omitempty"`
							PlaceholderHeight int               `json:"placeholderHeight,omitempty"`
							Oop               bool              `json:"oop,omitempty"`
						} `json:"placements"`
						Aps struct {
							Slots []struct {
								SlotName string  `json:"slotName"`
								Sizes    [][]int `json:"sizes"`
								SlotID   string  `json:"slotID"`
							} `json:"slots"`
						} `json:"aps"`
					} `json:"gpt"`
					Page struct {
						PageType string `json:"pageType"`
					} `json:"page"`
					Adex struct {
						Placements []struct {
							Options        string          `json:"options"`
							Id             string          `json:"id"`
							PushParameters [][]interface{} `json:"pushParameters"`
						} `json:"placements"`
					} `json:"adex"`
					Group struct {
						Ug string `json:"ug"`
					} `json:"group"`
				} `json:"advertisingConfig"`
				MetaData struct {
					Title       string `json:"title"`
					Headline    string `json:"headline"`
					Description string `json:"description"`
					Keywords    string `json:"keywords"`
					PlsUrl      struct {
						Href       string `json:"href"`
						PathEnding string `json:"pathEnding"`
						Path       string `json:"path"`
						Type       string `json:"type"`
					} `json:"plsUrl"`
					CrossChannelSearch string `json:"crossChannelSearch"`
					JsonLdData         struct {
						Graph []struct {
							ItemListElement []struct {
								Type     string `json:"@type"`
								Position int    `json:"position"`
								Item     struct {
									Id   string `json:"@id"`
									Name string `json:"name"`
								} `json:"item"`
							} `json:"itemListElement,omitempty"`
							Type         string `json:"@type"`
							Name         string `json:"name,omitempty"`
							Description  string `json:"description,omitempty"`
							Manufacturer string `json:"manufacturer,omitempty"`
							Model        string `json:"model,omitempty"`
							Offers       struct {
								LowPrice      string `json:"lowPrice"`
								HighPrice     string `json:"highPrice"`
								OfferCount    int    `json:"offerCount"`
								Type          string `json:"@type"`
								PriceCurrency string `json:"priceCurrency"`
							} `json:"offers,omitempty"`
						} `json:"@graph"`
						Context string `json:"@context"`
					} `json:"jsonLdData"`
					Breadcrumbs []struct {
						Label string `json:"label"`
						Href  string `json:"href,omitempty"`
					} `json:"breadcrumbs"`
					Thumbnail struct {
						Src    string `json:"src"`
						Width  int    `json:"width"`
						Height int    `json:"height"`
					} `json:"thumbnail"`
				} `json:"metaData"`
				Aggregations struct {
					Ab []struct {
						Key   string `json:"key"`
						Count int    `json:"count"`
					} `json:"ab"`
					St []struct {
						Key   string `json:"key"`
						Count int    `json:"count"`
					} `json:"st"`
					Con []struct {
						Key   string `json:"key"`
						Count int    `json:"count"`
					} `json:"con"`
					Clim []struct {
						Key   string `json:"key"`
						Count int    `json:"count"`
					} `json:"clim"`
					Ecol []struct {
						Key   string `json:"key"`
						Count int    `json:"count"`
					} `json:"ecol"`
					Fe []struct {
						Key   string `json:"key"`
						Count int    `json:"count"`
					} `json:"fe"`
					Dm []interface{} `json:"dm"`
					Tr []struct {
						Key   string `json:"key"`
						Count int    `json:"count"`
					} `json:"tr"`
					Ft []struct {
						Key   string `json:"key"`
						Count int    `json:"count"`
					} `json:"ft"`
					P []struct {
						Field1 string `json:"25.0,omitempty"`
						Field2 string `json:"1.0,omitempty"`
						Field3 string `json:"95.0,omitempty"`
						Field4 string `json:"50.0,omitempty"`
						Field5 string `json:"5.0,omitempty"`
						Field6 string `json:"99.0,omitempty"`
						Field7 string `json:"75.0,omitempty"`
					} `json:"p"`
					Sr []struct {
						Key   string `json:"key"`
						Count int    `json:"count"`
					} `json:"sr"`
				} `json:"aggregations"`
				SearchResults struct {
					SearchId              string `json:"searchId"`
					NumResultsTotal       int    `json:"numResultsTotal"`
					Page                  int    `json:"page"`
					NumPages              int    `json:"numPages"`
					ObsSearchResultsCount int    `json:"obsSearchResultsCount"`
					HasNextPage           bool   `json:"hasNextPage"`
					Items                 []struct {
						IsEyeCatcher bool `json:"isEyeCatcher,omitempty"`
						NumImages    int  `json:"numImages,omitempty"`
						Attr         struct {
							Cn   string `json:"cn"`
							Z    string `json:"z"`
							Loc  string `json:"loc"`
							Fr   string `json:"fr"`
							Pw   string `json:"pw"`
							Ft   string `json:"ft"`
							Ml   string `json:"ml"`
							Cc   string `json:"cc"`
							Tr   string `json:"tr"`
							Gi   string `json:"gi,omitempty"`
							Ecol string `json:"ecol"`
							Door string `json:"door"`
							Sc   string `json:"sc"`
							C    string `json:"c"`
							Emc  string `json:"emc,omitempty"`
							Pvo  string `json:"pvo,omitempty"`
							Nw   string `json:"nw,omitempty"`
							Eu   string `json:"eu,omitempty"`
							Yc   string `json:"yc,omitempty"`
							Bc   string `json:"bc,omitempty"`
						} `json:"attr,omitempty"`
						ShortTitle         string `json:"shortTitle,omitempty"`
						SubTitle           string `json:"subTitle,omitempty"`
						IsVideoEnabled     bool   `json:"isVideoEnabled,omitempty"`
						BatteryCertificate bool   `json:"batteryCertificate,omitempty"`
						HasElectricEngine  bool   `json:"hasElectricEngine,omitempty"`
						FinancePlans       []struct {
							Type          string `json:"type"`
							Url           string `json:"url"`
							ShowInGallery bool   `json:"showInGallery"`
							Offer         struct {
								BankName              string  `json:"bankName"`
								LoanBroker            string  `json:"loanBroker,omitempty"`
								LoanType              string  `json:"loanType"`
								DownPayment           int     `json:"downPayment"`
								CreditTerm            int     `json:"creditTerm"`
								YearlyMileage         int     `json:"yearlyMileage"`
								CreditAmount          int     `json:"creditAmount"`
								InterestRateNominal   float64 `json:"interestRateNominal"`
								InterestRateEffective float64 `json:"interestRateEffective"`
								MonthlyInstallment    float64 `json:"monthlyInstallment"`
								FinalInstallment      float64 `json:"finalInstallment"`
								TotalInterest         float64 `json:"totalInterest"`
								TotalAmount           float64 `json:"totalAmount"`
								Localized             struct {
									DownPayment              string `json:"downPayment"`
									DownPaymentInPercent     string `json:"downPaymentInPercent"`
									CreditTerm               string `json:"creditTerm"`
									CreditTermInYears        string `json:"creditTermInYears"`
									YearlyMileage            string `json:"yearlyMileage"`
									CreditAmount             string `json:"creditAmount"`
									InterestRateNominal      string `json:"interestRateNominal"`
									InterestRateEffective    string `json:"interestRateEffective"`
									MonthlyInstallment       string `json:"monthlyInstallment"`
									FinalInstallment         string `json:"finalInstallment"`
									TotalInterest            string `json:"totalInterest"`
									TotalAmount              string `json:"totalAmount"`
									Disclaimer               string `json:"disclaimer"`
									MinMonthlyInstallment    string `json:"minMonthlyInstallment,omitempty"`
									MinInterestRateEffective string `json:"minInterestRateEffective,omitempty"`
									MinInterestRateNominal   string `json:"minInterestRateNominal,omitempty"`
									MaxMonthlyInstallment    string `json:"maxMonthlyInstallment,omitempty"`
									MaxInterestRateEffective string `json:"maxInterestRateEffective,omitempty"`
									MaxInterestRateNominal   string `json:"maxInterestRateNominal,omitempty"`
								} `json:"localized"`
								MinMonthlyInstallment    float64 `json:"minMonthlyInstallment,omitempty"`
								MinInterestRateEffective float64 `json:"minInterestRateEffective,omitempty"`
								MinInterestRateNominal   float64 `json:"minInterestRateNominal,omitempty"`
								MaxMonthlyInstallment    float64 `json:"maxMonthlyInstallment,omitempty"`
								MaxInterestRateEffective float64 `json:"maxInterestRateEffective,omitempty"`
								MaxInterestRateNominal   float64 `json:"maxInterestRateNominal,omitempty"`
								BankId                   string  `json:"bankId,omitempty"`
							} `json:"offer"`
							BudgetStatus string `json:"budgetStatus"`
							Fallback     bool   `json:"fallback"`
							DownPayment  int    `json:"downPayment"`
							LoanDuration int    `json:"loanDuration"`
							Localized    struct {
								DownPayment string `json:"downPayment"`
							} `json:"localized"`
							BankId string `json:"bankId,omitempty"`
						} `json:"financePlans,omitempty"`
						SellerId    int `json:"sellerId,omitempty"`
						PriceRating struct {
							Rating             string   `json:"rating"`
							RatingLabel        string   `json:"ratingLabel"`
							ThresholdLabels    []string `json:"thresholdLabels"`
							VehiclePriceOffset int      `json:"vehiclePriceOffset"`
						} `json:"priceRating,omitempty"`
						Segment  string `json:"segment,omitempty"`
						Title    string `json:"title,omitempty"`
						Vc       string `json:"vc,omitempty"`
						Category string `json:"category,omitempty"`
						Id       int    `json:"id,omitempty"`
						Kba      struct {
							Hsn string `json:"hsn"`
							Tsn string `json:"tsn"`
						} `json:"kba,omitempty"`
						CustomDimensions struct {
							Field1 string `json:"10"`
						} `json:"customDimensions,omitempty"`
						Attributes [][]struct {
							Value string `json:"value"`
							Bold  bool   `json:"bold,omitempty"`
						} `json:"attributes,omitempty"`
						Availability interface{} `json:"availability"`
						ContactInfo  struct {
							TypeLocalized string `json:"typeLocalized"`
							Name          string `json:"name"`
							Location      string `json:"location"`
							Rating        *struct {
								Score float64 `json:"score"`
								Count int     `json:"count"`
							} `json:"rating"`
							Logo *struct {
								Src    string `json:"src"`
								SrcSet string `json:"srcSet"`
								Alt    string `json:"alt"`
							} `json:"logo"`
							HasContactPhones bool   `json:"hasContactPhones"`
							ContactPhone     string `json:"contactPhone"`
							Country          string `json:"country"`
							SellerType       string `json:"sellerType"`
						} `json:"contactInfo,omitempty"`
						EmailLink            string      `json:"emailLink,omitempty"`
						IsFinancingAvailable bool        `json:"isFinancingAvailable,omitempty"`
						IsNullLeasingBto     bool        `json:"isNullLeasingBto,omitempty"`
						LeasingRate          interface{} `json:"leasingRate"`
						Make                 string      `json:"make,omitempty"`
						Model                string      `json:"model,omitempty"`
						ObsUrl               string      `json:"obsUrl,omitempty"`
						PreviewImage         struct {
							Src    string `json:"src"`
							SrcSet string `json:"srcSet"`
							Alt    string `json:"alt"`
						} `json:"previewImage,omitempty"`
						Price struct {
							Gross         string  `json:"gross"`
							GrossAmount   float64 `json:"grossAmount"`
							GrossCurrency string  `json:"grossCurrency"`
							Net           string  `json:"net,omitempty"`
							NetAmount     float64 `json:"netAmount,omitempty"`
							Vat           string  `json:"vat,omitempty"`
						} `json:"price,omitempty"`
						RelativeUrl          string      `json:"relativeUrl,omitempty"`
						SupportedLeasingType interface{} `json:"supportedLeasingType"`
						Type                 string      `json:"type"`
						Highlights           []string    `json:"highlights,omitempty"`
						IsNew                bool        `json:"isNew,omitempty"`
						SlotId               string      `json:"slotId,omitempty"`
						SealDetails          struct {
							LocalizedName string `json:"localizedName"`
							InfoUrl       string `json:"infoUrl"`
							Icons         struct {
								FULLICON struct {
									Url    string `json:"url"`
									Width  int    `json:"width"`
									Height int    `json:"height"`
								} `json:"FULL_ICON"`
								SMALLICON struct {
									Url    string `json:"url"`
									Width  int    `json:"width"`
									Height int    `json:"height"`
								} `json:"SMALL_ICON"`
								MEDIUMICON struct {
									Url    string `json:"url"`
									Width  int    `json:"width"`
									Height int    `json:"height"`
								} `json:"MEDIUM_ICON"`
							} `json:"icons"`
							Elements []string `json:"elements"`
						} `json:"sealDetails,omitempty"`
					} `json:"items"`
					Meta struct {
						ObsSrpBaseUrl string `json:"obsSrpBaseUrl"`
					} `json:"meta"`
				} `json:"searchResults"`
				SeoCustomAttributesTest bool `json:"seoCustomAttributesTest"`
			} `json:"data"`
			PageUrl string `json:"pageUrl"`
			Status  string `json:"status"`
		} `json:"srp"`
		SrpSimilarTopAds struct {
			Data   interface{} `json:"data"`
			Status string      `json:"status"`
		} `json:"srp.similarTopAds"`
		SrpSrxAds struct {
			Data   interface{} `json:"data"`
			Status string      `json:"status"`
		} `json:"srp.srxAds"`
		Dsp struct {
			Data            interface{} `json:"data"`
			PageUrl         interface{} `json:"pageUrl"`
			ActiveSectionId interface{} `json:"activeSectionId"`
			PowerUnit       string      `json:"powerUnit"`
			Status          string      `json:"status"`
			ReferenceData   struct {
			} `json:"referenceData"`
			ModelData struct {
			} `json:"modelData"`
			TruckPriceType     string `json:"truckPriceType"`
			ObsBannerIsVisible bool   `json:"obsBannerIsVisible"`
			EditSavedSearch    struct {
				Status string      `json:"status"`
				Error  interface{} `json:"error"`
				Search interface{} `json:"search"`
				Name   string      `json:"name"`
			} `json:"editSavedSearch"`
		} `json:"dsp"`
		Vip struct {
			Ads struct {
			} `json:"ads"`
			Current struct {
			} `json:"current"`
			AdContacts struct {
			} `json:"adContacts"`
			SimilarAds struct {
			} `json:"similarAds"`
			SelectedLeasingRate struct {
			} `json:"selectedLeasingRate"`
			AdListingType                 string `json:"adListingType"`
			LeasingTrackingEmailPlacement string `json:"leasingTrackingEmailPlacement"`
			IsObsTabActive                bool   `json:"isObsTabActive"`
		} `json:"vip"`
	} `json:"search"`
	Ssp struct {
		SavedSearches struct {
			Status       string `json:"status"`
			ItemStatuses struct {
			} `json:"itemStatuses"`
		} `json:"savedSearches"`
		FollowedDealers struct {
			Items        []interface{} `json:"items"`
			Status       string        `json:"status"`
			ItemStatuses struct {
			} `json:"itemStatuses"`
		} `json:"followedDealers"`
	} `json:"ssp"`
}
