package results

import (
	"carscraper/pkg/adsdb"
	"carscraper/pkg/jobs"
)

type ResultsWriter struct {
	adapter IAdsResultsAdapter
	repo    ResultsRepository
}

func NewResultsWriter(iadapter IAdsResultsAdapter, rr ResultsRepository) *ResultsWriter {
	return &ResultsWriter{adapter: iadapter, repo: rr}
}

func (w ResultsWriter) WriteAds(ads []jobs.Ad, marketID uint, criteriaID uint) (*[]uint, error) {
	var dbAds []adsdb.Ad
	for _, ad := range ads {
		dbAd, err := w.adapter.ToActiveDBAd(ad, marketID, criteriaID)
		if err != nil {
			return nil, err
		}
		dbAds = append(dbAds, *dbAd)
	}
	adsIds, err := w.repo.Write(dbAds)
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
