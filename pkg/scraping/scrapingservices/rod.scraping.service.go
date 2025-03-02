package scrapingservices

import (
	"carscraper/pkg/amconfig"
	"carscraper/pkg/jobs"
	"carscraper/pkg/logging"
	"carscraper/pkg/scraping/icollector"
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
	context                       context.Context
	jobChannel                    chan jobs.SessionJob
	resultsChannel                chan jobs.AdsPageJobResult
	scrapingMapper                IScrapingMapper
	browserMapper                 MarketBrowserMapper
	browserLauncher               *launcher.Launcher
	urlBuilderMapper              urlbuilder.URLBuilderMapper
	browser                       *rod.Browser
	currentJobAvailabilityChannel chan bool
	loggingService                *logging.ScrapeLoggingService
	dockerBrowserURL              string
	useDockerBrowser              bool
}

func NewRodScrapingService(ctx context.Context, scrapingMapper IScrapingMapper, cfg amconfig.IConfig) *RodScrapingService {

	urlbBuilderMapper := urlbuilder.NewURLBuildMapper()

	autotrackURLBuilder := autotrack.NewURLBuilder()
	urlbBuilderMapper.AddBuilder("autotracknl", autotrackURLBuilder)

	//mobileDeURLBuilder := mobile.NewURLBuilder()
	autoscoutURLBuilder := autoscout.NewURLBuilder()
	urlbBuilderMapper.AddBuilder("autoscout", autoscoutURLBuilder)

	var br *rod.Browser

	//isProd := cfg.GetBool(amconfig.AppIsProd)
	useDockerBrowser := cfg.GetBool(amconfig.PageScraperUseDockerRod)
	dockerBrowserURL := cfg.GetString(amconfig.PageScraperDockerContainerURL)
	if useDockerBrowser {
		log.Println("Using docker browser at ", dockerBrowserURL)
		br = connectToDockerBrowser(dockerBrowserURL)
	} else {
		log.Println("Using local browser")
		br = startLocalBrowserWithMonitor()
		//br = startBrowser()
	}
	//br := startBrowser()

	return &RodScrapingService{
		context:                       ctx,
		jobChannel:                    make(chan jobs.SessionJob),
		scrapingMapper:                scrapingMapper,
		browserMapper:                 MarketBrowserMapper{marketBrowsers: nil},
		resultsChannel:                make(chan jobs.AdsPageJobResult),
		urlBuilderMapper:              *urlbBuilderMapper,
		browser:                       br,
		currentJobAvailabilityChannel: make(chan bool),
		loggingService:                logging.NewScrapeLoggingService(cfg),
	}
}

func (rss RodScrapingService) GetCurrentJobExecutionAvailabilityChannel() chan bool {
	return rss.currentJobAvailabilityChannel
}

func connectToDockerBrowser(url string) *rod.Browser {
	//l, err := launcher.NewManaged("pensive_mendel")
	//if err != nil {
	//	panic(err)
	//}
	//l.Headless(false).XVFB("--server-num=5", "--server-args=-screen 0 1600x900x16")
	//log.Println("connecting to docker")
	//u, err := launcher.ResolveURL("")
	//if err != nil {
	//	log.Println(" -------- ")
	//	panic(err)
	//}
	//
	//browser := rod.New().ControlURL(u).MustConnect()
	//return browser

	//l, err := launcher.NewManaged("http://rod-chromium:7317")
	//l, err := launcher.NewManaged("http://dev.auto-mall.ro:7317")
	l, err := launcher.NewManaged(url)
	if err != nil {
		panic(err)
	}
	l.Headless(false).XVFB("--server-num=5", "--server-args=-screen 0 1600x900x16")

	browser := rod.New().Client(l.MustClient()).Trace(true).MustConnect()
	return browser

}

func startLocalBrowserWithMonitor() *rod.Browser {
	l := launcher.New().
		Headless(false).
		//XVFB("--server-num=5", "--server-args=-screen 0 1600x900x16").
		Devtools(true)

	url := l.MustLaunch()

	// Trace shows verbose debug information for each action executed
	// SlowMotion is a debug related function that waits 2 seconds between
	// each action, making it easier to inspect what your code is doing.
	browser := rod.New().
		ControlURL(url).
		//MustIncognito().
		Trace(false).
		SlowMotion(2 * time.Second).
		MustConnect()

	// ServeMonitor plays screenshots of each tab. This feature is extremely
	// useful when debugging with headless mode.
	// You can also enable it with flag "-rod=monitor"
	launcher.Open(browser.ServeMonitor(""))
	return browser
}

