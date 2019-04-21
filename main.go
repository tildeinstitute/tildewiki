package main // import "github.com/gbmor/tildewiki"

import (
	"log"
	"net/http"
)

func main() {
	//mux := http.NewServeMux()
	http.HandleFunc("/", welcomeHandler)
	http.HandleFunc("/w/", makeHandler(viewHandler))
	http.HandleFunc("/edit/", makeHandler(editHandler))
	http.HandleFunc("/save/", makeHandler(saveHandler))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
