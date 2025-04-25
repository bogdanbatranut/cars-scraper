package autoklass

import "time"

type Response struct {
	ErrorCode struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"errorCode"`
	Response []struct {
		Id                    int    `json:"id"`
		IdMaster              int    `json:"idMaster"`
		Title                 string `json:"title"`
		Slug                  string `json:"slug"`
		CarBrandID            int    `json:"carBrandID"`
		CarModelID            int    `json:"carModelID"`
		CarSubModel           string `json:"carSubModel"`
		Vin                   string `json:"vin"`
		StandardPrice         int    `json:"standardPrice"`
		SalePrice             int    `json:"salePrice"`
		Discount              int    `json:"discount"`
		BranchID              int    `json:"branchID"`
		ConsultantID          int    `json:"consultantID"`
		Status                string `json:"status"`
		Fuel                  string `json:"fuel"`
		Pollution             string `json:"pollution"`
		Transmission          string `json:"transmission"`
		Traction              string `json:"traction"`
		Km                    int    `json:"km"`
		Cylinder              int    `json:"cylinder"`
		Power                 int    `json:"power"`
		Consumption           int    `json:"consumption"`
		EmissionsFrom         int    `json:"emissionsFrom"`
		EmissionsTo           int    `json:"emissionsTo"`
		Emissions             int    `json:"emissions"`
		DateManufacture       int64  `json:"dateManufacture"`
		DateFirstRegistration int    `json:"dateFirstRegistration"`
		CarBodyID             int    `json:"carBodyID"`
		Color                 string `json:"color"`
		Availability          int    `json:"availability"`
		MercedesBenzCertified int    `json:"mercedesBenzCertified"`
		IsOffer               int    `json:"isOffer"`
		Remat                 int    `json:"remat"`
		Vat                   int    `json:"vat"`
		IsActive              int    `json:"isActive"`
		Stock                 int    `json:"stock"`
		Description           string `json:"description"`
		DescriptionRtf        string `json:"description_rtf"`
		DescriptionAutovitRtf string `json:"description_autovit_rtf"`
		TransmissionDetails   string `json:"transmissionDetails"`
		IsReserved            int    `json:"isReserved"`
		StockOnline           int    `json:"stockOnline"`
		CreatedAt             int    `json:"createdAt"`
		UpdatedAt             int    `json:"updatedAt"`
		//LastChangeAM             time.Time   `json:"lastChangeAM"`
		//CreatedAM                time.Time   `json:"createdAM"`
		LastRTFAM                string      `json:"lastRTFAM"`
		ReservationExpiredAt     int         `json:"reservationExpiredAt"`
		ImportedImages           int         `json:"importedImages"`
		StorageFolder            string      `json:"storageFolder"`
		AmCarObject              string      `json:"amCarObject"`
		AutovitID                interface{} `json:"autovitID"`
		AutovitStatus            interface{} `json:"autovitStatus"`
		AutovitImageCollectionID interface{} `json:"autovitImageCollectionID"`
		AutovitGeneration        interface{} `json:"autovitGeneration"`
		AutovitVersion           interface{} `json:"autovitVersion"`
		AutovitURL               interface{} `json:"autovitURL"`
		AutovitUpdatedAt         interface{} `json:"autovitUpdatedAt"`
		AutovitCreatedAt         interface{} `json:"autovitCreatedAt"`
		CarVerticalFile          string      `json:"carVerticalFile"`
		Model                    struct {
			Id               int    `json:"id"`
			CarBrandID       int    `json:"carBrandID"`
			Name             string `json:"name"`
			LastImportedDate int    `json:"lastImportedDate"`
			CreatedAt        int    `json:"createdAt"`
			UpdatedAt        int    `json:"updatedAt"`
			DisplayPriority  int    `json:"displayPriority"`
			Slug             string `json:"slug"`
			DisplayName      string `json:"displayName"`
			SameNameList     string `json:"sameNameList"`
			ModelType        string `json:"modelType"`
			ForTestDrive     int    `json:"forTestDrive"`
			NotInFilters     int    `json:"notInFilters"`
		} `json:"model"`
		Brand struct {
			Id               int           `json:"id"`
			Name             string        `json:"name"`
			Slug             string        `json:"slug"`
			SearchPriority   int           `json:"searchPriority"`
			LastImportedDate int           `json:"lastImportedDate"`
			CreatedAt        int           `json:"createdAt"`
			UpdatedAt        int           `json:"updatedAt"`
			AliasAM          []interface{} `json:"aliasAM"`
		} `json:"brand"`
		CarsGallery []struct {
			Id               int         `json:"id"`
			CarID            int         `json:"carID"`
			Type             string      `json:"type"`
			CarMediaURL      string      `json:"carMediaURL"`
			AmKey            int         `json:"amKey"`
			LastImportedDate interface{} `json:"lastImportedDate"`
			CreatedAt        int         `json:"createdAt"`
			UpdatedAt        int         `json:"updatedAt"`
			LastImageAM      time.Time   `json:"lastImageAM"`
		} `json:"carsGallery"`
		MinSalePrice int `json:"minSalePrice"`
		HasBadge     int `json:"hasBadge"`
	} `json:"response"`
}
