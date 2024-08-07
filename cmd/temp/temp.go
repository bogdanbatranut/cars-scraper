package main

import (
	"carscraper/pkg/amconfig"
	"carscraper/pkg/errorshandler"
	"log"
)

func main() {
	cfg, err := amconfig.NewViperConfig()
	errorshandler.HandleErr(err)
	log.Println(cfg.GetBool(amconfig.AppIsDev))
	log.Println(cfg.GetBool(amconfig.AppIsProd))
}
