package main

import (
	"sync"
	"time"
)

// the in-memory page cache
var cachedPages = make(map[string]Page)

// prevent concurrent writes to the cache
var pmutex = &sync.RWMutex{}

// the in-memory index cache object
var indexCache = indexPage{}

// mutex for the index cache
var imutex = &sync.RWMutex{}

type cacher interface {
	cache()
	checkCache()
}

// Page cache object definition
type Page struct {
	Longname  string
	Shortname string
	Title     string
	Desc      string
	Author    string
	Modtime   time.Time
	Body      []byte
	Raw       pagedata
}

// index cache object definition
type indexPage struct {
	Modtime   time.Time
	LastTally time.Time
	Body      []byte
	Raw       pagedata
}

type pagedata []byte

// Creates a filled page object
func newPage(longname, shortname, title, author, desc string, modtime time.Time, body []byte, raw pagedata) *Page {

	return &Page{
		Longname:  longname,
		Shortname: shortname,
		Title:     title,
		Author:    author,
		Desc:      desc,
		Modtime:   modtime,
		Body:      body,
		Raw:       raw}

}

// Creates a page object with the minimal number of fields filled
func newBarePage(longname, shortname string) *Page {
	return &Page{
		Longname:  longname,
		Shortname: shortname,
		Title:     "",
		Author:    "",
		Desc:      "",
		Modtime:   time.Time{},
		Body:      nil,
		Raw:       nil,
	}
}
