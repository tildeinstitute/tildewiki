package main

import (
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

// handler for viewing content pages (not the index page)
func pageHandler(w http.ResponseWriter, r *http.Request, filename string) {
	// get the file name from the request name
	filename += ".md"
	// pull the page from cache
	pmutex.RLock()
	page := cachedPages[filename]
	pmutex.RUnlock()

	// see if it needs to be cached
	pingCache(page)

	// if the page doesn't exist, redirect to the index
	if page.Body == nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	// the etag header is used for browser-level caching.
	// sending the sha256 sum of the modtime in hexadecimal
	etag := fmt.Sprintf("%x", sha256.Sum256([]byte(page.Modtime.String())))

	// send the page to the client
	w.Header().Set("ETag", "\""+etag+"\"")
	w.Header().Set("Content-Type", htmlutf8)
	_, err := w.Write(page.Body)
	if err != nil {
		log.Printf("Error writing %s to HTTP stream: %v\n", filename, err)
		error500(w, r)
	}
}

// Handler for viewing the index page.
func indexHandler(w http.ResponseWriter, r *http.Request) {

	// check the index page's cache
	pingCache(indexCache)

	// the etag header is used for browser-level caching.
	// sending the sha256 sum of the modtime in hexadecimal
	etag := fmt.Sprintf("%x", sha256.Sum256([]byte(indexCache.Modtime.String())))

	// serve the index page
	w.Header().Set("ETag", "\""+etag+"\"")
	w.Header().Set("Content-Type", htmlutf8)
	_, err := w.Write(indexCache.Body)
	if err != nil {
		log.Printf("Error writing %s to HTTP stream: %v\n", confVars.indexFile, err)
		error500(w, r)
	}
}

// Serves the favicon as a URL.
// This is due to the default behavior of
// not serving naked paths but virtual ones.
func iconHandler(w http.ResponseWriter, r *http.Request) {

	// read the raw bytes of the image
	longname := confVars.assetsDir + "/" + confVars.iconPath
	icon, err := ioutil.ReadFile(longname)
	if err != nil {
		if os.IsNotExist(err) {
			log.Printf("Favicon file specified in config does not exist: /icon request 404\n")
			error404(w, r)
		}
		log.Printf("%v\n", err)
		error500(w, r)
	}

	// stat to get the mod time for the etag header
	stat, err := os.Stat(longname)
	if err != nil {
		log.Printf("Couldn't stat icon to send ETag header: %v\n", err)
	}

	// the etag header is used for browser-level caching.
	// sending the sha256 sum of the modtime in hexadecimal
	etag := fmt.Sprintf("%x", sha256.Sum256([]byte(stat.ModTime().String())))

	// check the mime type, then send
	// the bytes to the client
	w.Header().Set("ETag", "\""+etag+"\"")
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
	if !cssLocal([]byte(confVars.cssPath)) {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	// read the raw bytes of the stylesheet
	css, err := ioutil.ReadFile(confVars.cssPath)
	if err != nil {
		if os.IsNotExist(err) {
			log.Printf("CSS file specified in config does not exist: /css request 404\n")
			error404(w, r)
		}
		log.Printf("%v\n", err)
		error500(w, r)
	}

	// stat to get the mod time for the etag header
	stat, err := os.Stat(confVars.cssPath)
	if err != nil {
		log.Printf("Couldn't stat CSS file to send ETag header: %v\n", err)
	}

	// the etag header is used for browser-level caching.
	// sending the sha256 sum of the modtime in hexadecimal
	etag := fmt.Sprintf("%x", sha256.Sum256([]byte(stat.ModTime().String())))

	// send it to the client
	w.Header().Set("ETag", "\""+etag+"\"")
	w.Header().Set("Content-Type", cssutf8)
	_, err = w.Write(css)
	if err != nil {
		log.Printf("Error writing CSS file to HTTP stream: %v\n", err)
		error500(w, r)
	}
}

// Validate the request path, then pass everything on
// to the appropriate handler function.
func validatePath(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := confVars.validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			log.Printf("Invalid path requested :: %v\n", r.URL.Path)
			error404(w, r)
			return
		}
		fn(w, r, m[2])
	}
}

// this is a custom 500 page using a markdown doc
// in the assets directory.
// if the markdown doc can't be read, default to
// net/http's error handling
func error500(w http.ResponseWriter, _ *http.Request) {
	e500 := confVars.assetsDir + "/500.md"

	file, err := ioutil.ReadFile(e500)
	if err != nil {
		log.Printf("Tried to read 500.md: %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", htmlutf8)
	_, err = w.Write(render(file, "500: Internal Server Error"))
	if err != nil {
		log.Printf("Failed to write to HTTP stream: %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// this is a custom 404 page using a markdown doc
// in the assets directory.
// if the markdown doc can't be read, default to
// net/http's error handling
func error404(w http.ResponseWriter, r *http.Request) {
	e404 := confVars.assetsDir + "/404.md"

	file, err := ioutil.ReadFile(e404)
	if err != nil {
		log.Printf("Tried to read 404.md: %v\n", err)
		http.Error(w, err.Error(), http.StatusNotFound)
	}

	w.Header().Set("Content-Type", htmlutf8)
	_, err = w.Write(render(file, "404: Not Found"))
	if err != nil {
		log.Printf("Failed to write to HTTP stream: %v\n", err)
		error500(w, r)
	}
}
