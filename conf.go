package main

import (
	"html/template"
	"log"
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
