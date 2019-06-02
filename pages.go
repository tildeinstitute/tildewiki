package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/spf13/viper"
)

// Loads a given wiki page and returns a page object.
// Used for building the initial cache and re-caching.
func buildPage(filename string) (*Page, error) {

	// open the page into *os.File
	file, err := os.Open(filename)
	if err != nil {
		log.Printf("%v\n", err.Error())
		return nil, err
	}

	// the cleanup crew
	defer func() {
		err = file.Close()
		if err != nil {
			log.Printf("%v\n", err.Error())
		}
	}()

	// stat the file to get mod time later
	stat, err := file.Stat()
	if err != nil {
		log.Printf("Couldn't stat %s: %v\n", filename, err.Error())
	}

	// body holds the raw bytes from the file
	var body pagedata
	body, err = ioutil.ReadAll(file)
	if err != nil {
		log.Printf("%v\n", err.Error())
	}

	// extract the file name from the path
	_, shortname := filepath.Split(filename)

	// get meta info on file from the header comment
	title, desc, author := body.getMeta()
	if title == "" {
		title = shortname
	}
	if desc != "" {
		confVars.mu.RLock()
		desc = confVars.descSep + " " + desc
		confVars.mu.RUnlock()
	}
	if author != "" {
		author = "`by " + author + "`"
	}

	// longtitle is used in the <title> tags of the output html
	confVars.mu.RLock()
	longtitle := title + " " + confVars.titleSep + " " + confVars.wikiName
	confVars.mu.RUnlock()

	// store the raw bytes of the document after parsing
	// from markdown to HTML.
	// keep the unparsed markdown for future use (maybe gopher?)
	bodydata := render(body, longtitle)
	return newPage(filename, shortname, title, author, desc, stat.ModTime(), bodydata, body, false), nil
}

// Scan the page until reaching following fields in the
// header comment:
//		title:
//		author:
//		description:
func (body pagedata) getMeta() (string, string, string) {

	// a bit redundant, but scanner is simpler to use
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
			break
		}

	}

	return title, desc, author
}

// Checks the index page's cache. Returns true if the
// index needs to be re-cached.
// This method helps satisfy the cacher interface.
func (indexCache *indexPage) checkCache() bool {

	// if the last tally time is past the
	// interval in the config file, re-cache
	if interval, err := time.ParseDuration(viper.GetString("IndexRefreshInterval")); err == nil {
		imutex.RLock()
		if time.Since(indexCache.LastTally) > interval {
			imutex.RUnlock()
			return true
		}
		imutex.RUnlock()
	} else {
		log.Printf("Couldn't parse index refresh interval: %v\n", err.Error())
	}

	// if the stored mod time is different
	// from the file's modtime, re-cache
	confVars.mu.RLock()
	if stat, err := os.Stat(confVars.assetsDir + "/" + confVars.indexFile); err == nil {
		imutex.RLock()
		if stat.ModTime() != indexCache.Modtime {
			imutex.RUnlock()
			confVars.mu.RUnlock()
			return true
		}
		imutex.RUnlock()
	} else {
		log.Printf("Couldn't stat index page: %v\n", err.Error())
	}
	confVars.mu.RUnlock()

	// if the last tally time or stored mod time is zero, signal
	// to re-cache the index
	imutex.RLock()
	if indexCache.LastTally.IsZero() || indexCache.Modtime.IsZero() {
		imutex.RUnlock()
		return true
	}
	imutex.RUnlock()

	return false
}

// Re-caches the index page.
// This method helps satisfy the cacher interface.
func (indexCache *indexPage) cache() error {
	confVars.mu.RLock()
	body := render(genIndex(), confVars.wikiName+" "+confVars.titleSep+" "+confVars.wikiDesc)
	confVars.mu.RUnlock()
	if body == nil {
		return errors.New("indexPage.cache(): getting nil bytes")
	}
	imutex.Lock()
	indexCache.Body = body
	imutex.Unlock()
	return nil
}

