package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/spf13/viper"
)

// handler for viewing content pages (not the index page)
func pageHandler(w http.ResponseWriter, r *http.Request, filename string) {
	filename = filename + ".md"
	mutex.RLock()
	page := cachedPages[filename]
	mutex.RUnlock()

	if page.checkCache() {
		err := page.cache()
		if err != nil {
			log.Printf("Client requested %s, but couldn't update cache: %v", page.Shortname, err)
		}
	}

	if page.Body == nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, err := w.Write(page.Body)
	if err != nil {
		log.Printf("Error writing %s to HTTP stream: %v\n", filename, err)
		error500(w, r)
	}
}

// handler for viewing the index page
func indexHandler(w http.ResponseWriter, r *http.Request) {
	parsed := render(genIndex(), viper.GetString("CSS"), viper.GetString("Name")+" "+viper.GetString("Separator")+" "+viper.GetString("ShortDesc"))
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, err := w.Write(parsed)
	if err != nil {
		log.Printf("Error writing %s to HTTP stream: %v\n", viper.GetString("CSS"), err)
		error500(w, r)
	}
}

// serves the icon as a URL
func iconHandler(w http.ResponseWriter, r *http.Request) {
	longname := viper.GetString("AssetsDir") + "/" + viper.GetString("Icon")
	icon, err := ioutil.ReadFile(longname)
	if err != nil {
		if os.IsNotExist(err) {
			log.Println("Favicon file specified in config does not exist: /icon request 404")
			error404(w, r)
		}
		error500(w, r)
	}
	mime := iconType(longname)
	w.Header().Set("Content-Type", mime)
	_, err = w.Write(icon)
	if err != nil {
		log.Printf("Error writing favicon to HTTP stream: %v\n", err)
		error500(w, r)
	}
}

// serves local css file as a url
func cssHandler(w http.ResponseWriter, r *http.Request) {
	if !cssLocal() {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	css, err := ioutil.ReadFile(viper.GetString("CSS"))
	if err != nil {
		if os.IsNotExist(err) {
			log.Println("CSS file specified in config does not exist: /css request 404")
			error404(w, r)
		}
		log.Println("Can't read CSS file")
		error500(w, r)
	}
	w.Header().Set("Content-Type", "text/css; charset=utf-8")
	_, err = w.Write(css)
	if err != nil {
		log.Printf("Error writing CSS file to HTTP stream: %v\n", err)
		error500(w, r)
	}
}

// closure to validate the request paths (using the regex in main.go / tildewiki.yaml)
// then pass everything on to the appropriate handler function if it all checks out
func validatePath(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			log.Println("Invalid path requested :: " + r.URL.Path)
			error404(w, r)
			return
		}
		fn(w, r, m[2])
	}
}
