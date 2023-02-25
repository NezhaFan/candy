package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type client struct {
	request
	response
}

type request struct {
	method  string
	url     string
	timeout time.Duration
	headers map[string]string
	reqBody []byte
}

type response struct {
	Status   int
	RespBody []byte
}

func NewClient() *client {
	hc := &client{}
	hc.headers = make(map[string]string)
	hc.timeout = time.Second * 5
	return hc
}

func (hc *client) Get(api string, params map[string]string) *client {
	hc.method = http.MethodGet
	if len(params) > 0 {
		hc.url = api + "?" + string(hc.BuildQuery(params))
	} else {
		hc.url = api
	}

	return hc
}

func (hc *client) PostJson(api string, params any) *client {
	hc.method = http.MethodPost
	hc.url = api
	hc.AddHeader("Content-Type", "application/json")

	data, err := json.Marshal(params)
	if err != nil {
		fmt.Println("PostJson: " + err.Error())
	} else {
		hc.reqBody = data
	}
	return hc
}

func (hc *client) PostForm(api string, params map[string]string) *client {
	hc.method = http.MethodPost
	hc.url = api
	hc.AddHeader("Content-Type", "application/x-www-form-urlencoded")

	if len(params) > 0 {
		hc.reqBody = hc.BuildQuery(params)
	}

	return hc
}

func (hc *client) AddHeader(key, value string) *client {
	hc.headers[key] = value
	return hc
}

func (hc *client) SetTimeout(d time.Duration) *client {
	hc.timeout = d
	return hc
}

func (hc *client) BuildQuery(params map[string]string) []byte {
	if len(params) == 0 {
		return nil
	}
	var b bytes.Buffer
	for k, v := range params {
		b.Write([]byte(k + "=" + v + "&"))
	}

	return b.Bytes()[0 : b.Len()-1]
}

func (hc *client) Do() ([]byte, error) {
	req, err := http.NewRequest(hc.method, hc.url, bytes.NewBuffer(hc.reqBody))
	if err != nil {
		return nil, err
	}

	for k, v := range hc.headers {
		req.Header.Set(k, v)
	}

	client := &http.Client{Timeout: hc.timeout}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	hc.Status = resp.StatusCode
	hc.RespBody, err = io.ReadAll(resp.Body)
	return hc.RespBody, err
}