// Generate the front page of the wiki
func genIndex() []byte {

	var err error
	confVars.mu.RLock()
	indexpath := confVars.assetsDir + "/" + confVars.indexFile
	confVars.mu.RUnlock()

	// stat to check mod time
	stat, err := os.Stat(indexpath)
	if err != nil {
		log.Printf("Couldn't stat index: %v\n", err.Error())
	}

	// if the index file has been modified,
	// vaccuum up those bytes into the cache

	imutex.RLock()
	if indexCache.Modtime != stat.ModTime() {
		imutex.RUnlock()
		imutex.Lock()
		indexCache.Raw, err = ioutil.ReadFile(indexpath)
		imutex.Unlock()
		if err != nil {
			return []byte("Could not open \"" + indexpath + "\"")
		}

	} else {
		imutex.RUnlock()
	}

	// body holds the bytes of the generated index page being sent to the client.
	// create the byte array and the buffer used to write to it
	body := make([]byte, 0)
	buf := bytes.NewBuffer(body)

	// scan the file line by line until it finds the anchor
	// comment. replace the anchor comment with a list of
	// wiki pages sorted alphabetically by title.
	imutex.RLock()
	builder := bufio.NewScanner(bytes.NewReader(indexCache.Raw))
	imutex.RUnlock()
	builder.Split(bufio.ScanLines)

	for builder.Scan() {
		if bytes.Equal(builder.Bytes(), []byte("<!--pagelist-->")) {
			tallyPages(buf)
		} else {
			n, err := buf.Write(append(builder.Bytes(), byte('\n')))
			if err != nil || n == 0 {
				log.Printf("Error writing to buffer: %v\n", err.Error())
			}
		}
	}

	// the LastTally field lets us know
	// when the index was last generated
	// by this function.
	imutex.Lock()
	indexCache.LastTally = time.Now()
	imutex.Unlock()

	return buf.Bytes()
}

// Generate a list of pages for the index.
// Called by genIndex() when the anchor
// comment has been found.
func tallyPages(buf *bytes.Buffer) {

	// get a list of files in the directory specified
	// in the config file parameter "PageDir"
	confVars.mu.RLock()
	if files, err := ioutil.ReadDir(confVars.pageDir); err == nil {

		// entry is used in the loop to construct the markdown
		// link to the given page
		if len(files) == 0 {
			n, err := buf.WriteString("*No wiki pages! Add some content.*\n")
			if err != nil || n == 0 {
				log.Printf("Error writing to buffer: %v\n", err.Error())
			}
			confVars.mu.RUnlock()
			return
		}

		// true if reversing page order, otherwise don't reverse
		if confVars.reverseTally {
			for i := len(files) - 1; i >= 0; i-- {
				writeIndexLinks(files[i], buf)
			}
		} else {
			for _, f := range files {
				writeIndexLinks(f, buf)
			}
		}
	} else {
		n, err := buf.WriteString("*PageDir can't be read.*\n")
		if err != nil || n == 0 {
			log.Printf("Error writing to buffer: %v\n", err.Error())
		}
	}

	err := buf.WriteByte(byte('\n'))
	if err != nil {
		log.Printf("Error writing to buffer: %v\n", err.Error())
	}
	confVars.mu.RUnlock()
}

