package candy

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
)

var (
	httpLogger io.Writer = os.Stdout
)

func SetHTTPLogger(logger io.Writer) {
	httpLogger = logger
}

type client struct {
	// 请求相关
	Url       string            `json:"url"`
	Method    string            `json:"method"`
	Timeout   time.Duration     `json:"-"`
	Headers   map[string]string `json:"-"`
	ReqBody   []byte            `json:"-"`
	ReqJson   json.RawMessage   `json:"req,omitempty"`
	ReqString string            `json:"req_string,omitempty"`

	// 返回相关
	Cost       string          `json:"cost,omitempty"`
	RespBody   json.RawMessage `json:"resp,omitempty"`
	RespString string          `json:"resp_string,omitempty"`

	// 其他
	OpenLog bool   `json:"-"`
	ErrMsg  string `json:"ERROR,omitempty"`
}

func newClient() *client {
	c := &client{}
	c.Timeout = time.Second * 10
	c.Headers = make(map[string]string)
	c.OpenLog = true
	return c
}

func HTTPGet(api string, params any) *client {
	c := newClient()
	c.Url = api
	c.Method = http.MethodGet
	if params != nil {
		// 防止一些自定义的参数已经拼接
		if !strings.Contains(c.Url, "?") {
			c.Url += "?"
		}
		c.Url += string(buildQuery(params))
	}

	return c
}

func HTTPPostJSON(api string, params any) *client {
	c := newClient()
	c.Method = http.MethodPost
	c.Url = api
	c.Headers["Content-Type"] = "application/json"

	if params != nil {
		// 参数限定为结构体或者map等可Json化的
		data, err := json.Marshal(params)
		if err != nil {
			panic("Failed to marshal JSON response body:" + err.Error())
		}

		c.ReqBody = data
	}

	return c
}

func HTTPPostForm(api string, params any) *client {
	c := newClient()
	c.Method = http.MethodPost
	c.Url = api
	c.Headers["Content-Type"] = "application/x-www-form-urlencoded"

	if params != nil {
		c.ReqBody = buildQuery(params)
	}

	return c
}

func buildQuery(params interface{}) []byte {
	if params == nil {
		return nil
	}

	tp := reflect.TypeOf(params)
	val := reflect.ValueOf(params)

	var buff bytes.Buffer
	buff.Grow(512)

	switch tp.Kind() {
	case reflect.Map:
		switch data := params.(type) {
		case map[string]string:
			for k, v := range data {
				buff.WriteString(k + "=" + v + "&")
			}
		case map[string]any:
			for k, v := range data {
				buff.WriteString(k + "=" + fmt.Sprintf("%v", v) + "&")
			}
		}
	case reflect.Pointer:
		tp = tp.Elem()
		val = val.Elem()
		fallthrough
	case reflect.Struct:
		for i := 0; i < tp.NumField(); i++ {
			k := tp.Field(i).Tag.Get("json")
			if k != "" {
				buff.WriteString(k + "=" + fmt.Sprintf("%v", val.Field(i).Interface()) + "&")
			}
		}
	}

	if buff.Len() == 0 {
		panic("参数类型不支持：" + tp.Kind().String())
	}

	return buff.Bytes()[:buff.Len()-1]
}

func (c *client) AddHeader(key, value string) *client {
	c.Headers[key] = value
	return c
}

func (c *client) SetTimeout(d time.Duration) *client {
	c.Timeout = d
	return c
}

func (c *client) CloseLog() *client {
	c.OpenLog = false
	return c
}

func (c *client) log(start time.Time, errPtr *error) {
	// 记录错误
	if *errPtr != nil {
		c.ErrMsg = (*errPtr).Error()
	}

	// JSON和非JSON格式处理
	if json.Valid(c.ReqBody) {
		c.ReqJson = c.ReqBody
	} else {
		c.ReqString = string(c.ReqBody)
	}

	if !json.Valid(c.RespBody) {
		c.RespString = string(c.RespBody)
		c.RespBody = nil
	}

	// 记录耗时
	c.Cost = strconv.Itoa(int(time.Since(start).Milliseconds())) + "ms"

	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	if e := encoder.Encode(c); e != nil {
		httpLogger.Write([]byte("ERROR: " + e.Error() + "\n"))
	}

	httpLogger.Write(buffer.Bytes())
	httpLogger.Write([]byte("\n"))
}

func (c *client) Do() ([]byte, error) {
	var err error

	// 日志记录
	start := time.Now()
	if c.OpenLog {
		defer c.log(start, &err)
	}

	// 创建请求
	var req *http.Request
	req, err = http.NewRequest(c.Method, c.Url, bytes.NewBuffer(c.ReqBody))
	if err != nil {
		return nil, err
	}

	for k, v := range c.Headers {
		req.Header.Set(k, v)
	}

	// 发送请求
	client := &http.Client{Timeout: c.Timeout}
	var resp *http.Response
	resp, err = client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// 读取结果
	if resp.StatusCode == http.StatusOK {
		c.RespBody, err = io.ReadAll(resp.Body)
	} else {
		err = fmt.Errorf("%d %s", resp.StatusCode, http.StatusText(resp.StatusCode))
	}

	return c.RespBody, err
}

// 解析结果。 可以是结构体指针，也可以是 map[string]any 指针
func (c *client) DoAndUnmarshal(to interface{}) error {
	c.Do()
	return json.Unmarshal(c.RespBody, to)
}
