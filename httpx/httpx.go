package httpx

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

var (
	// 日志输出
	httpLogger io.Writer = os.Stdout
	// 默认超时时间
	httpTimeout = 10 * time.Second
)

func SetLogger(w io.Writer) {
	httpLogger = w
}

type httpLog struct {
	Time string `json:"time"`
	Cost int64  `json:"cost"`

	Request struct {
		Url    string      `json:"url"`
		Method string      `json:"method"`
		Data   interface{} `json:"data"`
	} `json:"request"`

	Response struct {
		Status int             `json:"status"`
		Data   json.RawMessage `json:"data,omitempty"`
	} `json:"response"`

	Error string `json:"ERROR,omitempty"`
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

func HTTPGet(url string, params map[string]string) *httpx {
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

func HTTPPost(url string, data interface{}) *httpx {
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

func HTTPPostForm(url string, data map[string]string) *httpx {
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
		defer x.printLog(time.Now())
	}

	x.send()

	return x.response, x.err
}

// 发送请求，解析返回的json
func (x *httpx) SendAndParse(to interface{}) error {
	if !x.closeLog {
		defer x.printLog(time.Now())
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

func (x *httpx) printLog(start time.Time) {
	var hlog httpLog

	hlog.Time = start.Format("2006-01-02 15:04:05")
	hlog.Cost = time.Since(start).Milliseconds()

	hlog.Request.Url = x.request.URL.String()
	hlog.Request.Method = x.request.Method
	hlog.Request.Data = x.requestData

	hlog.Response.Data = x.responseData
	if x.response != nil {
		hlog.Response.Status = x.response.StatusCode
	}

	if x.err != nil {
		hlog.Error = x.err.Error()
	}

	var buff bytes.Buffer
	enc := json.NewEncoder(&buff)
	enc.SetEscapeHTML(false)
	err := enc.Encode(&hlog)
	if err != nil {
		httpLogger.Write([]byte("ERROR: " + err.Error()))
	}
	httpLogger.Write(buff.Bytes())
	httpLogger.Write([]byte("\r\n"))
}
