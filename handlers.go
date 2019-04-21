package main // import "github.com/gbmor/tildewiki"

import "net/http"

func viewHandler(w http.ResponseWriter, r *http.Request, filename string) {
	p, err := loadPage(filename)
	if err != nil && p.Filename != "wiki.md" {
		http.Redirect(w, r, "/edit/"+filename, http.StatusFound)
		return
	}
	if err != nil && p.Filename == "wiki.md" {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	if p.Filename == "wiki.md" {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	renderTemplate(w, "view", p)
}

func welcomeHandler(w http.ResponseWriter, r *http.Request) {
	return //placeholder
}

func editHandler(w http.ResponseWriter, r *http.Request, filename string) {
	p, err := loadPage(filename)
	if err != nil {
		p, err := loadPage("template")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		p.Filename = filename
	}
	if filename != "wiki.md" {
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
