package main

import (
	"bufio"
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/spf13/viper"
)

// Page holds wiki page title and body
type Page struct {
	Filename string
	Title    string
	Body     []byte
}

// method to save a page after editing
func (p *Page) save() error {
	return ioutil.WriteFile(p.Filename, p.Body, 0600)
}

// loads a given wiki page and returns a page struct pointer
func loadPage(filename string) (*Page, error) {
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	title := getTitle(filename)
	parsed := render(body, viper.GetString("CSS"), title)
	return &Page{Filename: filename, Title: title, Body: parsed}, nil
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
	builder := bufio.NewScanner(index)
	builder.Split(bufio.ScanLines)
	for builder.Scan() {
		if builder.Text() == "<!--pagelist-->" {
			tmp := tallyPages()
			buf.WriteString(tmp + "\n")
		} else if builder.Text() != "<!--pagelist-->" {
			buf.WriteString(builder.Text() + "\n")
		} else {
			// schrodinger's HTML
			buf.WriteString(builder.Text() + "\n")
		}
	}
	return []byte(buf.String())
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

// pass a page to the parsed HTML template
func renderTemplate(w http.ResponseWriter, tmpl string, p *Page, r *http.Request) {
	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		error500(w, r)
		return
	}
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
	cachedPages[filename] = page.Body
}

// compare the size and timestamp of a cached page.
// if the size is different or the cached version is
// old, then reload the page into memory
func checkPageCache(filename string) []byte {
	longname := viper.GetString("PageDir") + "/" + filename
	if filename == viper.GetString("Index") {
		longname = viper.GetString("IndexDir") + "/" + filename
		filename = longname
	}
	oldPageSize := int64(len(cachedPages[filename]))
	newpage, err := os.Stat(longname)
	if err != nil {
		log.Println("checkPageCache() :: Can't stat " + filename)
		return cachedPages[filename]
	}

	if oldPageSize != newpage.Size() {
		cachePage(filename)
		log.Println("checkPageCache() :: Re-caching page " + longname)
	}
	return cachedPages[filename]
}

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
