package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	h "retargetly-exercise/handlers"
	"time"
)

func main() {
	port := flag.Int("port", 3000, "the port to start the web application on")

	http.HandleFunc("/login", h.LoginHandler)
	http.HandleFunc("/files/", h.FilesHandler)

	fmt.Printf("Starting server on port: %d\n", *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), nil))
}

type Message struct{
	Msg string `json:"message"`
}

func dashhandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Connection", "keep-alive")
	msg := Message{
		Msg: "First message",
	}
	msg2  := Message{
		Msg: "sending second line of data",
	}
	// fmt.Fprintf(w, msg)
	json.NewEncoder(w).Encode(msg)
  if f, ok := w.(http.Flusher); ok {
     f.Flush()
  } else {
     log.Println("Damn, no flush");
  }
  m, _ := time.ParseDuration("5s")
	time.Sleep(m)
	json.NewEncoder(w).Encode(msg2)

  //fmt.Fprintf(w, "sending second line of data")
}


