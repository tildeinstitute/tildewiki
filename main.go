package main // import "github.com/gbmor/tildewiki"

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
)

// TildeWiki version
const twvers = "0.6.1"

func main() {

	// show the logo, repo link, etc
	setUpUsTheWiki()

	// initialize the configuration
	initConfigParams()

	// set up logging if the config file params
	// are set
	if confVars.fileLogging {
		if llogfile, err := os.OpenFile(confVars.logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600); err == nil {
			log.SetOutput(llogfile)
			defer func() {
				err := llogfile.Close()
				if err != nil {
					log.Printf("Couldn't close log file: %v\n", err)
				}
			}()
		} else {
			log.Printf("Couldn't log to file: %v\n", err)
		}
	}
	// Tell Tildewiki to be quiet,
	// Supersedes file logging
	if confVars.quietLogging {
		if llogfile, err := os.Open("/dev/null"); err == nil {
			log.SetOutput(llogfile)
			defer func() {
				err := llogfile.Close()
				if err != nil {
					log.Printf("Couldn't close log file: %v\n", err)
				}
			}()
		} else {
			log.Printf("Couldn't quiet logging: %v\n", err)
		}
	}

	// fill the page cache
	log.Println("**NOTICE** Building initial cache ...")
	genPageCache()

	serv := http.NewServeMux()

	serv.HandleFunc("/", indexHandler)
	serv.HandleFunc(confVars.viewPath, validatePath(pageHandler))
	serv.HandleFunc("/css", cssHandler)
	serv.HandleFunc("/icon", iconHandler)
	serv.HandleFunc("/500", error500)
	serv.HandleFunc("/404", error404)

	log.Println("**NOTICE** Binding to " + confVars.port)

	// let the user know if using reversed page listings
	if confVars.reverseTally {
		log.Printf("**NOTICE** Using reversed page listings on index ... \n")
	}

	log.Fatal(http.ListenAndServe(confVars.port, handlers.CompressHandler(serv)))
}
