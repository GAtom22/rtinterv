package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	h "retargetly-exercise/handlers"
)

func main() {
	port := flag.Int("port", 3000, "the port to start the web application on")

	http.HandleFunc("/login", h.LoginHandler)
	http.HandleFunc("/files/", h.FilesHandler)

	fmt.Printf("Starting server on port: %d\n", *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), nil))
}
