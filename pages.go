package main

import (
	"bufio"
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// Loads a given wiki page and returns a page struct pointer.
// Used for building the initial cache and re-caching.
func loadPage(filename string) (*Page, error) {

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
	bytereader := bytes.NewReader(body)
	metafinder := bufio.NewScanner(bytereader)
	title := getTitle(metafinder)
	bytereader.Reset(body)
	author := getAuthor(metafinder)
	bytereader.Reset(body)
	desc := getDesc(metafinder)

	if title == "" {
		title = shortname
	}
	if desc != "" {
		desc = viper.GetString("DescSeparator") + " " + desc
	}

	// store the raw bytes of the document after parsing
	// from markdown to HTML.
	// keep the unparsed markdown for future use (maybe gopher?)
	return &Page{
		Longname:  filename,
		Shortname: shortname,
		Title:     title,
		Author:    author,
		Desc:      desc,
		Modtime:   stat.ModTime(),
		Body:      render(body, viper.GetString("CSS"), title),
		Raw:       body}, nil
}

// scan the page for the `title: ` field
// in the header comment. used in the construction
// of the page cache on startup
func getTitle(metafinder *bufio.Scanner) string {

	// scan the file line by line until it finds
	// the title: comment, return the value.
	//titlefinder := bufio.NewScanner(metafinder)
	for metafinder.Scan() {
		splitter := bytes.Split(metafinder.Bytes(), []byte(":"))
		if bytes.Equal(bytes.ToLower(splitter[0]), []byte("title")) {
			return string(bytes.TrimSpace(splitter[1]))
		}
	}

	return ""
}

// scan the page for the `description: ` field
// in the header comment. used in the construction
// of the page cache on startup
func getDesc(metafinder *bufio.Scanner) string {

	// scan the file line by line until it finds
	// the description: comment, return the value.
	//descfinder := bufio.NewScanner(metafinder)
	for metafinder.Scan() {
		splitter := bytes.Split(metafinder.Bytes(), []byte(":"))
		if bytes.Equal(bytes.ToLower(splitter[0]), []byte("description")) {
			return string(bytes.TrimSpace(splitter[1]))
		}
	}

	return ""
}

// scan the page for the `author: ` field
// in the header comment. used in the construction
// of the page cache on startup
func getAuthor(metafinder *bufio.Scanner) string {

	// scan the file line by line until it finds
	// the author: comment, return the value.
	//authfinder := bufio.NewScanner(metafinder)
	for metafinder.Scan() {
		splitter := bytes.Split(metafinder.Bytes(), []byte(":"))
		if bytes.Equal(bytes.ToLower(splitter[0]), []byte("author")) {
			return "`by " + string(bytes.TrimSpace(splitter[1])) + "`"
		}
	}

	return ""
}

// generate the front page of the wiki
func genIndex() []byte {
	indexpath := viper.GetString("AssetsDir") + "/" + viper.GetString("Index")

	// body holds the bytes of the generated index page being sent to the client.
	// create the byte array and the buffer used to write to it
	body := make([]byte, 0)
	buf := bytes.NewBuffer(body)

	index, err := ioutil.ReadFile(indexpath)
	if err != nil {
		return []byte("Could not open \"" + indexpath + "\"")
	}

	// scan the file line by line until it finds the anchor
	// comment. replace the anchor comment with a list of
	// wiki pages sorted alphabetically by title.
	builder := bufio.NewScanner(bytes.NewReader(index))
	builder.Split(bufio.ScanLines)
	for builder.Scan() {
		if bytes.Equal(builder.Bytes(), []byte("<!--pagelist-->")) {
			buf.Write(tallyPages())
		} else {
			buf.Write(append(builder.Bytes(), byte('\n')))
		}
	}

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
	if reverse {
		for i := len(files) - 1; i >= 0; i-- {
			f := files[i]
			mutex.RLock()
			page := cachedPages[f.Name()]
			mutex.RUnlock()

			if page.Body == nil {
				page.Shortname = f.Name()
				page.Longname = pagedir + "/" + f.Name()

				err := page.cache()
				if err != nil {
					log.Printf("Couldn't pull new page %s into cache: %v\n", page.Shortname, err)
				}
			}

			linkname := bytes.TrimSuffix([]byte(page.Shortname), []byte(".md"))
			buf.WriteString("* [" + page.Title + "](/" + viewpath + "/" + string(linkname) + ") " + page.Desc + " " + page.Author + "\n")
		}
	} else {

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
	indexpath := viper.GetString("AssetsDir") + "/" + viper.GetString("Index")
	indexname := viper.GetString("Index")
	pagedir := viper.GetString("PageDir")

	// build an array of all the (*os.FileInfo)'s
	// needed to build the cache
	indexpage, err := os.Stat(indexpath)
	if err != nil {
		log.Printf("Initial Cache Build :: Can't stat index page: %v\n", err)
	}
	wikipages, err := ioutil.ReadDir(pagedir)
	if err != nil {
		log.Printf("Initial Cache Build :: Can't read directory %s: %v\n", pagedir, err)
	}

	wikipages = append(wikipages, indexpage)

	// spawn a new goroutine for each entry, to cache
	// everything as quickly as possible
	for _, f := range wikipages {
		go func(f os.FileInfo) {

			var page Page
			page.Shortname = f.Name()

			// store any page with the same name as
			// the index page as its relative path
			// for the key.
			// this is to try to avoid collisions
			// by explicitly disallowing pages with
			// the same filename as the index.
			// later I'll cache the assets separately
			// but this works for now.
			if page.Shortname == indexname {
				page.Shortname = indexpath
				page.Longname = page.Shortname
			} else {
				page.Longname = pagedir + "/" + page.Shortname
			}

			err = page.cache()
			if err != nil {
				log.Println("Couldn't cache " + page.Shortname)
			}

			log.Println("Cached page " + page.Shortname)
		}(f)
	}

}
