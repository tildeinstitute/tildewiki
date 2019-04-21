package main // import "github.com/gbmor/tildewiki"

import (
	"bufio"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/russross/blackfriday"
)

var templates = template.Must(template.ParseFiles("edit.html", "view.html"))
var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")

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
	if filename == "" {
		filename = "wiki.md"
	} else {
		filename += ".md"
	}
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

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
