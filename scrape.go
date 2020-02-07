package main

import (
	"fmt"
	"time"

	"github.com/gocolly/colly"
	"github.com/gocolly/colly/debug"
)

type Parameter struct {
	AllOptions        string // div.gcode div.params table td.arg code
	Description       string // div.gcode div.params table td:nth-child(2) p:first
	Option            string // div.gcode div.params table td:nth-child(2) ul > li > code
	OptionDescription string // div.gcode div.params table td:nth-child(2) ul > li > p
}

type Example struct {
	Code        string // div.gcode div.examples code
	Description string // div.gcode div.examples p
}

type Command struct {
	Title           string   // div.gcode div.meta h1
	Code            string   // div.gcode div.usage code
	Description     string   // div.gcode div.long p
	FirmwareVersion string   // div.gcode div.meta span.label-success
	Notes           []string // div.gcode div.notes p
	Examples        []Example
	Parameters      []Parameter
}

func main() {
	// Instantiate a new collector with support for caching,
	// debugging, and scraping websites asyncronously.
	c := colly.NewCollector(
		colly.UserAgent("gcode.fyi"),
		colly.CacheDir("./.cache"),
		colly.Debugger(&debug.LogDebugger{}),
		colly.Async(true),
	)

	// Limit the maximum parallelism to 2.
	// Required to limit simultaneous requests.
	c.Limit(&colly.LimitRule{DomainGlob: "*", Parallelism: 2, RandomDelay: 5 * time.Second})

	// Scrape marlin.org for links to docs for each command.
	c.OnHTML("div.item > h2 > a", func(e *colly.HTMLElement) {
		marlinDocsLink := e.Attr("href")
		e.Request.Visit(marlinDocsLink)
	})

	c.OnHTML("div.gcode", func(e *colly.HTMLElement) {
		found := e.DOM.Find("> #gcode-header")
		if found != nil {
			fmt.Println("is command page")

		}
	})

	// Handle any errors that occur.
	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("[FAILED]", r.Request.URL, "failed:", r, "\nERROR:", err)
	})

	// Start scraping on the specified URL.
	c.Visit("https://marlinfw.org/meta/gcode/")

	// Wait until threads are finished.
	// Required when collector configured with Async(true).
	c.Wait()
}
