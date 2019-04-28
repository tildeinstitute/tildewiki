package main

import (
	"bytes"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"mime"
	"os"

	"github.com/spf13/viper"
)

// displays on startup
func setUpUsTheWiki() {
	fmt.Printf(`
   __  _ __    __             _ __   _
  / /_(_) /___/ /__ _      __(_) /__(_)
 / __/ / / __  / _ \ | /| / / / //_/ /
/ /_/ / / /_/ /  __/ |/ |/ / / ,< / /
\__/_/_/\__,_/\___/|__/|__/_/_/|_/_/ 

         :: TildeWiki 0.4 ::
    (c)2019 Ben Morrison (gbmor)
               GPL v3
  https://github.com/gbmor/tildewiki
    All Contributions Appreciated!
		`)
	fmt.Printf("\n")
}

// open the icon file and process it
func iconType(filename string) string {
	file, err := os.Open(filename)
	if err != nil {
		log.Println("Couldn't open icon, sending mime type image/unknown")
		return "image/unknown"
	}

	// defer closing and checking of the error returned from (*os.File).Close()
	defer func() {
		err := file.Close()
		if err != nil {
			log.Printf("Deferred closing of favicon resulted in error: %v\n", err)
		}
	}()

	// pull the metadata from the image so we know
	// what mime type to send in the http header later
	_, format, err := image.DecodeConfig(file)
	if err != nil {
		log.Println("Can't decode icon, sending mime type image/unknown")
		return "image/unknown"
	}
	mime := mime.TypeByExtension("." + format)
	return mime
}

// determine if using local or remote css
// by checking if it's a URL or not
func cssLocal() bool {
	css := []byte(viper.GetString("CSS"))
	if bytes.HasPrefix(css, []byte("http://")) || bytes.HasPrefix(css, []byte("https://")) {
		return false
	}
	return true
}