// Takes in a file and outputs a markdown link to it.
// Called by tallyPages() for each file in the pages
// directory.
func writeIndexLinks(f os.FileInfo, buf *bytes.Buffer) {
	var page *Page
	var err error
	if _, exists := cachedPages[f.Name()]; exists {
		// pull the page from the cache
		page, err = pullFromCache(f.Name())
		if err != nil {
			log.Printf("%v\n", err.Error())
		}
	} else {
		// if it hasn't been cached, cache it.
		// usually means the page is new.
		confVars.mu.RLock()
		newpage := newBarePage(confVars.pageDir+"/"+f.Name(), f.Name())
		confVars.mu.RUnlock()
		if err := newpage.cache(); err != nil {
			log.Printf("While caching page %v during the index generation, caught an error: %v\n", f.Name(), err.Error())
		}
		page, err = pullFromCache(f.Name())
		if err != nil {
			log.Printf("%v\n", err.Error())
		}
	}
	// get the URI path from the file name
	// and write the formatted link to the
	// bytes.Buffer
	linkname := bytes.TrimSuffix([]byte(page.Shortname), []byte(".md"))
	confVars.mu.RLock()
	n, err := buf.WriteString("* [" + page.Title + "](" + confVars.viewPath + string(linkname) + ") " + page.Desc + " " + page.Author + "\n")
	confVars.mu.RUnlock()
	if err != nil || n == 0 {
		log.Printf("Error writing to buffer: %v\n", err.Error())
	}
}

// Caches a page.
// This method helps satisfy the cacher interface.
func (page *Page) cache() error {

	// If buildPage() successfully returns a page
	// object ptr, then push it into the cache
	if newpage, err := buildPage(page.Longname); err == nil {
		pmutex.Lock()
		cachedPages[newpage.Shortname] = newpage
		pmutex.Unlock()
	} else {
		log.Printf("Couldn't cache %v: %v", page.Longname, err.Error())
		return err
	}
	return nil
}

// Compare the recorded modtime of a cached page to the
// modtime of the file on disk. If they're different,
// return `true`, indicating the cache needs
// to be refreshed. Also returns `true` if the
// page.Recache field is set to `true`.
// This method helps satisfy the cacher interface.
func (page *Page) checkCache() bool {

	if newpage, err := os.Stat(page.Longname); err == nil {
		if newpage.ModTime() != page.Modtime || page.Recache {
			return true
		}
	} else {
		log.Println("Can't stat " + page.Longname + ". Using cached copy...")
	}

	return false
}

// When TildeWiki first starts, pull all available pages
// into cache, saving their modification time as well to
// detect changes to a page.
func genPageCache() {

	// spawn a new goroutine for each entry, to cache
	// everything as quickly as possible
	confVars.mu.RLock()
	if wikipages, err := ioutil.ReadDir(confVars.pageDir); err == nil {
		var wg sync.WaitGroup
		for _, f := range wikipages {

			wg.Add(1)

			go func(f os.FileInfo) {
				confVars.mu.RLock()
				page := newBarePage(confVars.pageDir+"/"+f.Name(), f.Name())
				confVars.mu.RUnlock()
				if err := page.cache(); err != nil {
					log.Printf("While generating initial cache, caught error for %v: %v\n", f.Name(), err.Error())
				}
				log.Printf("Cached page %v\n", page.Shortname)

				wg.Done()
			}(f)
		}

		wg.Wait()

	} else {
		log.Printf("Initial cache build :: Can't read directory: %s\n", err.Error())
		log.Printf("**NOTICE** TildeWiki's cache may not function correctly until this is resolved.\n")
		log.Printf("\tPlease verify the directory in tildewiki.yml is correct and restart TildeWiki\n")
	}
	confVars.mu.RUnlock()
}

// Wrapper function to check the cache
// of any cacher type, and if true,
// re-cache the data
func pingCache(c cacher) {

	if c.checkCache() {
		if err := c.cache(); err != nil {
			log.Printf("Pinged cache, received error while caching: %v\n", err.Error())
		}
	}
}

// Pulling from cache is its own function.
// Less worrying about mutexes.
func pullFromCache(filename string) (*Page, error) {

	pmutex.RLock()
	if page, ok := cachedPages[filename]; ok {
		pmutex.RUnlock()
		return page, nil
	}
	pmutex.RUnlock()

	return nil, fmt.Errorf("error pulling %v from cache", filename)
}

// Blanks stored modtimes for the page cache.
// Used to trigger a forced re-cache on the
// next page load.
func triggerRecache() {
	for _, v := range cachedPages {
		v.Recache = true
	}
}
