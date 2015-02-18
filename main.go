package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

func main() {

	var port = flag.String("port", ":8080", "network port to receive http requests over")

	flag.Parse()

	http.HandleFunc("/", index)
	err := http.ListenAndServe(*port, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func index(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintln(w, "This does nothing :(")
}
