package db

import (
	couch "github.com/fjl/go-couchdb"
)

var client *couch.Client
var db *couch.DB

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
	return nil
}

func CreateUser() error {
	_, err := db.Put("bob", "bob", "bob")
	if err != nil {
		return err
	}
	return nil
}