package main

import (
	"carscraper/pkg/amconfig"
	"carscraper/pkg/errorshandler"
	"carscraper/pkg/logging"
	"carscraper/pkg/scraping/scrapingservices"
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	// to change the flags on the default logger
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	log.Println("starting page scraping service...")

	cfg, err := amconfig.NewViperConfig()
	errorshandler.HandleErr(err)
	ctx, cancel := context.WithCancel(context.Background())

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	done := make(chan bool, 1)

	go func() {
		sig := <-sigs
		fmt.Println()
		fmt.Println(sig)
		log.Println("canceling")
		cancel()
		done <- true
	}()

	fmt.Println("awaiting signal")

	loggingService := logging.NewScrapeLoggingService(cfg)

	scrapingMapper := scrapingservices.NewScrapingAdaptersMapper(loggingService)

	rodScrapingService := scrapingservices.NewRodScrapingService(ctx, scrapingMapper, cfg)
	rodScrapingService.Start()

	collyScrapingService := scrapingservices.NewCollyScrapingService(ctx, scrapingMapper)
	collyScrapingService.Start()

	jsonScrapingService := scrapingservices.NewJSONScrapingService(ctx, scrapingMapper)
	jsonScrapingService.Start()

	//sjh := scrapingservices.NewSessionJobHandler(ctx, cfg, rodScrapingService, collyScrapingService, jsonScrapingService)
	sjh := scrapingservices.NewSessionJobHandler(ctx, cfg,
		scrapingservices.WithMarketService("autovit", jsonScrapingService),
		scrapingservices.WithMarketService("mobile.de", collyScrapingService),
		scrapingservices.WithMarketService("bmw.de", collyScrapingService),
		scrapingservices.WithMarketService("autoscout", rodScrapingService),
		scrapingservices.WithMarketService("autotracknl", rodScrapingService),
		scrapingservices.WithMarketService("olx", jsonScrapingService),
		scrapingservices.WithMarketService("oferte_bmw", jsonScrapingService),
		scrapingservices.WithMarketService("tiriacauto", collyScrapingService),
		scrapingservices.WithMarketService("mercedes-benz.ro", jsonScrapingService),
		scrapingservices.WithMarketService("mercedes-benz.de", jsonScrapingService),
		scrapingservices.WithMarketService("autoklass.ro", jsonScrapingService),
	)

	sjh.Start()

	<-done
	fmt.Println("exiting")
}