func startBrowser() *rod.Browser {
	//l := launcher.MustNewManaged("http://dev.auto-mall.ro:7317")
	l, err := launcher.NewManaged("http://dev.auto-mall.ro:7317")
	if err != nil {
		panic(err)
	}
	l.Headless(false).XVFB("--server-num=5", "--server-args=-screen 0 1600x900x16")

	browser := rod.New().Client(l.MustClient()).Trace(false).MustConnect()
	return browser
}

func (rss RodScrapingService) StartFake() {
	log.Println("Rod Scraping Service Fake Start")
	//rss.browser = startBrowser()

	//launcher.Open(rss.browser.ServeMonitor(""))
	go func() {
		for {
			select {
			case job := <-rss.jobChannel:
				rss.processFakeJob(job)
			case <-rss.context.Done():
				log.Println("Rod Scraping Service Fake Terminating...")
				return
			}
		}
	}()
}

func (rss RodScrapingService) Start() {
	log.Println("Rod Scraping Service Start")

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

func (rss RodScrapingService) processFakeJob(job jobs.SessionJob) {
	log.Println("ROD SLEEPING")
	time.Sleep(2 * time.Second)

	isLastPage := job.Market.PageNumber > 2

	results := icollector.AdsResults{
		Ads:        nil,
		IsLastPage: isLastPage,
		Error:      errors.New("FAKE JOB"),
	}

	adResult := jobs.AdsPageJobResult{
		RequestedScrapingJob: job,
		PageNumber:           job.Market.PageNumber,
		IsLastPage:           results.IsLastPage,
		Success:              results.Error == nil,
		Data:                 results.Ads,
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
}

func (rss RodScrapingService) processJob(job jobs.SessionJob) {
	//urlBuilder := autoscout.NewURLBuilder(job.Criteria)
	criteriaLog, err := rss.loggingService.GetCriteriaLog(job.SessionID, job.CriteriaID, job.MarketID)
	if err != nil {
		log.Println(err.Error())
	}
	pageLog, err := rss.loggingService.CreatePageLog(criteriaLog, job, "", job.Market.PageNumber)
	if err != nil {
		panic(err)
	}
	log.Println("ROD Executing :", job.ToString())
	urlBuilder := rss.urlBuilderMapper.GetURLBuilder(job.Market.Name)
	url := urlBuilder.GetURL(job)
	if url == nil {
		panic(errors.New("could not build url for scraping"))
	}
	log.Println("ROD Getting data from URL : ", *url)

	err = rss.loggingService.PageLogSetVisitURL(pageLog, *url)
	if err != nil {
		panic(err)
	}

	var page *rod.Page
	//page = rss.browser.SlowMotion(1 * time.Second).MustPage(*url).MustWaitDOMStable()
	page, err = rss.browser.Page(proto.TargetCreateTarget{
		URL:                     *url,
		Width:                   nil,
		Height:                  nil,
		BrowserContextID:        "",
		EnableBeginFrameControl: false,
		NewWindow:               false,
		Background:              false,
		ForTab:                  false,
	})
	if err != nil {
		log.Println(err)
	}
	//page = rss.browser.SlowMotion(1 * time.Second).MustPage(*url).MustWaitDOMStable()
	adapter := rss.scrapingMapper.GetRodMarketAdsAdapter(job.Market.Name)
	if err != nil {
		panic(err)
	}

	results := adapter.GetAds(page)
	if results.Ads != nil {
		log.Println("ROD SERVICE Got ", len(*results.Ads), " ads")
	} else {
		log.Println("ROD SERVICE Got NULL ads")

	}
	err = page.Close()
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
		err2 := rss.loggingService.PageLogSetPageScraped(pageLog, 0, adResult.IsLastPage)
		if err2 != nil {
			log.Println(err2.Error())
		}
		return
	}

	err2 := rss.loggingService.PageLogSetPageScraped(pageLog, len(*adResult.Data), adResult.IsLastPage)
	if err2 != nil {
		log.Println(err2.Error())
	}

	if adResult.IsLastPage {
		err := closeAllPages(rss.browser)
		if err != nil {
			panic(err)
		}

		if err != nil {
			return
		}
	}

	go func(res jobs.AdsPageJobResult) {
		go func(tmpch chan jobs.AdsPageJobResult, r jobs.AdsPageJobResult) {
			tmpch <- r
		}(rss.resultsChannel, res)
	}(adResult)
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
