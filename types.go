package main

import (
	"regexp"
	"sync"
	"time"
)

// The in-memory page cache
var pageCache = &pagesCache{
	mu:   new(sync.RWMutex),
	pool: make(map[string]*Page),
}

// The in-memory index cache
var indexCache = &indexCacheBlk{
	mu:   new(sync.RWMutex),
	page: new(indexPage),
}

// indexPage and Page types implement
// this interface, currently.
type cacher interface {
	cache() error
	checkCache() bool
}

type ipCtxKey int

const ctxKey ipCtxKey = iota

type pagesCache struct {
	mu   *sync.RWMutex
	pool map[string]*Page
}

type indexCacheBlk struct {
	mu   *sync.RWMutex
	page *indexPage
}

type confParams struct {
	mu                   sync.RWMutex
	port                 string
	pageDir              string
	assetsDir            string
	cssPath              string
	viewPath             string
	indexRefreshInterval string
	wikiName             string
	wikiDesc             string
	descSep              string
	titleSep             string
	iconPath             string
	indexFile            string
	reverseTally         bool
	validPath            *regexp.Regexp
	quietLogging         bool
	fileLogging          bool
	logFile              string
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
	Recache   bool
}

// Index cache object definition
type indexPage struct {
	Modtime   time.Time
	LastTally time.Time
	Body      []byte
	Raw       pagedata
}

// Type alias for methods and readability
// in certain situations.
type pagedata []byte

// Creates a filled page object
func newPage(longname, shortname, title, author, desc string, modtime time.Time, body []byte, raw pagedata, recache bool) *Page {
	return &Page{
		Longname:  longname,
		Shortname: shortname,
		Title:     title,
		Author:    author,
		Desc:      desc,
		Modtime:   modtime,
		Body:      body,
		Raw:       raw,
		Recache:   recache}

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
