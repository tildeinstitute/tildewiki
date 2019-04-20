package main // import "github.com/gbmor/tildewiki"

import (
	"bufio"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"runtime"
	"strings"

	"github.com/russross/blackfriday"
)

var templates = template.Must(template.ParseFiles("edit.html", "view.html"))
var validPath = regexp.MustCompile("^/(edit|save|)/([a-zA-Z0-9]+)$")

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
	mdfile, err := os.Open(filename)
	var title string
	if err == nil {
		titlefinder := bufio.NewScanner(mdfile)
		for titlefinder.Scan() {
			splitter := strings.Split(titlefinder.Text(), ":")
			if splitter[0] == "title" {
				title = strings.TrimSpace(splitter[1])
				break
			}
		}
	} else {
		title = filename
	}
	parsed := blackfriday.Run(body)
	return &Page{Filename: filename, Title: title, Body: parsed}, nil
}

func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return
	}
	renderTemplate(w, "view", p)
}

func editHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title}
	}
	renderTemplate(w, "edit", p)
}

func saveHandler(w http.ResponseWriter, r *http.Request, title string) {
	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	err := p.save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/"+title, http.StatusFound)
}

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		fn(w, r, m[2])
	}
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	go func() {
		http.HandleFunc("/", makeHandler(viewHandler))
		runtime.Gosched()
	}()
	go func() {
		http.HandleFunc("/edit/", makeHandler(editHandler))
		runtime.Gosched()
	}()
	go func() {
		http.HandleFunc("/save/", makeHandler(saveHandler))
		runtime.Gosched()
	}()
	log.Fatal(http.ListenAndServe(":8080", nil))
}
