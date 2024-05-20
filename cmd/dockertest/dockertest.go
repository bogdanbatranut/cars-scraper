package main

import (
	"fmt"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/input"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/utils"
)

func main() {
	l := launcher.New().
		Headless(false).
		Devtools(true)

	defer l.Cleanup()

	url := l.MustLaunch()

	// Trace shows verbose debug information for each action executed
	// SlowMotion is a debug related function that waits 2 seconds between
	// each action, making it easier to inspect what your code is doing.
	browser := rod.New().
		ControlURL(url).
		Trace(true).
		SlowMotion(2 * time.Second).
		MustConnect()

	// ServeMonitor plays screenshots of each tab. This feature is extremely
	// useful when debugging with headless mode.
	// You can also enable it with flag "-rod=monitor"
	launcher.Open(browser.ServeMonitor(""))

	defer browser.MustClose()

	page := browser.MustPage("https://github.com/")

	page.MustElement("input").MustInput("git").MustType(input.Enter)

	text := page.MustElement(".codesearch-results p").MustText()

	fmt.Println(text)

	utils.Pause() // pause goroutine
}

func main__() {
	// This example is to launch a browser remotely, not connect to a running browser remotely,
	// to connect to a running browser check the "../connect-browser" example.
	// Rod provides a docker image for beginners, run the below to start a launcher.Manager:
	//
	//     docker run --rm -p 7317:7317 ghcr.io/go-rod/rod
	//
	// For available CLI flags run: docker run --rm ghcr.io/go-rod/rod rod-manager -h
	// For more information, check the doc of launcher.Manager
	//l := launcher.MustNewManaged("http://dev.auto-mall.ro:7317")
	l := launcher.MustNewManaged("")

	// You can also set any flag remotely before you launch the remote browser.
	// Available flags: https://peter.sh/experiments/chromium-command-line-switches
	l.Set("disable-gpu").Delete("disable-gpu")

	// Launch with headful mode
	l.Headless(false).XVFB("--server-num=5", "--server-args=-screen 0 1600x900x16")

	browser := rod.New().Client(l.MustClient()).MustConnect()

	// You may want to start a server to watch the screenshots of the remote browser.
	//launcher.Open(browser.ServeMonitor("dev.auto-mall.ro:7317"))

	fmt.Println(
		browser.MustPage("https://developer.mozilla.org").MustEval("() => document.title"),
	)

	// Launch another browser with the same docker container.
	ll := launcher.MustNewManaged("http://dev.auto-mall.ro:7317")

	// You can set different flags for each browser.
	ll.Set("disable-sync").Delete("disable-sync")

	anotherBrowser := rod.New().Client(ll.MustClient()).MustConnect()

	fmt.Println(
		anotherBrowser.MustPage("https://go-rod.github.io").MustEval("() => document.title"),
	)

	utils.Pause()
}
