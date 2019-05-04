package main

import (
	"bytes"
	"fmt"
)

// displays on startup
func setUpUsTheWiki() {
	fmt.Printf(`
   __  _ __    __             _ __   _
  / /_(_) /___/ /__ _      __(_) /__(_)
 / __/ / / __  / _ \ | /| / / / //_/ /
/ /_/ / / /_/ /  __/ |/ |/ / / ,< / /
\__/_/_/\__,_/\___/|__/|__/_/_/|_/_/ 

        :: TildeWiki ` + twvers + ` ::
    (c)2019 Ben Morrison (gbmor)
               GPL v3
  https://github.com/gbmor/tildewiki
    All Contributions Appreciated!
		`)
	fmt.Printf("\n")
}

// determine if using local or remote css
// by checking if it's a URL or not
func cssLocal(css []byte) bool {
	if bytes.HasPrefix(css, []byte("http://")) || bytes.HasPrefix(css, []byte("https://")) {
		return false
	}
	return true
}
