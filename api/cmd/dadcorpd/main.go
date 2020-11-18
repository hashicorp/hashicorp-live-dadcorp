package main

import (
	"log"
	"net/http"
	"os"

	"dadcorp.dev/api"
)

func main() {
	storer, err := api.NewStorer()
	if err != nil {
		log.Println("Error setting up storer:", err.Error())
		os.Exit(1)
	}
	a := api.API{
		Storer: storer,
	}

	http.Handle("/", a.Server(""))
	err = http.ListenAndServe(":12345", nil)
	if err != nil {
		log.Println("Error listening and serving:", err.Error())
		os.Exit(1)
	}
}
