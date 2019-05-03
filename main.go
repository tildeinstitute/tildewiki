package main // import "github.com/gbmor/tildewiki"

import (
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/spf13/viper"
)

// TildeWiki version
const twvers = "0.4"

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

// content-type constants
const htmlutf8 = "text/html; charset=utf-8"
const cssutf8 = "text/css; charset=utf-8"

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
	// show the logo, repo link, etc
	setUpUsTheWiki()

	// set up logging if the config file params
	// are set
	if err := logSetup(); err != nil {
		log.Printf("Couldn't set up non-default logging: %v\n", err)
	}

	// fill the page cache
	log.Println("**NOTICE** Building initial cache ...")
	genPageCache()

	viewpath := "/" + viper.GetString("ViewPath") + "/"

	http.HandleFunc("/", indexHandler)
	http.HandleFunc(viewpath, validatePath(pageHandler))
	http.HandleFunc("/css", cssHandler)
	http.HandleFunc("/icon", iconHandler)

	port := ":" + viper.GetString("Port")
	log.Println("**NOTICE** Binding to " + port)

	// let the user know if using reversed page listings
	if viper.GetBool("ReverseTally") {
		log.Printf("**NOTICE** Using reversed page listings on index ... \n")
	}

	log.Fatal(http.ListenAndServe(port, nil))
}
