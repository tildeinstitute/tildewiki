package main

import (
	"bufio"
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/spf13/viper"
)

// Loads a given wiki page and returns a page struct pointer.
// Used for building the initial cache and re-caching.
func loadPage(filename string) (*Page, error) {

	// read the raw bytes
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Println("Couldn't read " + filename)
		return nil, err
	}

	// stat the file to get mod time later
	filestat, err := os.Stat(filename)
	if err != nil {
		log.Println("Couldn't stat " + filename)
	}

	// extract the file name from the path
	var shortname string
	filebyte := []byte(filename)
	for i := len(filebyte) - 1; i > 0; i-- {
		if filebyte[i] == byte('/') {
			shortname = string(filebyte[i+1:])
			break
		}
	}

	// get meta info on file from the header comment
	title := getTitle(filename)
	author := getAuthor(filename)
	desc := getDesc(filename)

	// store the raw bytes of the document after parsing
	// from markdown to HTML.
	// keep the unparsed markdown for future use (maybe gopher?)
	parsed := render(body, viper.GetString("CSS"), title)
	return &Page{
		Longname:  filename,
		Shortname: shortname,
		Title:     title,
		Author:    author,
		Desc:      desc,
		Modtime:   filestat.ModTime(),
		Body:      parsed,
		Raw:       body}, nil
}

// scan the page for the `title: ` field
// in the header comment. used in the construction
// of the page cache on startup
func getTitle(filename string) string {
	mdfile, err := os.Open(filename)
	if err != nil {
		return filename
	}

	// defer closing and checking of the error returned from (*os.File).Close()
	defer func() {
		err := mdfile.Close()
		if err != nil {
			log.Printf("Deferred closing of %s resulted in error: %v\n", filename, err)
		}
	}()

	// scan the file line by line until it finds
	// the title: comment, return the value.
	titlefinder := bufio.NewScanner(mdfile)
	for titlefinder.Scan() {
		splitter := strings.Split(titlefinder.Text(), ":")
		if strings.ToLower(splitter[0]) == "title" {
			return strings.TrimSpace(splitter[1])
		}
	}
	return filename
}

// scan the page for the `description: ` field
// in the header comment. used in the construction
// of the page cache on startup
func getDesc(filename string) string {
	mdfile, err := os.Open(filename)
	if err != nil {
		return ""
	}

	// defer closing and checking of the error returned from (*os.File).Close()
	defer func() {
		err := mdfile.Close()
		if err != nil {
			log.Printf("Deferred closing of %s resulted in error: %v\n", filename, err)
		}
	}()

	// scan the file line by line until it finds
	// the description: comment, return the value.
	descfinder := bufio.NewScanner(mdfile)
	for descfinder.Scan() {
		splitter := strings.Split(descfinder.Text(), ":")
		if strings.ToLower(splitter[0]) == "description" {
			return strings.TrimSpace(splitter[1])
		}
	}
	return ""
}

// scan the page for the `author: ` field
// in the header comment. used in the construction
// of the page cache on startup
func getAuthor(filename string) string {
	mdfile, err := os.Open(filename)
	if err != nil {
		return ""
	}

	// defer closing and checking of the error returned from (*os.File).Close()
	defer func() {
		err := mdfile.Close()
		if err != nil {
			log.Printf("Deferred closing of %s resulted in error: %v\n", filename, err)
		}
	}()

	// scan the file line by line until it finds
	// the author: comment, return the value.
	authfinder := bufio.NewScanner(mdfile)
	for authfinder.Scan() {
		splitter := strings.Split(authfinder.Text(), ":")
		if strings.ToLower(splitter[0]) == "author" {
			return "`by " + strings.TrimSpace(splitter[1]) + "`"
		}
	}
	return ""
}

// generate the front page of the wiki
func genIndex() []byte {

	// body holds the bytes of the generated index page being sent to the client.
	// create the byte array and the buffer used to write to it
	body := make([]byte, 0)
	buf := bytes.NewBuffer(body)
	index, err := os.Open(viper.GetString("AssetsDir") + "/" + viper.GetString("Index"))
	if err != nil {
		return []byte("Could not open \"" + viper.GetString("AssetsDir") + "/" + viper.GetString("Index") + "\"")
	}

	// defer closing and checking of the error returned from (*os.File).Close()
	defer func() {
		err := index.Close()
		if err != nil {
			log.Printf("Deferred closing of %s resulted in error: %v\n", viper.GetString("Index"), err)
		}
	}()

	// scan the file line by line until it finds the anchor
	// comment. replace the anchor comment with a list of
	// wiki pages sorted alphabetically by title.
	builder := bufio.NewScanner(index)
	builder.Split(bufio.ScanLines)
	for builder.Scan() {
		if builder.Text() == "<!--pagelist-->" {
			tmp := tallyPages()
			buf.WriteString(tmp + "\n")
		} else {
			buf.WriteString(builder.Text() + "\n")
		}
	}
	return buf.Bytes()
}

// generate a list of pages for the front page
func tallyPages() string {
	// pagelist and its associated buffer hold the links
	// displayed on the index page
	pagelist := make([]byte, 0, 1)
	buf := bytes.NewBuffer(pagelist)
	pagedir := viper.GetString("PageDir")
	viewpath := viper.GetString("ViewPath")

	// get a list of files in the director specified
	// in the config file parameter "PageDir"
	files, err := ioutil.ReadDir(pagedir)
	if err != nil {
		return "*PageDir can't be read.*"
	}
	// entry is used in the loop to construct the markdown
	// link to the given page
	var entry string
	if len(files) == 0 {
		return "*No wiki pages! Add some content.*"
	}
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
		entry = "* [" + page.Title + "](/" + viewpath + "/" + string(linkname) + ") :: " + page.Desc + " " + page.Author + "\n"
		buf.WriteString(entry)
	}
	return buf.String()
}

// used when refreshing the cached copy
// of a single page
func (page *Page) cache() error {
	// loadPage() is defined in this file.
	// it reads the file and builds the Page struct
	page, err := loadPage(page.Longname)
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
	indexpage, err := os.Stat(viper.GetString("AssetsDir") + "/" + viper.GetString("Index"))
	if err != nil {
		log.Println("Initial Cache Build :: Can't stat index page")
	}
	wikipages, err := ioutil.ReadDir(viper.GetString("PageDir"))
	if err != nil {
		log.Println("Initial Cache Build :: Can't read directory " + viper.GetString("PageDir"))
	}
	wikipages = append(wikipages, indexpage)

	// spawn a new goroutine for each entry, to cache
	// everything as quickly as possible
	for _, f := range wikipages {
		go func(f os.FileInfo) {
			var page Page
			shortname := f.Name()
			var longname string
			// store any page with the same name as
			// the index page as its relative path
			// for the key.
			// this is to try to avoid collisions
			// by explicitly disallowing pages with
			// the same filename as the index
			// later I'll cache the assets separately
			// but this works for now.
			if shortname == viper.GetString("Index") {
				shortname = viper.GetString("AssetsDir") + "/" + viper.GetString("Index")
				longname = shortname
			} else {
				longname = viper.GetString("PageDir") + "/" + f.Name()
			}
			page.Longname = longname
			page.Shortname = shortname
			err = page.cache()
			if err != nil {
				log.Println("Couldn't cache " + page.Shortname)
			}
			log.Println("Cached page " + page.Shortname)
		}(f)
	}
}
