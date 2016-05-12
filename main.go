package main

import (
	"flag"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"html/template"
	"log"
	"net/http"
	"github.com/dsluis/budget/db"
    "fmt"
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	Salt string
	CookieSecret string
}

var store = sessions.NewCookieStore([]byte("take-me-out-of-code"))
var templates = template.Must(template.ParseFiles("views/user/login.html", "views/user/create.html", "views/home/index.html"))

func main() {

	var port = flag.String("port", ":8080", "network port to receive http requests over")
	var configPath = flag.String("config","config.json","path to json config file")

	flag.Parse()
	
	file, err := ioutil.ReadFile(*configPath)
	if err != nil {
		log.Fatal("Missing Config File")
	}
	var config Config
	if err := json.Unmarshal(file,&config); err != nil {
		log.Fatal("Invalid Config File: " + err.Error() )
	}
	
	db.Config.Salt = config.Salt
	
	if err := db.Connect(); err != nil {
		log.Fatal("Connect: ", err )
	}

	router := mux.NewRouter()

	router.HandleFunc("/", index)
	router.HandleFunc("/user/login", loginView).Methods("GET")
	router.HandleFunc("/user/login", loginAction).Methods("POST")
	router.HandleFunc("/user/create", createView).Methods("GET")
	router.HandleFunc("/user/create", createAction).Methods("POST")

	http.Handle("/", router)
	err = http.ListenAndServe(*port, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func index(w http.ResponseWriter, req *http.Request) {

	if ! authorize(w,req) {
		return
	}

	err := templates.ExecuteTemplate(w, "index.html", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func loginView(w http.ResponseWriter, req *http.Request) {
	session, _ := store.Get(req, "user")

	data := ViewData{"", nil}
	if flashes := session.Flashes("feedback"); len(flashes) > 0 {
		data.Feedback = flashes[0].(string)
		session.Save(req, w)
	}

	err := templates.ExecuteTemplate(w, "login.html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func loginAction(w http.ResponseWriter, req * http.Request) {
    session, _ := store.Get(req, "user")
    
    if err := req.ParseForm(); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    name := req.PostForm.Get( "user" )
    if name == "" {
        return
    }
    pw := req.PostForm.Get( "pass" )
    if pw == "" {
        return
    }
	
    userID, err := db.LoginUser(name,pw); 
    
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    session.Values["user_id"] = userID
    session.Save(req,w)
    
    fmt.Fprintln(w, userID)
    return
}
func createView(w http.ResponseWriter, req *http.Request) {

	err := templates.ExecuteTemplate(w, "create.html", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func createAction(w http.ResponseWriter, req *http.Request) {
	session, _ := store.Get(req, "user")
    
    if err := req.ParseForm(); err != nil {
       http.Error(w, err.Error(), http.StatusBadRequest)
       return
    }
    name := req.PostForm.Get("user")
    if name == "" {
        http.Error(w, "Please pass a username", http.StatusBadRequest)
        return
    }
    password := req.PostForm.Get("pass")
    if password == "" {
        http.Error(w, "Please enter a password", http.StatusBadRequest)
        return
    }
    
	//todo: validate
	//todo: create account
    user := db.User{ name, password }
    
    if err := db.CreateUser(user); err != nil {
        http.Error(w,err.Error(),http.StatusInternalServerError)
        return
    }

	session.AddFlash("Successfully Created Account", "feedback")
	session.Save(req, w)
	http.Redirect(w, req, "/user/login", 302)
}

func authorize(w http.ResponseWriter, req *http.Request) bool {
	session, _ := store.Get(req, "user")

	if _, exists := session.Values["user_id"]; !exists {
		http.Redirect(w, req, "/user/login",302)
		return false
	}
	return true
}
