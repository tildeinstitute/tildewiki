package main

import (
	"bufio"
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/spf13/viper"
)

// the in-memory page cache
var cachedPages = make(map[string]Page)

// prevent concurrent writes to the cache
var mutex = &sync.RWMutex{}

// Page struct for caching
type Page struct {
	Longname  string
	Shortname string
	Title     string
	Desc      string
	Author    string
	Modtime   time.Time
	Body      []byte
	Raw       []byte
}

var indexCache = indexPage{}
var inmutex = &sync.RWMutex{}

type indexPage struct {
	Modtime   time.Time
	LastTally time.Time
	Body      []byte
	Raw       []byte
}

// Creates a page struct
func newPage(longname, shortname, title, author, desc string, modtime time.Time, body, raw []byte) *Page {

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

// Loads a given wiki page and returns a page struct pointer.
// Used for building the initial cache and re-caching.
func buildPage(filename string) (*Page, error) {

	// open the page into *os.File
	file, err := os.Open(filename)
	if err != nil {
		log.Printf("%v\n", err)
		return nil, err
	}

	// the cleanup crew
	defer func() {
		err := file.Close()
		if err != nil {
			log.Printf("%v\n", err)
		}
	}()

	// stat the file to get mod time later
	stat, err := file.Stat()
	if err != nil {
		log.Printf("Couldn't stat %s: %v\n", filename, err)
	}

	// body holds the raw bytes from the file
	body, err := ioutil.ReadAll(file)
	if err != nil {
		log.Printf("%v\n", err)
	}

	// extract the file name from the path
	_, shortname := filepath.Split(filename)

	// get meta info on file from the header comment
	title, desc, author := getMeta(body)

	if title == "" {
		title = shortname
	}
	if desc != "" {
		desc = viper.GetString("DescSeparator") + " " + desc
	}
	if author != "" {
		author = "`by " + author + "`"
	}

	// longtitle is used in the <title> tags of the output html
	longtitle := title + " " + viper.GetString("TitleSeparator") + " " + viper.GetString("Name")

	// store the raw bytes of the document after parsing
	// from markdown to HTML.
	// keep the unparsed markdown for future use (maybe gopher?)
	bodydata := render(body, viper.GetString("CSS"), longtitle)
	return newPage(filename, shortname, title, author, desc, stat.ModTime(), bodydata, body), nil
}

// scan the page to the following fields in the
// header comment:
//		title:
//		author:
//		description:
func getMeta(body []byte) (string, string, string) {

	bytereader := bytes.NewReader(body)
	metafinder := bufio.NewScanner(bytereader)
	var title, desc, author string

	// scan the file line by line until it finds
	// the comments.
	for metafinder.Scan() {

		splitter := bytes.Split(metafinder.Bytes(), []byte(":"))

		switch string(bytes.ToLower(splitter[0])) {
		case "title":
			title = string(bytes.TrimSpace(splitter[1]))
		case "description":
			desc = string(bytes.TrimSpace(splitter[1]))
		case "author":
			author = string(bytes.TrimSpace(splitter[1]))
		default:
			continue
		}
		if title != "" && desc != "" && author != "" {
			return title, desc, author
		}

	}

	return title, desc, author
}

// generate the front page of the wiki
func genIndex() []byte {

	var err error
	indexpath := viper.GetString("AssetsDir") + "/" + viper.GetString("Index")

	stat, err := os.Stat(indexpath)
	if err != nil {
		log.Printf("Couldn't stat index: %v\n", err)
	}

	if indexCache.Modtime != stat.ModTime() {
		inmutex.Lock()
		indexCache.Raw, err = ioutil.ReadFile(indexpath)
		inmutex.Unlock()
		if err != nil {
			return []byte("Could not open \"" + indexpath + "\"")
		}

	}

	// body holds the bytes of the generated index page being sent to the client.
	// create the byte array and the buffer used to write to it
	body := make([]byte, 0)
	buf := bytes.NewBuffer(body)

	// scan the file line by line until it finds the anchor
	// comment. replace the anchor comment with a list of
	// wiki pages sorted alphabetically by title.
	inmutex.RLock()
	builder := bufio.NewScanner(bytes.NewReader(indexCache.Raw))
	inmutex.RUnlock()
	builder.Split(bufio.ScanLines)

	for builder.Scan() {
		if bytes.Equal(builder.Bytes(), []byte("<!--pagelist-->")) {
			buf.Write(tallyPages())
		} else {
			buf.Write(append(builder.Bytes(), byte('\n')))
		}
	}

	inmutex.Lock()
	indexCache.LastTally = time.Now()
	inmutex.Unlock()

	return buf.Bytes()
}

// generate a list of pages for the front page
func tallyPages() []byte {

	// pagelist and its associated buffer hold the links
	// displayed on the index page
	pagelist := make([]byte, 0, 1)
	buf := bytes.NewBuffer(pagelist)
	pagedir := viper.GetString("PageDir")
	viewpath := viper.GetString("ViewPath")
	reverse := viper.GetBool("ReverseTally")

	// get a list of files in the director specified
	// in the config file parameter "PageDir"
	files, err := ioutil.ReadDir(pagedir)
	if err != nil {
		return []byte("*PageDir can't be read.*")
	}

	// entry is used in the loop to construct the markdown
	// link to the given page
	if len(files) == 0 {
		return []byte("*No wiki pages! Add some content.*")
	}

	// if the config file says to reverse the page listing order
	if reverse {

		for i := len(files) - 1; i >= 0; i-- {
			f := files[i]

			// pull the page from the cache
			mutex.RLock()
			page := cachedPages[f.Name()]
			mutex.RUnlock()

			// if it hasn't been cached, cache it.
			// usually means the page is new.
			if page.Body == nil {
				page.Shortname = f.Name()
				page.Longname = pagedir + "/" + f.Name()

				err := page.cache()
				if err != nil {
					log.Printf("Couldn't pull new page %s into cache: %v\n", page.Shortname, err)
				}
			}

			// get the URI path from the file name
			// and write the formatted link to the
			// bytes.Buffer
			linkname := bytes.TrimSuffix([]byte(page.Shortname), []byte(".md"))
			buf.WriteString("* [" + page.Title + "](/" + viewpath + "/" + string(linkname) + ") " + page.Desc + " " + page.Author + "\n")
		}
	} else {

		// if the config file says to NOT reverse the page listing order
		for _, f := range files {

			// pull the page from the cache
			mutex.RLock()
			page := cachedPages[f.Name()]
			mutex.RUnlock()

			// if it hasn't been cached, cache it.
			// usually means the page is new.
			if page.Body == nil {
				page.Shortname = f.Name()
				page.Longname = pagedir + "/" + f.Name()

				err := page.cache()
				if err != nil {
					log.Printf("Couldn't pull new page %s into cache: %v\n", page.Shortname, err)
				}
			}

			// get the URI path from the file name
			// and write the formatted link to the
			// bytes.Buffer
			linkname := bytes.TrimSuffix([]byte(page.Shortname), []byte(".md"))
			buf.WriteString("* [" + page.Title + "](/" + viewpath + "/" + string(linkname) + ") " + page.Desc + " " + page.Author + "\n")
		}
	}
	buf.WriteByte(byte('\n'))
	return buf.Bytes()
}

// used when refreshing the cached copy
// of a single page
func (page *Page) cache() error {

	// buildPage() is defined in this file.
	// it reads the file and builds the Page struct
	page, err := buildPage(page.Longname)
	if err != nil {
		return err
	}

	mutex.Lock()
	cachedPages[page.Shortname] = *page
	mutex.Unlock()

	return nil

}

// compare the recorded modtime of a cached page to the
// modtime of the file. if they're different,
// return `true`, indicating the cache needs
// to be refreshed.
func (page *Page) checkCache() bool {

	newpage, err := os.Stat(page.Longname)
	if err != nil {
		log.Println("Can't stat " + page.Longname + ". Using cached copy...")
		return false
	}

	if newpage.ModTime() != page.Modtime {
		return true
	}

	return false

}

// When TildeWiki first starts, pull all available pages
// into cache, saving their modification time as well to
// determine when to re-load the page.
func genPageCache() {

	// build an array of all the (*os.FileInfo)'s
	// needed to build the cache
	pagedir := viper.GetString("PageDir")
	wikipages, err := ioutil.ReadDir(pagedir)
	if err != nil {
		log.Printf("Initial Cache Build :: Can't read directory %s\n", pagedir)
		panic(err)
	}

	// spawn a new goroutine for each entry, to cache
	// everything as quickly as possible
	for _, f := range wikipages {

		go func(f os.FileInfo) {

			page := newPage(pagedir+"/"+f.Name(), f.Name(), "", "", "", time.Time{}, nil, nil)

			err = page.cache()
			if err != nil {
				log.Println("Couldn't cache " + page.Shortname)
			}

			log.Println("Cached page " + page.Shortname)
		}(f)
	}

}
