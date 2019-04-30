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
	// get the file name from the request name
	filename = filename + ".md"
	// pull the page from cache
	mutex.RLock()
	page := cachedPages[filename]
	mutex.RUnlock()

	// see if it needs to be cached
	if page.checkCache() {
		err := page.cache()
		if err != nil {
			log.Printf("Client requested %s, but couldn't update cache: %v", page.Shortname, err)
		}
	}

	// if the page doesn't exist, redirect to the index
	if page.Body == nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	// send the page to the client
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, err := w.Write(page.Body)
	if err != nil {
		log.Printf("Error writing %s to HTTP stream: %v\n", filename, err)
		error500(w, r)
	}
}

// Handler for viewing the index page.
// Renders the index markdown file into HTML
// and sends it to the client.
// Calls genIndex() for each request. I need to work
// on caching the index page.
func indexHandler(w http.ResponseWriter, r *http.Request) {
	parsed := render(genIndex(), viper.GetString("CSS"), viper.GetString("Name")+" "+viper.GetString("Separator")+" "+viper.GetString("ShortDesc"))
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, err := w.Write(parsed)
	if err != nil {
		log.Printf("Error writing %s to HTTP stream: %v\n", viper.GetString("CSS"), err)
		error500(w, r)
	}
}

// Serves the favicon as a URL.
// This is due to the default behavior of
// not serving naked paths but virtual ones.
func iconHandler(w http.ResponseWriter, r *http.Request) {

	// read the raw bytes of the image
	longname := viper.GetString("AssetsDir") + "/" + viper.GetString("Icon")
	icon, err := ioutil.ReadFile(longname)
	if err != nil {
		if os.IsNotExist(err) {
			log.Printf("Favicon file specified in config does not exist: /icon request 404\n")
			error404(w, r)
		}
		log.Printf("%v\n", err)
		error500(w, r)
	}

	// check the mime type, then send
	// the bytes to the client
	w.Header().Set("Content-Type", http.DetectContentType(icon))
	_, err = w.Write(icon)
	if err != nil {
		log.Printf("Error writing favicon to HTTP stream: %v\n", err)
		error500(w, r)
	}
}

// Serves the local css file as a url.
// This is due to the default behavior of
// not serving naked paths but virtual ones.
func cssHandler(w http.ResponseWriter, r *http.Request) {

	// check if using local or remote CSS.
	// if remote, don't bother doing anything
	// and redirect requests to /
	if !cssLocal([]byte(viper.GetString("CSS"))) {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	// read the raw bytes of the stylesheet
	css, err := ioutil.ReadFile(viper.GetString("CSS"))
	if err != nil {
		if os.IsNotExist(err) {
			log.Printf("CSS file specified in config does not exist: /css request 404\n")
			error404(w, r)
		}
		log.Printf("%v\n", err)
		error500(w, r)
	}

	// send it to the client
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
			log.Printf("Invalid path requested :: %v\n", r.URL.Path)
			error404(w, r)
			return
		}
		fn(w, r, m[2])
	}
}
