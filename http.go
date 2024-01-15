package candy

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/valyala/bytebufferpool"
)

var (
	httpDefaultLogger  *Logger = nil
	httpDefaultTimeout uint8   = 5 // 秒
)

type Request struct {
	Time    string `json:"time"`
	Cost    uint16 `json:"cost"`
	timeout uint8
	needLog bool
	Url     string          `json:"url"`
	Method  string          `json:"method"`
	Body    json.RawMessage `json:"body,omitempty"`
	header  map[string]string
	Err     string `json:"error,omitempty"`
	buf     []byte
}

func Get(url string, params map[string]string) *Request {
	r := Request{
		Method:  http.MethodGet,
		Url:     url,
		header:  map[string]string{},
		timeout: httpDefaultTimeout,
	}
	if len(params) > 0 {
		var buff bytes.Buffer
		buff.Grow(512)
		for k, v := range params {
			buff.WriteString(k + "=" + v + "&")
		}
		r.Url += "?" + buff.String()[:buff.Len()-1]
	}

	return &r
}

func PostJson(url string, data interface{}) *Request {
	r := Request{
		Method:  http.MethodPost,
		Url:     url,
		header:  map[string]string{"Content-Type": "application/json"},
		timeout: httpDefaultTimeout,
	}
	if data != nil {
		bs, _ := json.Marshal(data)
		r.Body = json.RawMessage(bs)
	}

	return &r
}

func PostForm(url string, data map[string]any) *Request {
	r := Request{
		Method:  http.MethodPost,
		Url:     url,
		header:  map[string]string{"Content-Type": "application/x-www-form-urlencoded"},
		timeout: httpDefaultTimeout,
	}
	if data != nil {
		var buff bytes.Buffer
		buff.Grow(512)
		for k, v := range data {
			buff.WriteString(k + "=" + fmt.Sprintf("%v", v) + "&")
		}
		r.Body = buff.Bytes()[:buff.Len()-1]
	}

	return &r
}

func (r *Request) OpenLog() *Request {
	r.needLog = true
	return r
}

func (r *Request) SetHeader(key, value string) *Request {
	r.header[key] = value
	return r
}

func (r *Request) SetTimeout(second uint8) *Request {
	httpDefaultTimeout = second
	return r
}

func (r *Request) Do() (*http.Response, error) {
	req, err := http.NewRequest(r.Method, r.Url, bytes.NewReader(r.Body))
	if err != nil {
		return nil, err
	}

	for k, v := range r.header {
		req.Header.Set(k, v)
	}

	client := &http.Client{Timeout: time.Duration(r.timeout) * time.Second}
	resp, err := client.Do(req)

	if resp.StatusCode > http.StatusCreated {
		err = errors.New("http status code: " + resp.Status)
	}

	if err != nil {
		r.Err = err.Error()
	}

	if r.needLog || err != nil {
		r.Time = time.Now().Format("2006-01-02 15:04:05")
		// 非json则套一下变成json格式
		if len(r.Body) > 0 && r.Body[0] != '{' {
			buf := bytebufferpool.Get()
			buf.Write([]byte(`{"":"`))
			buf.Write(r.Body)
			buf.Write([]byte(`"}`))
			r.Body = json.RawMessage(buf.Bytes())
			bytebufferpool.Put(buf)
		}
		b, _ := json.Marshal(r)
		httpDefaultLogger.Write(b)
		httpDefaultLogger.Write([]byte("\r\n"))
	}

	return resp, err
}

func BuildQuery(data map[string]string) []byte {
	buf := bytebufferpool.Get()
	defer bytebufferpool.Put(buf)
	for k, v := range data {
		buf.WriteString(k + "=" + v + "&")
	}
	return buf.Bytes()[:buf.Len()-1]
}
