package main

import (
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

// handler for viewing content pages (not the index page)
func pageHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	filename := vars["pageReq"]
	filename += ".md"

	page, err := pullFromCache(filename)
	if err != nil {
		log.Printf("%v\n", err)
	}

	pingCache(page)

	if page.Body == nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	etag := fmt.Sprintf("%x", sha256.Sum256([]byte(page.Modtime.String())))

	w.Header().Set("ETag", "\""+etag+"\"")
	w.Header().Set("Content-Type", htmlutf8)
	w.Header().Set("Link", "</>; rel=\"contents\", </css>; rel=\"stylesheet\"")
	_, err = w.Write(page.Body)
	if err != nil {
		log500(w, r, err)
		return
	}
	log200(r)
}

// Handler for viewing the index page.
func indexHandler(w http.ResponseWriter, r *http.Request) {
	pingCache(indexCache)

	etag := fmt.Sprintf("%x", sha256.Sum256([]byte(indexCache.page.Modtime.String())))

	w.Header().Set("ETag", "\""+etag+"\"")
	w.Header().Set("Content-Type", htmlutf8)
	w.Header().Set("Link", "</>; rel=\"contents\", </css>; rel=\"stylesheet\"")
	_, err := w.Write(indexCache.page.Body)
	if err != nil {
		log500(w, r, err)
		return
	}
	log200(r)
}

// Serves the favicon as a URL.
// This is due to the default behavior of
// not serving naked paths but virtual ones.
func iconHandler(w http.ResponseWriter, r *http.Request) {
	confVars.mu.RLock()
	assetsDir := confVars.assetsDir
	iconPath := confVars.iconPath
	confVars.mu.RUnlock()

	longname := assetsDir + "/" + iconPath
	icon, err := ioutil.ReadFile(longname)
	if err != nil {
		if os.IsNotExist(err) {
			log.Printf("Favicon file specified in config does not exist: /icon request 404\n")
			error404(w, r)
			return
		}
		log500(w, r, err)
		return
	}

	stat, err := os.Stat(longname)
	if err != nil {
		log.Printf("Couldn't stat icon to send ETag header: %v\n", err.Error())
	}

	etag := fmt.Sprintf("%x", sha256.Sum256([]byte(stat.ModTime().String())))

	w.Header().Set("ETag", "\""+etag+"\"")
	w.Header().Set("Content-Type", http.DetectContentType(icon))
	_, err = w.Write(icon)
	if err != nil {
		log500(w, r, err)
		return
	}
	log200(r)
}

// Serves the local css file as a url.
// This is due to the default behavior of
// not serving naked paths but virtual ones.
func cssHandler(w http.ResponseWriter, r *http.Request) {
	confVars.mu.RLock()
	cssPath := confVars.cssPath
	confVars.mu.RUnlock()

	// check if using local or remote CSS.
	// if remote, don't bother doing anything
	// and redirect requests to /
	if !cssLocal([]byte(cssPath)) {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	css, err := ioutil.ReadFile(cssPath)
	if err != nil {
		if os.IsNotExist(err) {
			log.Printf("CSS file specified in config does not exist: /css request 404\n")
			error404(w, r)
			return
		}
		log500(w, r, err)
		return
	}

	stat, err := os.Stat(cssPath)
	if err != nil {
		log.Printf("Couldn't stat CSS file to send ETag header: %v\n", err.Error())
	}

	etag := fmt.Sprintf("%x", sha256.Sum256([]byte(stat.ModTime().String())))

	w.Header().Set("ETag", "\""+etag+"\"")
	w.Header().Set("Content-Type", cssutf8)
	_, err = w.Write(css)
	if err != nil {
		log500(w, r, err)
		return
	}
	log200(r)
}
