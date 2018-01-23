package main

import (
	"fmt"
	"net/url"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
)

var rootURL *url.URL

// When no element found in [URL|Doc]Queue, sleeps for idleTime before try again
const idleTime = time.Duration(1) * time.Second

// If have tried with empty queue for more than emptyThshld times,
// considers the jobs related are done
const emptyThshld = 3

// Number of web page crawlers,
// for this solution, number of page parsers could only be 1,
// otherwise would need to also synchronize on site.(SiteMap)
const numCrawler = 5

func main() {
	urlQ := &URLQueue{0, make([]string, 0), new(sync.Mutex)}
	docQ := &DocQueue{0, make([]*goquery.Document, 0), new(sync.Mutex)}
	site := make(SiteMap)

	root := "https://monzo.com/"
	rootURL, _ = url.Parse(root)

	urlQ.PutURL(root)

	var wg sync.WaitGroup
	wg.Add(1 + numCrawler)

	for i := 0; i < numCrawler; i++ {
		go crawler(urlQ, docQ, &wg)
	}
	go parser(docQ, urlQ, &site, &wg)
	go monitor(docQ, urlQ)

	wg.Wait()
	fmt.Println("[main] dump site map to hello-monzo.out...")
	site.DumpToFile("hello-monzo.out")
}

// Crawls a web page pointed by a URL fetched from urlQ,
// puts the pointer to retrieved HTML file into docQ
func crawler(urlQ *URLQueue, docQ *DocQueue, wg *sync.WaitGroup) {
	empty := 0
	for {
		u, ok := urlQ.GetURL()
		if !ok {
			if empty >= emptyThshld {
				break
			} else {
				empty++
				time.Sleep(idleTime)
			}
		} else {
			empty = 0
			doc, err := goquery.NewDocument(u)
			if err != nil {
				fmt.Printf("[crawler] error tried to fetch %s: %v\n", u, err)
			} else {
				docQ.PutDoc(doc)
				//fmt.Println("[crawler] fetched page ", u)
			}
		}
	}

	wg.Done()
	fmt.Println("[crawler] quits!")
}

// Parses a web page,
// puts any internal, not-visited URLs into urlQ,
// meanwhile records the relation of parent URL -> parsed child URL
func parser(docQ *DocQueue, urlQ *URLQueue, site *SiteMap, wg *sync.WaitGroup) {
	empty := 0
	for {
		doc, ok := docQ.GetDoc()
		if !ok {
			if empty >= emptyThshld {
				break
			} else {
				empty++
				time.Sleep(idleTime)
			}
		} else {
			empty = 0
			doc.Find("a").Each(func(inx int, a *goquery.Selection) {
				href, exists := a.Attr("href")
				if exists {
					parsed, sameDomain := ofSameDomain(href)
					// to rule out several fake URLs of the format "tel:[0-9]+"
					if parsed.Scheme == "http" || parsed.Scheme == "https" {
						inMap := site.PutInMap(doc.Url, parsed, sameDomain)
						if sameDomain && !inMap {
							urlQ.PutURL(parsed.String())
						}
					}
				}
			})
		}
	}

	wg.Done()
	fmt.Println("[parser] quits!")
}

// Checks wheter rawurl is an internal URL,
// when a relative path is found, resolves it again the root path
func ofSameDomain(rawurl string) (*url.URL, bool) {
	rawURL, err := url.Parse(rawurl)
	if err != nil {
		return nil, false
	}

	if rawURL.Hostname() == "" {
		return rootURL.ResolveReference(rawURL), true
	}

	return rawURL, rawURL.Hostname() == rootURL.Hostname()
}

// Monitors the growth of docQ and urlQ
func monitor(docQ *DocQueue, urlQ *URLQueue) {
	for {
		now := time.Now().Format(time.Stamp)
		fmt.Printf("[%s] len(document queue) = %d len(URL queue) = %d\n", now, docQ.counter, urlQ.counter)

		time.Sleep(time.Second)
	}
}
