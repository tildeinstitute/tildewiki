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
	author := getAuthor(filename)
	desc := getDesc(filename)
	parsed := render(body, viper.GetString("CSS"), title)
	return &Page{Longname: filename, Title: title, Author: author, Desc: desc, Modtime: filestat.ModTime(), Body: parsed, Raw: body}, nil
}

// scan the page for the `title: ` field
// in the header comment
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
// in the header comment
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
// in the header comment
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
	var title string
	var desc string
	var auth string
	var tmp string
	var name string
	var shortname string
	if len(files) == 0 {
		return "*No wiki pages! Add some content.*"
	}
	for _, f := range files {
		name = f.Name()
		title = getTitle(pagedir + "/" + name)
		desc = getDesc(pagedir + "/" + name)
		auth = getAuthor(pagedir + "/" + name)
		shortname = string(name[:len(name)-3])
		tmp = "* [" + title + "](/" + viewpath + "/" + shortname + ") :: " + desc + " " + auth + "\n"
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
	pagestruct.Raw = page.Raw
	pagestruct.Title = page.Title
	pagestruct.Desc = page.Desc
	pagestruct.Author = page.Author
	pagestruct.Modtime = page.Modtime
	pagestruct.Longname = page.Longname
	mutex.Lock()
	cachedPages[filename] = pagestruct
	mutex.Unlock()
}

// compare the recorded modtime of a cached page to the
// modtime of the file. if they're different,
// re-cache the page.
func checkPageCache(filename string) Page {
	mutex.RLock()
	pages := cachedPages[filename]
	mutex.RUnlock()

	newpage, err := os.Stat(pages.Longname)
	if err != nil {
		log.Println("checkPageCache() :: Can't stat " + pages.Longname + ". Using cached copy...")
		return pages
	}

	if newpage.ModTime() != pages.Modtime {
		cachePage(filename)
		log.Println("checkPageCache() :: Re-caching page " + pages.Longname)
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
