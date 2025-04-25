package oferte_bmw

import (
	"carscraper/pkg/scraping/icollector"

	"github.com/go-rod/rod"
)

type OferteBMWRodAdapter struct {
}

func NewOferteBMWRodAdapter() *OferteBMWRodAdapter {
	return &OferteBMWRodAdapter{}
}

func (a OferteBMWRodAdapter) GetAds(page *rod.Page) *icollector.AdsResults {
	return nil
}
