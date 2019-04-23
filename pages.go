package main

import (
	"bufio"
	"bytes"
	"io/ioutil"
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

func (p *Page) save() error {
	return ioutil.WriteFile(p.Filename, p.Body, 0600)
}

func loadPage(filename string) (*Page, error) {
	filename = filename + ".md"
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	title := getTitle(filename)
	parsed := render(body, viper.GetString("CSS"), title)
	return &Page{Filename: filename, Title: title, Body: parsed}, nil
}

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

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page, r *http.Request) {
	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		error500(w, r)
		return
	}
}
