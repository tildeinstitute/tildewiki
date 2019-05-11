package main // import "github.com/gbmor/tildewiki"

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/handlers"
)

// TildeWiki version
const twvers = "0.6.2"

// Makes the deferred close functions for the log file
// block until exit
var closelog = make(chan bool, 1)

func main() {

	// watch for SIGINT aka ^C
	// close the log file then exit
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	go func() {
		for sigint := range c {
			log.Printf("\n\nCaught %v. Cleaning up ...\n", sigint)

			if confVars.fileLogging {
				// signal to close the log file
				closelog <- true
				time.Sleep(50 * time.Millisecond)
			}

			close(closelog)
			os.Exit(0)
		}
	}()

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

	// signal to close the log file
	if confVars.fileLogging || confVars.quietLogging {
		closelog <- true
		close(closelog)
	}
}
