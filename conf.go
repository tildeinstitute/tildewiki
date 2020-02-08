package main

import (
	"log"
	"regexp"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// content-type constants
const htmlutf8 = "text/html; charset=utf-8"
const cssutf8 = "text/css; charset=utf-8"

// Config object initialization
var confVars = &confParams{}

// (Re-)Populates config object
func setConfVars() {
	confVars.port = ":" + viper.GetString("Port")
	confVars.pageDir = viper.GetString("PageDir")
	confVars.assetsDir = viper.GetString("AssetsDir")
	confVars.cssPath = viper.GetString("CSS")
	confVars.viewPath = "/" + viper.GetString("ViewPath") + "/"
	confVars.indexRefreshInterval = viper.GetString("IndexRefreshInterval")
	confVars.wikiName = viper.GetString("Name")
	confVars.wikiDesc = viper.GetString("ShortDesc")
	confVars.descSep = viper.GetString("DescSeparator")
	confVars.titleSep = viper.GetString("TitleSeparator")
	confVars.iconPath = viper.GetString("Icon")
	confVars.indexFile = viper.GetString("Index")
	confVars.reverseTally = viper.GetBool("ReverseTally")
	confVars.validPath = regexp.MustCompile(viper.GetString("ValidPath"))
	confVars.quietLogging = viper.GetBool("QuietLogging")
	confVars.fileLogging = viper.GetBool("FileLogging")
	confVars.logFile = viper.GetString("LogFile")
}

// Sets the basic parameters for the default viper (config library) instance
func initConfigParams() {
	conf := viper.GetViper()

	conf.SetConfigType("yaml")
	conf.SetConfigName("tildewiki")

	conf.AddConfigPath(".")
	conf.AddConfigPath("$HOME/.config/")
	conf.AddConfigPath("/usr/local/tildewiki/")
	conf.AddConfigPath("/etc/")
	conf.AddConfigPath("/usr/local/etc/")

	err := conf.ReadInConfig()
	if err != nil {
		log.Fatalf("Config file error: %s\n", err.Error())
	}

	setConfVars()

	conf.WatchConfig()
	conf.OnConfigChange(func(e fsnotify.Event) {
		log.Println("**NOTICE** Config file change detected: ", e.Name)
		setConfVars()
		triggerRecache()
	})

}
