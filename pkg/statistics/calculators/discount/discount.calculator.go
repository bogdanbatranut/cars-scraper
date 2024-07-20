package discount

import (
	"carscraper/pkg/adsdb"
	"carscraper/pkg/repos"
	"carscraper/pkg/statistics/calculators/age"
	"carscraper/pkg/statistics/calculators/helpers"
)

func GetCalculatedAverageDealerDiscountPercent(repo repos.IAdsRepository, sellerID uint) float64 {
	sellerAds := repo.GetSellerAds(sellerID)
	if sellerAds == nil || len(*sellerAds) == 0 {
		return 0
	}
	var total float64
	for _, sellerAd := range *sellerAds {
		_, discountPercent := CalculateAdDiscount(sellerAd)
		total += discountPercent
		// for each ad, calulate avg discount percent
	}
	avg := total / float64(len(*sellerAds))
	return helpers.ToFixed(avg, 2)
}

func CalculateAverageDealersDiscount(ads []adsdb.Ad) map[uint]float64 {

	var dealersWithAdsDailyDiscountsAmmounts map[uint][]float64
	var dealersAverageDailyDiscountsAmmounts map[uint]float64

	dealersWithAdsDailyDiscountsAmmounts = make(map[uint][]float64)
	dealersAverageDailyDiscountsAmmounts = make(map[uint]float64)

	for _, ad := range ads {
		dealersWithAdsDailyDiscountsAmmounts[ad.SellerID] = append(dealersWithAdsDailyDiscountsAmmounts[ad.SellerID],
			CalculateDailyDiscount(ad),
		)
	}

	for dealerID, dailyDiscounts := range dealersWithAdsDailyDiscountsAmmounts {
		dealersAverageDailyDiscountsAmmounts[dealerID] = helpers.ToFixed(helpers.Average(dailyDiscounts), 2)
	}
	return dealersAverageDailyDiscountsAmmounts
}

func CalculateDailyDiscount(ad adsdb.Ad) float64 {
	firstPrice := ad.Prices[0].Price
	lastPrice := ad.Prices[len(ad.Prices)-1].Price
	adDuration := age.CalculateAdAge(ad)
	discountAvgAmount := (firstPrice - lastPrice) / adDuration
	return helpers.ToFixed(float64(discountAvgAmount), 2)
}

func CalculateAdDiscount(ad adsdb.Ad) (int, float64) {
	if len(ad.Prices) == 0 {
		return 0, 0
	}
	discVal := ad.Prices[0].Price - ad.Prices[len(ad.Prices)-1].Price
	discPercent := float64(discVal) / float64(ad.Prices[0].Price) * 100
	ro := helpers.ToFixed(discPercent, 2)
	return discVal, ro
}
