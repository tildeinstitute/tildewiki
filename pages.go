package main // import "github.com/gbmor/tildewiki"

import (
	"bufio"
	"bytes"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/russross/blackfriday"
)

var templates = template.Must(template.ParseFiles("templates/edit.html", "templates/view.html"))
var validPath = regexp.MustCompile("^/(edit|save|w)/([a-zA-Z0-9]+)$")

// Page holds wiki page title and body
type Page struct {
	Filename string
	Title    string
	Body     []byte
}

func (p *Page) save() error {
	return ioutil.WriteFile(p.Filename, p.Body, 0600)
}

func loadPage(filename string) (*Page, error) {
	filename = "pages/" + filename + ".md"
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	title := getTitle(filename)
	parsed := blackfriday.Run(body)
	return &Page{Filename: filename, Title: title, Body: parsed}, nil
}

func getTitle(filename string) string {
	mdfile, err := os.Open(filename)
	if err == nil {
		titlefinder := bufio.NewScanner(mdfile)
		for titlefinder.Scan() {
			splitter := strings.Split(titlefinder.Text(), ":")
			if splitter[0] == "title" {
				return strings.TrimSpace(splitter[1])
			}
		}
	}
	return filename
}

func genIndex() []byte {
	body := make([]byte, 0)
	buf := bytes.NewBuffer(body)
	index, err := os.Open("wiki.md")
	if err != nil {
		return []byte("Could not open \"wiki.md\"")
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

func tallyPages() string {
	pagelist := make([]byte, 0, 1)
	buf := bytes.NewBuffer(pagelist)
	files, err := ioutil.ReadDir("./pages/")
	if err != nil {
		return "*Pages either don't exist or can't be read.*"
	}
	var title string
	var tmp string
	var name string
	var shortname string
	for _, f := range files {
		title = getTitle(f.Name())
		name = string(f.Name())
		shortname = string(name[:len(name)-3])
		tmp = "* [" + title + "](/w/" + shortname + ")\n"
		buf.WriteString(tmp)
	}
	return buf.String()
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
