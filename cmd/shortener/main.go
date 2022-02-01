package main

import (
	"log"

	"github.com/vstdy0/go-project/cmd/shortener/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
