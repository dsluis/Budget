package db

import (
	couch "github.com/fjl/go-couchdb"
	"net/http"
	"encoding/json"
	"io/ioutil"
)

var client *couch.Client
var db *couch.DB
var user_db *couch.DB
var uuids []string = make([]string,0, 100)
var queue []string = uuids

func Connect() error {
	var err error
	client, err = couch.NewClient("http://localhost:5984",nil)
	if err != nil {
		return err
	}
	db, err = client.EnsureDB("budget")
	if( err != nil ) {
		return err
	}
	user_db, err = client.EnsureDB("user_db"
	if( err != nil ) {
		return err
	})
	return nil
}

func CreateUser() error {
	_, err := db.Put("bob", "bob", "")
	if err != nil {
		return err
	}
	return nil
}

func getUUID() error {
	res, err := http.Get("http://127.0.0.1:5984/_uuids?count=10")
	if err != nil {
		log.Fatal(err)
	}
	json.Unmarshal(res, &uuids)
	if err != nil {
		return err
	}
}
