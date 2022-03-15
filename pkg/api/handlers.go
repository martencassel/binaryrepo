package api

import (
	"fmt"
	"net/http"
)

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome!")
}

func CreateRepo(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Create Repo")
}

