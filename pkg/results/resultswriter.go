package results

import (
	"carscraper/pkg/adsdb"
	"carscraper/pkg/jobs"
	"carscraper/pkg/repos"
)

type ResultsWriter struct {
	adapter IAdsResultsAdapter
	repo    repos.AdsRepository
}

func NewResultsWriter(iadapter IAdsResultsAdapter, adsRepo repos.AdsRepository) *ResultsWriter {
	return &ResultsWriter{adapter: iadapter, repo: adsRepo}
}

func (w ResultsWriter) WriteAds(ads []jobs.Ad, marketID uint, criteriaID uint) (*[]uint, error) {
	var dbAds []adsdb.Ad
	for _, ad := range ads {
		if ad.Price == 0 {
			continue
		}
		dbAd, err := w.adapter.ToActiveDBAd(ad, marketID, criteriaID)
		if err != nil {
			return nil, err
		}
		dbAds = append(dbAds, *dbAd)
	}
	adsIds, err := w.repo.Upsert(dbAds)
	if err != nil {
		return nil, err
	}
	return adsIds, nil
}

func (w ResultsWriter) GetAllAdsIDs(marketID uint, criteriaID uint) *[]uint {
	return w.repo.GetAllAdsIDs(marketID, criteriaID)
}

func (w ResultsWriter) DeleteAd(adID uint) {
	w.repo.DeleteAd(adID)
}
