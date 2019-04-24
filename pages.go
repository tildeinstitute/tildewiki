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
	title := getTitle(filename)
	parsed := render(body, viper.GetString("CSS"), title)
	return &Page{Filename: filename, Title: title, Modtime: filestat.ModTime(), Body: parsed}, nil
}

// scan the page for the `title: ` field
// in the header comment
func getTitle(filename string) string {
	mdfile, err := os.Open(filename)
	if err != nil {
		return filename
	}
	titlefinder := bufio.NewScanner(mdfile)
	for titlefinder.Scan() {
		splitter := strings.Split(titlefinder.Text(), ":")
		if splitter[0] == "title" {
			return strings.TrimSpace(splitter[1])
		}
	}
	return filename
}

// generate the front page of the wiki
func genIndex() []byte {
	body := make([]byte, 0)
	buf := bytes.NewBuffer(body)
	index, err := os.Open(viper.GetString("IndexDir") + "/" + viper.GetString("Index"))
	if err != nil {
		return []byte("Could not open \"" + viper.GetString("IndexDir") + "/" + viper.GetString("Index") + "\"")
	}
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
	var title string
	var tmp string
	var name string
	var shortname string
	if len(files) == 0 {
		return "*No wiki pages! Add some content.*"
	}
	for _, f := range files {
		title = getTitle(pagedir + "/" + f.Name())
		name = f.Name()
		shortname = string(name[:len(name)-3])
		tmp = "* [" + title + "](/" + viewpath + "/" + shortname + ")\n"
		buf.WriteString(tmp)
	}
	return buf.String()
}

// Pull a page into memory
func cachePage(filename string) {
	var longname string
	if filename != viper.GetString("IndexDir")+"/"+viper.GetString("Index") {
		longname = viper.GetString("PageDir") + "/" + filename
	} else {
		longname = filename
	}

	page, err := loadPage(longname)
	if err != nil {
		log.Println("cachePage() :: Can't cache " + filename)
		return
	}
	if page == nil {
		panic("cachePage() :: Call to loadPage() returned nil")
	}
	var pagestruct Page
	pagestruct.Body = page.Body
	pagestruct.Title = page.Title
	pagestruct.Modtime = page.Modtime
	mutex.Lock()
	cachedPages[filename] = pagestruct
	mutex.Unlock()
}

// compare the size and timestamp of a cached page.
// if the size is different or the cached version is
// old, then reload the page into memory
func checkPageCache(filename string) Page {
	longname := viper.GetString("PageDir") + "/" + filename
	if filename == viper.GetString("Index") {
		longname = viper.GetString("IndexDir") + "/" + filename
		filename = longname
	}

	mutex.RLock()
	pages := cachedPages[filename]
	mutex.RUnlock()

	newpage, err := os.Stat(longname)
	if err != nil {
		log.Println("checkPageCache() :: Can't stat " + filename + ". Using cached copy...")
		return pages
	}

	if newpage.ModTime() != pages.Modtime {
		cachePage(filename)
		log.Println("checkPageCache() :: Re-caching page " + longname)
	}
	return cachedPages[filename]
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
	var tmp string
	for _, f := range wikipages {
		tmp = f.Name()
		if tmp == viper.GetString("Index") {
			tmp = viper.GetString("IndexDir") + "/" + viper.GetString("Index")
		}
		cachePage(tmp)
		log.Println("genPageCache() :: Cached page " + tmp)
	}
}
