package marketadapters

import (
	"carscraper/pkg/scraping/icollector"

	"github.com/go-rod/rod"
)

type IRodMarketAdsAdapter interface {
	GetAds(page *rod.Page) *icollector.AdsResults
}
