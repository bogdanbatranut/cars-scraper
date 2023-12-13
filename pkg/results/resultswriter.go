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

func (w ResultsWriter) WriteAds(ads []jobs.Ad, marketID uint) error {
	var dbAds []adsdb.Ad
	for _, ad := range ads {
		dbAd, err := w.adapter.ToActiveDBAd(ad, marketID)
		if err != nil {
			return err
		}
		dbAds = append(dbAds, *dbAd)
	}
	err := w.repo.WriteResults(dbAds)
	if err != nil {
		return err
	}
	return nil
}
