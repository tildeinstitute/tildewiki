package main // import "github.com/gbmor/tildewiki"

import (
	"log"
	"net/http"
	"sync"
	"time"
)

// TildeWiki version
const twvers = "development version"

// Page struct for caching
type Page struct {
	Filename string
	Title    string
	Desc     string
	Modtime  time.Time
	Body     []byte
	Raw      []byte
}

// the in-memory page cache
var cachedPages = make(map[string]Page)

// initialize the basic configuration and
// assign the parsed templates and compiled regex
// to these two globals
//var templates, validPath = initConfigParams()
var validPath = initConfigParams()

// prevent concurrent writes to the cache
var mutex = &sync.RWMutex{}

func main() {
	// fill the page cache
	genPageCache()

	http.HandleFunc("/", welcomeHandler)
	http.HandleFunc("/w/", makeHandler(viewHandler))
	http.HandleFunc("/css", cssHandler)
	http.HandleFunc("/icon", iconHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
