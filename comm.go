package goasana

import (
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

type connInfo struct {
	simpleConn bool
	apiKey     string
	token      string
}

var client http.Client

func (conn connInfo) SendRequestWithFilters(method string, url string, filters map[string]string) (data []byte, err error) {
	var body io.ReadCloser
	if len(filters) != 0 {
		url += "?"
		params := make([]string, 0)
		for key, val := range filters {
			params = append(params, key+"="+val)
		}
		url += strings.Join(params, "&")
	}
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return
	}
	if conn.simpleConn {
		req.SetBasicAuth(conn.apiKey, "")
	} else {
		return nil, errors.New("Oatuh not implemented")
	}
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

func (conn connInfo) SendRequest(method string, url string) (data []byte, err error) {
	return conn.SendRequestWithFilters(method, url, nil)
}
