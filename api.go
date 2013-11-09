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
	api_key string = "2qnhTBYf.uKGAHHYy4PbaAQUMXs2Ux5c"
	uri     string = "https://app.asana.com/api/1.0/users/me"
)

func GetMe() (me *User, err error) {
	req, err := http.NewRequest("GET", uri, nil)
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
