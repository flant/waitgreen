package main

import (
	"flag"
	"log"

	"waitgreen/modules/config"
	"waitgreen/modules/wg"
)

var (
	configfile string
	vBuild     string
	cnf        config.Config
)

func init() {
	flag.StringVar(&configfile, "config", "main.yml", "Read configuration from this file")
	flag.StringVar(&configfile, "f", "main.yml", "Read configuration from this file")
	flag.Parse()

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	log.Println("Bootstrap: build num.", vBuild)

	cnf = config.Parse(configfile)
	log.Println("Bootstrap: successful parsing config file.")

}

func main() {
	wg.Run(cnf)
}
