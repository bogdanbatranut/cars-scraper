package main

import (
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/utils"
)

func main() {
	u := launcher.MustResolveURL("http://dev.auto-mall.ro:7317")

	browser := rod.New().ControlURL(u).MustConnect()
	launcher.Open(browser.ServeMonitor(""))
	utils.Pause()
}
