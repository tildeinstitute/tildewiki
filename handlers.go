package main

import (
	"net/http"

	"github.com/spf13/viper"
)

func viewHandler(w http.ResponseWriter, r *http.Request, filename string) {
	p, err := loadPage(viper.GetString("PageDir") + filename)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	//renderTemplate(w, "view", p)
	w.Header().Set("Content-Type", "text/html")
	w.Write(p.Body)
}

func welcomeHandler(w http.ResponseWriter, r *http.Request) {
	parsed := render(genIndex(), viper.GetString("CSS"), viper.GetString("Name")+" "+viper.GetString("Separator")+" "+viper.GetString("ShortDesc"))
	//reader := bytes.NewReader(parsed)
	//http.ServeContent(w, r, "index.html", time.Now(), reader)
	w.Header().Set("Content-Type", "text/html")
	w.Write(parsed)
}

func editHandler(w http.ResponseWriter, r *http.Request, filename string) {
	p, err := loadPage(viper.GetString("PageDir") + "/" + filename)
	if err != nil {
		p, err := loadPage(viper.GetString("IndexDir") + viper.GetString("PageTmpl"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		p.Filename = filename
	}
	if filename != viper.GetString("Index") {
		renderTemplate(w, "edit", p)
	} else {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
}

func saveHandler(w http.ResponseWriter, r *http.Request, filename string) {
	body := r.FormValue("body")
	filename = r.FormValue("filename") + ".md"
	p := &Page{Filename: filename, Body: []byte(body)}
	err := p.save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/"+filename, http.StatusFound)
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
