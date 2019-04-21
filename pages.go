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

var templates = template.Must(template.ParseFiles("edit.html", "view.html"))
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
	filename = "./pages/" + filename + ".md"
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
	body := make([]byte, 0, 1)
	buf := bytes.NewBuffer(body)
	index, err := os.Open("wiki.md")
	if err != nil {
		builder := bufio.NewScanner(index)
		builder.Split(bufio.ScanLines)
		for builder.Scan() {
			if builder.Text() != "<!--#pagelist-->" {
				buf.WriteString(builder.Text())
			}
			if builder.Text() == "<!--#pagelist-->" {
				tmp := string(tallyPages())
				buf.WriteString(tmp)
			}
		}
	}
	return []byte(buf.String())
}

func tallyPages() []byte {
	pagelist := make([]byte, 0, 1)
	buf := bytes.NewBuffer(pagelist)
	files, err := ioutil.ReadDir("./pages/")
	if err != nil {
		return []byte("<strong>Pages either don't exist or can't be read.</strong>")
	}
	var title string
	var tmp string
	var shortname []byte
	for _, f := range files {
		title = getTitle(f.Name())
		shortname = []byte(f.Name())
		tmp = "<a href=\"/w/" + string(shortname[:len(shortname)-3]) + "\">" + title + "</a><br />\n"
		buf.WriteString(tmp)
	}
	return []byte(buf.String())
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
