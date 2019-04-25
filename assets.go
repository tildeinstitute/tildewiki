package main

import (
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"mime"
	"os"

	"github.com/spf13/viper"
)

// open the icon file and process it
func iconType(filename string) string {
	file, err := os.Open(filename)
	if err != nil {
		log.Println("iconType() :: Couldn't open icon, sending mime type image/unknown")
		return "image/unknown"
	}
	defer file.Close()
	_, format, err := image.DecodeConfig(file)
	if err != nil {
		log.Println("iconType() :: Can't decode icon, sending mime type image/unknown")
		return "image/unknown"
	}
	mime := mime.TypeByExtension("." + format)
	log.Println("iconType() :: " + mime)
	return mime
}

// determine if using local or remote css
func cssLocal() bool {
	css := viper.GetString("CSS")
	cssbyte := []byte(css)
	if string(cssbyte[:7]) == "http://" || string(cssbyte[:8]) == "https://" {
		return false
	}
	return true
}
