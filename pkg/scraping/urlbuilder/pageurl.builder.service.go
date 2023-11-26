package urlbuilder

import "fmt"

type IPageURLBuilder interface {
	GetURL(pageURL PageURL) string
}

type PageURLBuilderImplementations struct {
	implementations map[string]IPageURLBuilder
}

func NewPageURLBuilderImplementations() PageURLBuilderImplementations {
	implementations := PageURLBuilderImplementations{}
	implementations.addPageURLBuilderImplementations()
	return implementations
}

func (pbi PageURLBuilderImplementations) addPageURLBuilderImplementations() {
	pbi.implementations["autovit"] = AutovitPageURLBuilder{}
	pbi.implementations["mobile.de"] = MobileDEPageURLBuilder{}
}

type AutovitPageURLBuilder struct {
}

func (apb AutovitPageURLBuilder) GetURL(pageURL PageURL) string {
	pURL := fmt.Sprintf("%s/autoturisme/%s/%s/de-la%d?search%%5Bfilter_enum_fuel_type%%5D=%s&search%%5Bfilter_float_mileage%%3Ato%%5D=%d&page=%d",
		pageURL.MarketURL, pageURL.CarBrand, pageURL.CarModel, *pageURL.YearFrom, *pageURL.Fuel, *pageURL.KmTo, pageURL.PageNumber)
	return pURL
}

type MobileDEPageURLBuilder struct {
}

func (mdepb MobileDEPageURLBuilder) GetURL(pageURL PageURL) string {
	// https://www.mobile.de/ro/automobil/porsche-panamera/vhc:car,srt:price,sro:asc,ms1:20100_4_,frn:2018,ful:hybrid,mlx:125000
	//https://www.mobile.de/ro/automobil/mercedes-benz-clasa-gle/vhc:car,srt:price,sro:asc,ms1:17200_-58_,frn:2019,ful:diesel!hybrid,mlx:125000
	//https://www.mobile.de/ro/automobil/mercedes-benz-gle-350/vhc:car,pgn:2,pgs:10,srt:price,sro:asc,ms1:17200_251_,frn:2019,ful:diesel!hybrid,mlx:125000,dmg:false
	pURL := fmt.Sprintf("%s/ro/automobil/%s-%s/vhc:car,pgn:%d,srt:price,sro:asc,ms1:17200_-58_,frn%d,ful:%s,mlx:%d",
		pageURL.MarketURL, pageURL.CarBrand, pageURL.CarModel, pageURL.PageNumber, *pageURL.YearFrom, *pageURL.Fuel, *pageURL.KmTo)
	return pURL
}
