package icollector

import "github.com/go-rod/rod"

type IRodPageProcessor interface {
	ProcessPage(page *rod.Page) AdsResults
}
