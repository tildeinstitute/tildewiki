package main // import "github.com/gbmor/tildewiki"

import (
	"log"
	"net/http"
	"runtime"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	go func() {
		http.HandleFunc("/", welcomeHandler)
	}()
	go func() {
		http.HandleFunc("/w/", makeHandler(viewHandler))
	}()
	//http.HandleFunc("/edit/", makeHandler(editHandler))
	//http.HandleFunc("/save/", makeHandler(saveHandler))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
