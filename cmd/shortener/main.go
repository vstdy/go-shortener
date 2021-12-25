package main

import (
	"github.com/vstdy0/go-project/api"
	"log"
	"net/http"
)

func main() {
	r := api.Router()
	log.Fatal(http.ListenAndServe(":8080", r))
}
