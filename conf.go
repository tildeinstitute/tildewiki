package main

import (
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// Sets the basic parameters for the default viper (config library) instance
func initConfigParams() (*template.Template, *regexp.Regexp) {
	conf := viper.GetViper()

	// type of config file to look for
	conf.SetConfigType("yaml")
	// name of config file to look for
	conf.SetConfigName("tildewiki")
	// location of said config file
	conf.AddConfigPath("/usr/local/etc/")
	conf.AddConfigPath("/etc/")
	conf.AddConfigPath("$HOME/.tildewiki")
	conf.AddConfigPath(".")

	err := conf.ReadInConfig()
	if err != nil {
		log.Fatalln("Config file error: ", err)
	}

	// WatchConfig() is a function provided by blackfriday that watches the config
	// file for any changes and automatically reloads it if needed
	conf.WatchConfig()
	conf.OnConfigChange(func(e fsnotify.Event) {
		log.Println("Config file change detected: ", e.Name)
	})

	// Parse the HTML template file(s) and compile the regex path validation)
	var Templates = template.Must(template.ParseFiles(viper.GetString("TmplDir")+"/edit.html", viper.GetString("TmplDir")+"/view.html"))
	var ValidPath = regexp.MustCompile(viper.GetString("ValidPath"))

	return Templates, ValidPath
}

// this is just a custom 500 page using a markdown doc
// in the primary data/config directory.
// if the markdown doc can't be read, default to
// net/http's error handling
func error500(w http.ResponseWriter, r *http.Request) {
	e500 := viper.GetString("IndexDir") + "/500.md"
	file, err := ioutil.ReadFile(e500)
	if err != nil {
		http.NotFound(w, r)
	}
	w.Header().Set("Content-Type", "text/html")
	w.Write(render(file, viper.GetString("CSS"), "500 Error"))
}

// same as the 500 page
func error404(w http.ResponseWriter, r *http.Request) {
	e404 := viper.GetString("IndexDir") + "/404.md"
	file, err := ioutil.ReadFile(e404)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "text/html")
	w.Write(render(file, viper.GetString("CSS"), "404 Error"))
}
