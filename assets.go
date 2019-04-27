package main

import (
	"bytes"
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
	defer func() {
		err := file.Close()
		if err != nil {
			log.Printf("iconType() :: Deferred closing of image resulted in error: %v\n", err)
		}
	}()
	_, format, err := image.DecodeConfig(file)
	if err != nil {
		log.Println("iconType() :: Can't decode icon, sending mime type image/unknown")
		return "image/unknown"
	}
	mime := mime.TypeByExtension("." + format)
	return mime
}

// determine if using local or remote css
func cssLocal() bool {
	css := []byte(viper.GetString("CSS"))
	if bytes.HasPrefix(css, []byte("http://")) || bytes.HasPrefix(css, []byte("https://")) {
		return false
	}
	return true
}
