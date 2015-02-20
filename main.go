package main

import (
	"flag"
	"fmt"
	"github.com/gorilla/sessions"
	"log"
	"net/http"
)

var store = sessions.NewCookieStore([]byte("take-me-out-of-code"))

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

	session, _ := store.Get(req, "user")

	if user_id, exists := session["user_id"]; !exists {
		http.Redirect(w, req, "/login", 302)
		return
	}

	fmt.Fprintln(w, "This does nothing :(")
}
