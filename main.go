package main // import "github.com/gbmor/tildewiki"

import (
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/spf13/viper"
)

// TildeWiki version
const twvers = "0.3"

// Page struct for caching
type Page struct {
	Longname  string
	Shortname string
	Title     string
	Desc      string
	Author    string
	Modtime   time.Time
	Body      []byte
	Raw       []byte
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
	setUpUsTheWiki()
	log.Println("Building initial cache ...")
	// fill the page cache
	genPageCache()

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/w/", validatePath(pageHandler))
	http.HandleFunc("/css", cssHandler)
	http.HandleFunc("/icon", iconHandler)

	port := ":" + viper.GetString("Port")
	log.Println("Binding to " + port)
	log.Fatal(http.ListenAndServe(port, nil))
}
