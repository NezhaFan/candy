package candy

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/valyala/bytebufferpool"
)

type Request struct {
	timeout uint8             // 请求超时（秒）
	Url     string            `json:"url"`
	Method  string            `json:"method"`
	Body    json.RawMessage   `json:"body,omitempty"`
	Header  map[string]string `json:"header,omitempty"`
}

func HttpGet(url string, params map[string]string) *Request {
	r := Request{
		Method: http.MethodGet,
		Url:    url,
	}

	if len(params) > 0 {
		buf := bytebufferpool.Get()
		defer bytebufferpool.Put(buf)
		for k, v := range params {
			buf.WriteString(k + "=" + v + "&")
		}
		s := buf.String()[:buf.Len()-1]

		r.Url += "?" + s
	}

	return &r
}

func HttpPostJson(url string, data any) *Request {
	r := Request{
		Method: http.MethodPost,
		Url:    url,
		Header: map[string]string{"Content-Type": "application/json"},
	}
	if data != nil {
		bs, _ := json.Marshal(data)
		r.Body = json.RawMessage(bs)
	}

	return &r
}

func HttpPostForm(url string, data map[string]any) *Request {
	r := Request{
		Method: http.MethodPost,
		Url:    url,
		Header: map[string]string{"Content-Type": "application/x-www-form-urlencoded"},
		Body:   buildForm(data),
	}
	return &r
}

func HttpPutJson(url string, data any) *Request {
	r := Request{
		Method: http.MethodPut,
		Url:    url,
		Header: map[string]string{"Content-Type": "application/json"},
	}
	if data != nil {
		bs, _ := json.Marshal(data)
		r.Body = json.RawMessage(bs)
	}

	return &r
}

func HttpPutForm(url string, data map[string]any) *Request {
	r := Request{
		Method: http.MethodPut,
		Url:    url,
		Header: map[string]string{"Content-Type": "application/x-www-form-urlencoded"},
		Body:   buildForm(data),
	}
	return &r
}

func HttpDelete(url string) *Request {
	r := Request{
		Method: http.MethodDelete,
		Url:    url,
	}
	return &r
}

func HttpHead(url string) *Request {
	r := Request{
		Method: http.MethodHead,
		Url:    url,
	}
	return &r
}

func buildForm(data map[string]any) []byte {
	if len(data) == 0 {
		return nil
	}
	buf := bytebufferpool.Get()
	defer bytebufferpool.Put(buf)
	for k, v := range data {
		buf.WriteString(k + "=" + fmt.Sprintf("%v", v) + "&")
	}
	return buf.Bytes()[:buf.Len()-1]
}

func (r *Request) SetHeader(key, value string) *Request {
	if r.Header == nil {
		r.Header = make(map[string]string)
	}
	r.Header[key] = value
	return r
}

func (r *Request) SetTimeout(second uint8) *Request {
	r.timeout = second
	return r
}

func (r *Request) Do() (*http.Response, error) {
	req, err := http.NewRequest(r.Method, r.Url, bytes.NewReader(r.Body))
	if err != nil {
		return nil, err
	}

	for k, v := range r.Header {
		req.Header.Set(k, v)
	}

	client := &http.Client{Timeout: time.Duration(r.timeout) * time.Second}
	if r.timeout == 0 {
		client.Timeout = time.Second * 10
	} else {
		client.Timeout = time.Duration(r.timeout) * time.Second
	}

	resp, err := client.Do(req)

	// 判断状态码
	if resp.StatusCode > http.StatusCreated && err == nil {
		err = errors.New("HTTP返回状态码异常: " + resp.Status)
	}

	return resp, err
}

func (r *Request) DoAndParse(data any) error {
	resp, err := r.Do()
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, data)
}
