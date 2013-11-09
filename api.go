package goasana

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

var client http.Client

type Workspace struct {
	Id   int
	Name string
}

type User struct {
	Email      string
	Name       string
	Photo      map[string]string
	Workspaces []Workspace
}

const (
	main_uri string = "https://app.asana.com/api/1.0"
	user_uri string = "/users"
	me_uri   string = "/me"
)

var api_key string

func setApiKey(key string) {
	api_key = key
}

func init() {
	key, err := ioutil.ReadFile("api.key")
	if err != nil {
		panic(err)
	}
	api_key = string(key)
}

func GetMe() (me *User, err error) {
	req, err := http.NewRequest("GET", main_uri+user_uri+me_uri, nil)
	if err != nil {
		return
	}
	req.SetBasicAuth(api_key, "")
	ris, err := client.Do(req)
	if err != nil {
		return
	}
	data, err := ioutil.ReadAll(ris.Body)
	if err != nil {
		return
	}
	var temp struct {
		Data User
	}
	json.Unmarshal(data, &temp)
	return &temp.Data, nil
}

func GetUsers() (users []User, err error) {
	req, err := http.NewRequest("GET", main_uri+user_uri, nil)
	if err != nil {
		return
	}
	req.SetBasicAuth(api_key, "")
	ris, err := client.Do(req)
	if err != nil {
		return
	}
	data, err := ioutil.ReadAll(ris.Body)
	if err != nil {
		return
	}
	var temp struct {
		Data []User
	}
	json.Unmarshal(data, &temp)
	return temp.Data, nil
}
