package mobiledecollycollector

import (
	"carscraper/pkg/jobs"
	"fmt"

	"github.com/gocolly/colly"
)

type MobileDECollyCollector struct {
}

func NewMobileDECollyCollector() *MobileDECollyCollector {
	return &MobileDECollyCollector{}
}

func (collector MobileDECollyCollector) GetCollyCollector(job jobs.SessionJob) *colly.Collector {

	collyCollector := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36"),
	)

	collyCollector.OnRequest(func(req *colly.Request) {
		fmt.Println("MOBILE Visiting", req.URL.String())
		fmt.Println("applying headers")
		req.Headers.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
		req.Headers.Add("Accept-Encoding", "gzip, deflate, br")
		req.Headers.Add("Accept-Language", "en-GB,en;q=0.9")
		req.Headers.Add("Sec-Ch-Ua", "\"Google Chrome\";v=\"119\", \"Chromium\";v=\"119\", \"Not?A_Brand\";v=\"24\"")
		req.Headers.Add("Sec-Ch-Ua-Mobile", "?0")
		req.Headers.Add("Sec-Ch-Ua-Platform", "\"macOS\"")
		req.Headers.Add("Sec-Fetch-Dest", "document")
		req.Headers.Add("Sec-Fetch-Mode", "navigate")
		req.Headers.Add("Sec-Fetch-Site", "none")
		req.Headers.Add("Sec-Fetch-User", "?1")
		req.Headers.Add("Upgrade-Insecure-Requests", "1")
		req.Headers.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36")

	})

	//err := collyCollector.Visit(url)
	//if err != nil {
	//	return nil, err
	//}
	//collyCollector.Wait()

	return collyCollector
}
