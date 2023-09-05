package httpx

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

var (
	// 默认超时时间
	httpTimeout = 10 * time.Second
)

func SetLogger(w io.Writer) {
	httpLogger = w
}

type httpx struct {
	request  *http.Request
	response *http.Response

	requestData  interface{}
	responseData json.RawMessage

	closeLog bool
	timeout  time.Duration
	err      error
}

func Get(url string, params map[string]string) *httpx {
	var x httpx
	if len(params) > 0 {
		var buff bytes.Buffer
		buff.Grow(512)
		for k, v := range params {
			buff.WriteString(k + "=" + v + "&")
		}

		url += "?" + buff.String()[:buff.Len()-1]
	}

	x.request, x.err = http.NewRequest(http.MethodGet, url, nil)
	if x.err == nil {
		x.timeout = httpTimeout
	}

	return &x
}

func Post(url string, data interface{}) *httpx {
	var x httpx
	var body io.Reader
	if data != nil {
		bs, _ := json.Marshal(data)
		x.requestData = json.RawMessage(bs)
		body = bytes.NewReader(bs)
	}

	x.request, x.err = http.NewRequest(http.MethodPost, url, body)
	if x.err == nil {
		x.request.Header.Set("Content-Type", "application/json")
		x.timeout = httpTimeout
	}

	return &x
}

func PostForm(url string, data map[string]string) *httpx {
	var x httpx
	var body io.Reader
	if data != nil {
		var buff bytes.Buffer
		buff.Grow(512)
		for k, v := range data {
			buff.WriteString(k + "=" + fmt.Sprintf("%v", v) + "&")
		}
		bs := buff.Bytes()[:buff.Len()-1]
		x.requestData = string(bs)
		body = bytes.NewReader(bs)
	}

	x.request, x.err = http.NewRequest(http.MethodPost, url, body)
	if x.err == nil {
		x.request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		x.timeout = httpTimeout
	}

	return &x
}

func (x *httpx) SetHeader(key, val string) *httpx {
	x.request.Header.Set(key, val)
	return x
}

func (x *httpx) SetTimeout(timeout time.Duration) *httpx {
	x.timeout = timeout
	return x
}

func (x *httpx) CloseLog() *httpx {
	x.closeLog = true
	return x
}

// 发送请求，无法记录返回的数据，须主动关闭Body
func (x *httpx) Send() (*http.Response, error) {
	if !x.closeLog {
		defer printLog(time.Now(), x)
	}

	x.send()

	return x.response, x.err
}

// 发送请求，并且通过回调函数处理，内部返回数据，用来日志记录
func (x *httpx) SendAndDeal(callback func(*http.Response, error) []byte) {
	if !x.closeLog {
		defer printLog(time.Now(), x)
	}

	x.send()
	x.responseData = callback(x.response, x.err)
	x.request.Body.Close()
}

// 发送请求，解析返回的json
func (x *httpx) SendAndParse(to interface{}) error {
	if !x.closeLog {
		defer printLog(time.Now(), x)
	}

	x.send()

	if x.err != nil {
		return x.err
	}

	var bs []byte
	bs, x.err = io.ReadAll(x.response.Body)
	if x.err != nil {
		return x.err
	}
	x.response.Body.Close()

	x.responseData = bs
	x.err = json.Unmarshal(bs, to)

	return x.err
}

func (x *httpx) send() {
	if x.err != nil || x.response != nil {
		return
	}

	// 发送请求
	client := &http.Client{Timeout: x.timeout}
	x.response, x.err = client.Do(x.request)

	if x.err == nil && x.response.StatusCode != http.StatusOK {
		x.err = fmt.Errorf("%d %s", x.response.StatusCode, http.StatusText(x.response.StatusCode))
	}
}
