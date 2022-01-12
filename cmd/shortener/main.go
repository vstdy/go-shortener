package main

import (
	"github.com/vstdy0/go-project/cmd/root"
	"log"
)

func main() {
	cmd, err := root.NewRootCmd()
	if err != nil {
		log.Fatal(err)
	}
	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
