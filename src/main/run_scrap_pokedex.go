package main

import (
	"github.com/jotitan/pokedex/pokedex"
	"github.com/jotitan/pokedex/server"
	"log"
	"os"
	"strings"
)

func main() {
	if len(os.Args) < 3 {
		log.Fatal("Must specify action and at least the folder")
	}

	if strings.EqualFold("scrap", os.Args[1]) {
		pokedex.Scrap(os.Args[2])
	}
	if strings.EqualFold("scrapd", os.Args[1]) {
		pokedex.ScrapDresseur(os.Args[2])
	}
	if strings.EqualFold("scrape", os.Args[1]) {
		pokedex.ScrapEneryies(os.Args[2])
	}
	if strings.EqualFold("scrapi", os.Args[1]) {
		pokedex.ScrapImage(os.Args[2])
	}
	if strings.EqualFold("analyze", os.Args[1]) {
		pokedex.Analyze(os.Args[2])
	}
	if strings.EqualFold("serve", os.Args[1]) {
		server.RunServer(os.Args[2])
	}
}
