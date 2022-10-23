package main

import (
	"fmt"
	"log"
	"net/http"
)

const port = "80"

type Config struct {}

func main() {
	app := Config{}

	log.Printf("Starting server on port %s\n", port)

	// define http server
	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: app.routes(),
	}

	// start the server
	err := server.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}

}
