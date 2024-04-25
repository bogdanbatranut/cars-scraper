package scrapingservices

import (
	"carscraper/pkg/amconfig"
	"carscraper/pkg/jobs"
	"carscraper/pkg/scraping/markets/autoscout"
	"carscraper/pkg/scraping/markets/autotrack"
	"carscraper/pkg/scraping/urlbuilder"
	"context"
	"errors"
	"log"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/proto"
	"github.com/go-rod/rod/lib/utils"
	"github.com/ysmood/gson"
)

type RodBrowserManager struct {
	browsers []rod.Browser
}

type MarketBrowserMapper struct {
	//marketBrowsers map[string]*rod.Browser
	marketBrowsers *rod.Browser
}

func (mapper MarketBrowserMapper) addMarketBrowser(market string, browser *rod.Browser) {
	mapper.marketBrowsers = browser
	//mapper.marketBrowsers[market] = browser
}

func (mapper MarketBrowserMapper) getMarketBrowser(market string) *rod.Browser {
	return mapper.marketBrowsers
	//return mapper.marketBrowsers[market]
}

type RodScrapingService struct {
	context          context.Context
	jobChannel       chan jobs.SessionJob
	resultsChannel   chan jobs.AdsPageJobResult
	scrapingMapper   IScrapingMapper
	browserMapper    MarketBrowserMapper
	browserLauncher  *launcher.Launcher
	urlBuilderMapper urlbuilder.URLBuilderMapper
	browser          *rod.Browser
}

func NewRodScrapingService(ctx context.Context, scrapingMapper IScrapingMapper, cfg amconfig.IConfig) *RodScrapingService {

	urlbBuilderMapper := urlbuilder.NewURLBuildMapper()

	autotrackURLBuilder := autotrack.NewURLBuilder()
	urlbBuilderMapper.AddBuilder("autotracknl", autotrackURLBuilder)

	//mobileDeURLBuilder := mobile.NewURLBuilder()
	autoscoutURLBuilder := autoscout.NewURLBuilder()
	urlbBuilderMapper.AddBuilder("autoscout", autoscoutURLBuilder)

	//br := startLocalBrowserWithMonitor()
	br := startBrowser()
	return &RodScrapingService{
		context:        ctx,
		jobChannel:     make(chan jobs.SessionJob),
		scrapingMapper: scrapingMapper,
		browserMapper: MarketBrowserMapper{
			marketBrowsers: nil},
		resultsChannel:   make(chan jobs.AdsPageJobResult),
		urlBuilderMapper: *urlbBuilderMapper,
		browser:          br,
	}
}

func startLocalBrowserWithMonitor() *rod.Browser {
	l := launcher.New().
		Headless(true).
		Devtools(false)

	url := l.MustLaunch()

	// Trace shows verbose debug information for each action executed
	// SlowMotion is a debug related function that waits 2 seconds between
	// each action, making it easier to inspect what your code is doing.
	browser := rod.New().
		ControlURL(url).
		MustIncognito().
		Trace(false).
		//SlowMotion(2 * time.Second).
		MustConnect()

	// ServeMonitor plays screenshots of each tab. This feature is extremely
	// useful when debugging with headless mode.
	// You can also enable it with flag "-rod=monitor"

	return browser
}

func startBrowser() *rod.Browser {
	l := launcher.MustNewManaged("http://dev.auto-mall.ro:7317")
	l.Headless(true).XVFB("--server-num=5", "--server-args=-screen 0 1600x900x16")

	browser := rod.New().Client(l.MustClient()).MustIncognito().Trace(false).MustConnect()
	return browser
}

func (rss RodScrapingService) Start() {
	log.Println("Rod Scraping Service Start")
	//rss.browser = startBrowser()

	//launcher.Open(rss.browser.ServeMonitor(""))
	go func() {
		for {
			select {
			case job := <-rss.jobChannel:
				rss.processJob(job)
			case <-rss.context.Done():
				log.Println("Rod Scraping Service Terminating...")
				return
			}
		}
	}()
}

func (rss RodScrapingService) AddJob(job jobs.SessionJob) {
	rss.jobChannel <- job
}

func (rss RodScrapingService) GetResultsChannel() *chan jobs.AdsPageJobResult {
	return &rss.resultsChannel
}

func (rss RodScrapingService) processJob(job jobs.SessionJob) {
	//urlBuilder := autoscout.NewURLBuilder(job.Criteria)
	urlBuilder := rss.urlBuilderMapper.GetURLBuilder(job.Market.Name)
	url := urlBuilder.GetURL(job)
	if url == nil {
		panic(errors.New("could not build url for scraping"))
	}
	log.Println("Getting data from URL : ", *url)
	var page *rod.Page
	//page = rss.browser.SlowMotion(1 * time.Second).MustPage(*url).MustWaitDOMStable()
	page = rss.browser.SlowMotion(1 * time.Second).MustPage(*url).MustWaitDOMStable()
	adapter := rss.scrapingMapper.GetRodMarketAdsAdapter(job.Market.Name)
	results := adapter.GetAds(page)
	err := page.Close()
	if err != nil {
		log.Println(err)
		//return
	}
	adResult := jobs.AdsPageJobResult{
		RequestedScrapingJob: job,
		PageNumber:           job.Market.PageNumber,
		IsLastPage:           results.IsLastPage,
		Success:              results.Error == nil,
		Data:                 results.Ads,
	}
	if results.Ads == nil {
		return
	}
	if adResult.IsLastPage {
		err := closeAllPages(rss.browser)
		if err != nil {
			panic(err)
		}
	}
	go func(res jobs.AdsPageJobResult) {
		go func(tmpch chan jobs.AdsPageJobResult, r jobs.AdsPageJobResult) {
			tmpch <- r
		}(rss.resultsChannel, res)
	}(adResult)

	log.Println("ROD pushed results to channel")

}

func closeAllPages(browser *rod.Browser) error {
	pages, err := browser.Pages()
	if err != nil {
		return err
	}
	for _, page := range pages {
		err := page.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

func ScreenShot(page *rod.Page) {

	// capture entire browser viewport, returning jpg with quality=90
	img, err := page.ScrollScreenshot(&rod.ScrollScreenshotOptions{
		Format:  proto.PageCaptureScreenshotFormatJpeg,
		Quality: gson.Int(90),
	})
	if err != nil {
		panic(err)
	}

	_ = utils.OutputFile("rod_scraping.jpg", img)
}
