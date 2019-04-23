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

func initConfigParams() (*template.Template, *regexp.Regexp) {
	conf := viper.GetViper()
	conf.SetConfigType("yaml")
	conf.SetConfigName("tildewiki")
	conf.AddConfigPath("/usr/local/etc/")
	conf.AddConfigPath("/etc/")
	conf.AddConfigPath("$HOME/.tildewiki")
	conf.AddConfigPath(".")

	err := conf.ReadInConfig()
	if err != nil {
		log.Fatalln("Config file error: ", err)
	}

	conf.WatchConfig()
	conf.OnConfigChange(func(e fsnotify.Event) {
		log.Println("Config file change detected: ", e.Name)
	})

	var Templates = template.Must(template.ParseFiles(viper.GetString("TmplDir")+"/edit.html", viper.GetString("TmplDir")+"/view.html"))
	var ValidPath = regexp.MustCompile(viper.GetString("ValidPath"))

	return Templates, ValidPath
}

func error500(w http.ResponseWriter, r *http.Request) {
	e500 := viper.GetString("IndexDir") + "/500.md"
	file, err := ioutil.ReadFile(e500)
	if err != nil {
		http.NotFound(w, r)
	}
	w.Header().Set("Content-Type", "text/html")
	w.Write(render(file, viper.GetString("CSS"), "500 Error"))
}

func error404(w http.ResponseWriter, r *http.Request) {
	e404 := viper.GetString("IndexDir") + "/404.md"
	file, err := ioutil.ReadFile(e404)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "text/html")
	w.Write(render(file, viper.GetString("CSS"), "404 Error"))
}
