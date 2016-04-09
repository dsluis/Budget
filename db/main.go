package db

import (
	"encoding/json"
	couch "github.com/fjl/go-couchdb"
	"log"
	"net/http"
    "io/ioutil"
    "errors"
)

var client *couch.Client
var db *couch.DB
var userDb *couch.DB
var uuids []string

type ViewResponse struct {
    TotalRows int `json:"total_rows"`
    Offset int `json:"offset"`
    Rows []interface{} `json:"rows"`
}
//todo: store password hash instead
type User struct {
	Username string
	Password string
}

func Connect() error {
	var err error
	client, err = couch.NewClient("http://localhost:5984", nil)
	if err != nil {
		return err
	}
	db, err = client.EnsureDB("budget")
	if err != nil {
		return err
	}
	userDb, err = client.EnsureDB("budget_users")
	if err != nil {
		return err
	}
	return nil
}

func LoginUser(user string, pass string) (string,error) {
    options := make( couch.Options )
    options["key"] = user
    resp := ViewResponse{}
    if err := userDb.View("_design/user","user",&resp, options); err != nil {
        return "", err
    }
    if len(resp.Rows) == 0 {
        return "",errors.New("User does not exist")
    }
    row := resp.Rows[0].(map[string]interface{})
    id, _ := row["id"].(string)
    return id,nil
}

func CreateUser(u User) error {
	uuid, err := nextUUID()
	if err != nil {
		return err
	}
	_, err = userDb.Put(uuid, u, "")
	if err != nil {
		return err
	}
	return nil
}

func nextUUID() (string, error) {
	if len(uuids) == 0 {
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
    defer res.Body.Close()
    body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
		return err
	}
	var result map[string][]string
	err = json.Unmarshal(body, &result)
	if err != nil {
		log.Fatal(err)
		return err
	}
	uuids = result["uuids"]
	return nil
}
