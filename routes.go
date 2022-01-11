package main

import (
	"log"

	"github.com/gorilla/mux"
)

func addRoutes() {
	r := mux.NewRouter()
	log.Print(r)
}
