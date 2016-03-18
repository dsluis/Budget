package db

import (
	couch "github.com/fjl/go-couchdb"
	"net/http"
	"encoding/json"
    "log"
)

var client *couch.Client
var db *couch.DB
var user_db *couch.DB
var uuids []string

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
	user_db, err = client.EnsureDB("user_db")
	if err != nil  {
		return err
	}
	return nil
}

func CreateUser() error {
    uuid, err := nextUUID()
	_, err = db.Put(uuid, "bob", uuid)
	if err != nil {
		return err
	}
	return nil
}

func nextUUID() (string,error) {
    if len( uuids ) == 0 {
        err := getUUID()
        if err != nil {
            return "", err
        }
    }
    uuid := uuids[0]
    uuids = uuids[1:]
    
    return uuid, nil
}
func getUUID() error {
	res, err := http.Get("http://127.0.0.1:5984/_uuids?count=100")
	if err != nil {
		log.Fatal(err)
        return err
	}
    var body []byte
    _,err = res.Body.Read(body)
    if err != nil {
        log.Fatal(err)
        return err
    }
    res.Body.Close()
    var result map[string][]string
	err = json.Unmarshal(body, &result)
    if err != nil {
        log.Fatal(err)
        return err
    }
    queue, _ = result["uuids"]
    return nil
}
