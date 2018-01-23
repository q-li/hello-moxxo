package main

import (
	"sync"

	"github.com/PuerkitoBio/goquery"
)

// A URLQueue holds the URLs waiting to be fetched
type URLQueue struct {
	counter int
	queue   []string
	lock    *sync.Mutex
}

// PutURL puts an URL into the queue
func (u *URLQueue) PutURL(url string) {
	u.lock.Lock()
	defer u.lock.Unlock()

	u.counter++
	u.queue = append(u.queue, url)
}

// GetURL returns an URL from the queue, if there is still any,
// otherwise the returned indicator will be false
func (u *URLQueue) GetURL() (string, bool) {
	var ok bool
	var url string

	u.lock.Lock()
	defer u.lock.Unlock()

	if u.counter != 0 {
		ok = true
		url = u.queue[0]

		u.counter--
		u.queue = u.queue[1:]
	} else {
		ok = false
	}

	return url, ok
}

// A DocQueue holds pointers to HTML files waiting to be parsed
type DocQueue struct {
	counter int
	queue   []*goquery.Document
	lock    *sync.Mutex
}

// PutDoc puts a new pointer of fetched HTML file into the queue
func (d *DocQueue) PutDoc(doc *goquery.Document) {
	d.lock.Lock()
	defer d.lock.Unlock()

	d.counter++
	d.queue = append(d.queue, doc)
}

// GetDoc returns a HTML file from the queue, if there is still any,
// otherwise the returned indicator will be false
func (d *DocQueue) GetDoc() (*goquery.Document, bool) {
	var ok bool
	var doc *goquery.Document

	d.lock.Lock()
	defer d.lock.Unlock()

	if d.counter != 0 {
		ok = true
		doc = d.queue[0]

		d.counter--
		d.queue = d.queue[1:]
	} else {
		ok = false
	}

	return doc, ok
}
