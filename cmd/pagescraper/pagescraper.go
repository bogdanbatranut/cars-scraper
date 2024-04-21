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

	sjh := scrapingservices.NewSessionJobHandler(ctx, cfg, rodScrapingService, collyScrapingService)
	sjh.StartWithoutMQ()

	<-done
	fmt.Println("exiting")
}
