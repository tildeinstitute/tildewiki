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

// loads a given wiki page and returns a page struct pointer
func loadPage(filename string) (*Page, error) {
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Println("loadPage() :: Couldn't read " + filename)
		return nil, err
	}
	filestat, err := os.Stat(filename)
	if err != nil {
		log.Println("loadPage() :: Couldn't stat " + filename)
	}
	var shortname string
	filebyte := []byte(filename)
	for i := len(filebyte) - 1; i > 0; i-- {
		if string(filebyte[i]) == "/" {
			shortname = string(filebyte[i+1:])
		}
	}
	title := getTitle(filename)
	author := getAuthor(filename)
	desc := getDesc(filename)
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
	defer mdfile.Close()
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
	defer mdfile.Close()
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
	defer mdfile.Close()
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
	body := make([]byte, 0)
	buf := bytes.NewBuffer(body)
	index, err := os.Open(viper.GetString("IndexDir") + "/" + viper.GetString("Index"))
	if err != nil {
		return []byte("Could not open \"" + viper.GetString("IndexDir") + "/" + viper.GetString("Index") + "\"")
	}
	defer index.Close()
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
	pagelist := make([]byte, 0, 1)
	buf := bytes.NewBuffer(pagelist)
	pagedir := viper.GetString("PageDir")
	viewpath := viper.GetString("ViewPath")
	files, err := ioutil.ReadDir(pagedir)
	if err != nil {
		return "*Pages either don't exist or can't be read.*"
	}
	var entry string
	if len(files) == 0 {
		return "*No wiki pages! Add some content.*"
	}
	for _, f := range files {
		mutex.RLock()
		page := cachedPages[f.Name()]
		mutex.RUnlock()
		if page.Body == nil {
			page.Shortname = f.Name()
			page.Longname = pagedir + "/" + f.Name()
			page.cache()
		}
		linkname := []byte(page.Shortname)
		entry = "* [" + page.Title + "](/" + viewpath + "/" + string(linkname[:len(linkname)-3]) + ") :: " + page.Desc + " " + page.Author + "\n"
		buf.WriteString(entry)
	}
	return buf.String()
}

// used when refreshing the cached copy
// of a single page
func (page *Page) cache() {
	page, err := loadPage(page.Longname)
	if err != nil {
		log.Println("Page.reCache() :: Couldn't reload " + page.Longname)
	}
	mutex.Lock()
	cachedPages[page.Shortname] = *page
	mutex.Unlock()
}

// compare the recorded modtime of a cached page to the
// modtime of the file. if they're different,
// re-cache the page.
func (page *Page) checkCache() {
	newpage, err := os.Stat(page.Longname)
	if err != nil {
		log.Println("Page.checkCache() :: Can't stat " + page.Longname + ". Using cached copy...")
		return
	}
	if newpage.ModTime() != page.Modtime {
		page.cache()
		log.Println("Page.checkCache() :: Re-caching page " + page.Longname)
	}
}

// when tildewiki first starts, pull all available pages
// into cache, saving their modification time as well to
// determine when to re-load the page.
func genPageCache() {
	indexpage, err := os.Stat(viper.GetString("IndexDir") + "/" + viper.GetString("Index"))
	if err != nil {
		log.Println("genPageCache() :: Can't stat index page")
	}
	wikipages, err := ioutil.ReadDir(viper.GetString("PageDir"))
	if err != nil {
		log.Println("genPageCache() :: Can't read directory " + viper.GetString("PageDir"))
	}
	wikipages = append(wikipages, indexpage)
	var shortname string
	var longname string
	var page Page
	for _, f := range wikipages {
		shortname = f.Name()
		if shortname == viper.GetString("Index") {
			shortname = viper.GetString("IndexDir") + "/" + viper.GetString("Index")
			longname = shortname
		} else {
			longname = viper.GetString("PageDir") + "/" + f.Name()
		}
		page.Longname = longname
		page.Shortname = shortname
		page.cache()
		log.Println("genPageCache() :: Cached page " + page.Shortname)
	}
}
