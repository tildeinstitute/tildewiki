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

	confVars.mu.RLock()
	filog := confVars.fileLogging
	qlog := confVars.quietLogging
	reversed := confVars.reverseTally
	portnum := confVars.port
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
		closelog <- true
		close(closelog)
	}
}
