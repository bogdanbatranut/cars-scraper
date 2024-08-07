package main

import (
	"carscraper/pkg/amconfig"
	"carscraper/pkg/errorshandler"
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg, err := amconfig.NewViperConfig()
	errorshandler.HandleErr(err)
	_, cancel := context.WithCancel(context.Background())

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

	testvar := cfg.GetString(amconfig.TestVar)

	log.Println("testvar :", testvar)

	log.Println("Is Dev: ", cfg.GetBool(amconfig.AppIsDev))
	log.Println("IS PROD", cfg.GetBool(amconfig.AppIsProd))

}
