package main

import (
	"log"
	"os"
)

func init() {

	// show the logo, repo link, etc
	setUpUsTheWiki()

	// initialize the configuration
	initConfigParams()

	// set up logging if the config file params
	// are set
	if confVars.fileLogging && !confVars.quietLogging {
		if llogfile, err := os.OpenFile(confVars.logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600); err == nil {
			log.SetOutput(llogfile)

			go func() {
				<-closelog
				log.Printf("Closing log file ...\n")
				err := llogfile.Close()
				if err != nil {
					log.Printf("Couldn't close log file: %v\n", err)
				}
			}()

		} else {
			log.Printf("Couldn't log to file: %v\n", err)
		}
	}

	// Tell TildeWiki to be quiet,
	if confVars.quietLogging {
		if llogfile, err := os.Open("/dev/null"); err == nil {
			log.SetOutput(llogfile)

			go func() {
				// I don't know why I'm bothering to do this for /dev/null
				// ...
				// whatever
				<-closelog
				log.Printf("Closing log file ...\n")
				err := llogfile.Close()
				if err != nil {
					log.Printf("Couldn't close log file: %v\n", err)
				}
			}()

		} else {
			log.Printf("Couldn't quiet logging: %v\n", err)
		}
	}
}
