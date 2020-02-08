package main // import "github.com/gbmor/tildewiki"

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

// TildeWiki version
const twvers = "0.6.3"

// Makes the deferred close functions for the log file
// block until exit
var closelog = make(chan struct{}, 1)

func main() {
	confVars.mu.RLock()
	filog := confVars.fileLogging
	portnum := confVars.port
	qlog := confVars.quietLogging
	reversed := confVars.reverseTally
	viewPath := confVars.viewPath
	confVars.mu.RUnlock()

	// watch for SIGINT aka ^C
	// close the log file then exit
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	go func() {
		for sigint := range c {
			log.Printf("\n\nCaught %v. Cleaning up ...\n", sigint)

			if filog {
				// signal to close the log file
				closelog <- struct{}{}
				time.Sleep(50 * time.Millisecond)
			}

			close(closelog)
			os.Exit(0)
		}
	}()

	// fill the page cache
	log.Println("**NOTICE** Building initial cache ...")
	genPageCache()

	serv := mux.NewRouter().StrictSlash(true)

	serv.Path("/").HandlerFunc(indexHandler)
	serv.Path(viewPath + "{pageReq:[a-zA-Z0-9_-]+}").HandlerFunc(pageHandler)
	serv.Path("/css").HandlerFunc(cssHandler)
	serv.Path("/icon").HandlerFunc(iconHandler)
	serv.Path("/500").HandlerFunc(error500)
	serv.Path("/404").HandlerFunc(error404)

	// let the user know if using reversed page listings
	if reversed {
		log.Printf("**NOTICE** Using reversed page listings on index ... \n")
	}

	log.Println("**NOTICE** Binding to " + portnum)

	server := &http.Server{
		Handler:      handlers.CompressHandler(ipMiddleware(serv)),
		Addr:         portnum,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	err := server.ListenAndServe()
	if err != nil {
		log.Printf("%v\n", err.Error())
	}

	// signal to close the log file
	if filog || qlog {
		closelog <- struct{}{}
		close(closelog)
	}
}
