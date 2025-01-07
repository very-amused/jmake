package main

import (
	"log"
	"os"
)

func main() {
	// Output to stderr
	log.SetOutput(os.Stderr)

	// Parse jmake.conf
	jmake, err := ParseConfig()
	if err != nil {
		log.Fatalln("Failed to parse jmake.toml:", err)
	}
	jmake.Keys, err = ParseConfigKeys()
	if err != nil {
		log.Fatalln("Failed to parse key order from jmake.toml:", err)
	}

	// Generate output files
	jmake.Generate()
}
