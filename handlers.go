package main

import (
	"net/http"

	"github.com/spf13/viper"
)

// handler for viewing content pages (not the index page)
func viewHandler(w http.ResponseWriter, r *http.Request, filename string) {
	filename = filename + ".md"
	page := checkPageCache(filename)

	if page == nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	_, err := w.Write(page)
	if err != nil {
		error500(w, r)
	}
}

// handler for viewing the index page
func welcomeHandler(w http.ResponseWriter, r *http.Request) {
	parsed := render(genIndex(), viper.GetString("CSS"), viper.GetString("Name")+" "+viper.GetString("Separator")+" "+viper.GetString("ShortDesc"))
	w.Header().Set("Content-Type", "text/html")
	_, err := w.Write(parsed)
	if err != nil {
		error500(w, r)
	}
}

// handler for requests to edit a page
/* func editHandler(w http.ResponseWriter, r *http.Request, filename string) {
	p, err := loadPage(viper.GetString("PageDir") + "/" + filename)
	if err != nil {
		p, err := loadPage(viper.GetString("IndexDir") + viper.GetString("PageTmpl"))
		if err != nil {
			error500(w, r)
			return
		}
		p.Filename = filename
	}
	if filename != viper.GetString("Index") {
		renderTemplate(w, "edit", p, r)
	} else {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
}

// saves a page after editing
func saveHandler(w http.ResponseWriter, r *http.Request, filename string) {
	body := r.FormValue("body")
	filename = r.FormValue("filename") + ".md"
	p := &Page{Filename: filename, Body: []byte(body)}
	err := p.save()
	if err != nil {
		error500(w, r)
		return
	}
	http.Redirect(w, r, "/"+filename, http.StatusFound)
}

// closure to validate the request paths (using the regex in main.go / tildewiki.yaml)
// then pass everything on to the appropriate handler function if it all checks out
func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			error404(w, r)
			return
		}
		fn(w, r, m[2])
	}
}*/
