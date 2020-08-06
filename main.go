package main

import(
	"net/http"
	"fmt"
	//"log"
	"flag"
	"retargetly-exercise/handlers"
	"retargetly-exercise/dataprocessing"
)

func main(){
	port := flag.Int("port", 3000, "the port to start the CYOA web application on")

	http.HandleFunc("/login", handlers.LoginHandler)
	http.HandleFunc("/files/", handlers.FilesHandler)

	fmt.Printf("Starting server on port: %d\n", *port)
	//log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), nil))

	dataprocessing.StartProcess("./data/file1.tsv")
	// dataprocessing.ProcessFile("./data/test1.tsv")
}