package main

import (
	"io/ioutil"
	"log"
	"net/http"

	"github.com/spf13/viper"
)

// handler for viewing content pages (not the index page)
func viewHandler(w http.ResponseWriter, r *http.Request, filename string) {
	filename = filename + ".md"
	page := checkPageCache(filename)

	if page.Body == nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	_, err := w.Write(page.Body)
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

// serves the icon as a URL
func iconHandler(w http.ResponseWriter, r *http.Request) {
	longname := viper.GetString("IndexDir") + "/" + viper.GetString("Icon")
	icon, err := ioutil.ReadFile(longname)
	if err != nil {
		log.Println("iconHandler() :: Couldn't read icon file")
		_, err = w.Write(nil)
		if err != nil {
			log.Println("iconHandler() :: 500 :: Couldn't write to http stream")
			error500(w, r)
		}
	}
	mime := iconType(longname)
	w.Header().Set("Content-Type", mime)
	_, err = w.Write(icon)
	if err != nil {
		error500(w, r)
	}
}

// serves local css file as a url
func cssHandler(w http.ResponseWriter, r *http.Request) {
	if !cssLocal() {
		log.Println("cssHandler() :: Using remote CSS, nothing to serve...")
		http.Redirect(w, r, "/", 302)
		return
	}
	css, err := ioutil.ReadFile(viper.GetString("CSS"))
	if err != nil {
		log.Println("cssHandler() :: Can't read CSS file")
		http.Redirect(w, r, "/", 302)
	}
	w.Header().Set("Content-Type", "text/css")
	_, err = w.Write(css)
	if err != nil {
		error500(w, r)
	}
}

// closure to validate the request paths (using the regex in main.go / tildewiki.yaml)
// then pass everything on to the appropriate handler function if it all checks out
func validatePath(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			error404(w, r)
			return
		}
		fn(w, r, m[2])
	}
}
