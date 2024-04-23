package main

import (
	"carscraper/pkg/amconfig"
	"carscraper/pkg/errorshandler"
	"carscraper/pkg/scraping/scrapingservices"
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {

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

	scrapingMapper := scrapingservices.NewScrapingAdaptersMapper()

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
		scrapingservices.WithMarketService("autoscout", rodScrapingService),
		scrapingservices.WithMarketService("autotracknl", rodScrapingService),
		scrapingservices.WithMarketService("olx", jsonScrapingService))

	sjh.Start()

	<-done
	fmt.Println("exiting")
}
