package goasana

import (
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

var client http.Client

var api_key string

func init() {
	key, err := ioutil.ReadFile("api.key")
	if err != nil {
		panic(err)
	}
	api_key = string(key)
}

func setApiKey(key string) {
	api_key = key
}

func SendRequestWithFilters(method string, url string, filters map[string]string) (data []byte, err error) {
	var body io.ReadCloser
	if len(filters) != 0 {
		url += "?"
		params := make([]string, len(filters))
		for key, val := range filters {
			params = append(params, key+"="+val)
		}
		url += strings.Join(params, "&")
	}
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return
	}
	req.SetBasicAuth(api_key, "")
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	data, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func SendRequest(method string, url string) (data []byte, err error) {
	return SendRequestWithFilters(method, url, nil)
}
