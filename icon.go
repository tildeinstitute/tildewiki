package main

import (
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"log"
	"mime"
	"os"
)

// Guess image format from gif/jpeg/png/webp
func guessImageFormat(r io.Reader) (format string, err error) {
	_, format, err = image.DecodeConfig(r)
	return
}

// Guess image mime types from gif/jpeg/png/webp
func guessImageMimeTypes(r io.Reader) string {
	format, _ := guessImageFormat(r)
	if format == "" {
		return ""
	}
	return mime.TypeByExtension("." + format)
}

// open the icon file and process it
func iconType(filename string) string {
	r, err := os.Open(filename)
	if err != nil {
		log.Println("iconType() :: Couldn't open icon, sending mime type image/unknown")
		return "image/unknown"
	}
	defer r.Close()
	mime := guessImageMimeTypes(r)
	log.Println("iconType() :: " + mime)
	return mime
}
