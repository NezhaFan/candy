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
	// 日志输出
	httpLogger io.Writer = os.Stdout
	// 默认超时时间
	httpTimeout = 10 * time.Second
)

func SetHTTPLogger(logger io.Writer) {
	httpLogger = logger
}

func SetHTTPTimeout(timeout time.Duration) {
	httpTimeout = timeout
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
	openLog bool   `json:"-"`
	ErrMsg  string `json:"ERROR,omitempty"`
}

func newClient() *client {
	c := &client{}
	c.Timeout = httpTimeout
	c.Headers = make(map[string]string)
	c.openLog = true
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

// params参数只接受 结构体、结构体指针、map[string]string类型
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

// params参数只接受 结构体、结构体指针、map[string]string类型
func HTTPGet(api string, params any) *client {
	c := newClient()
	c.Url = api
	c.Method = http.MethodGet
	if params != nil {
		// 防止一些自定义的参数已经拼接
		if !strings.Contains(c.Url, "?") {
			c.Url += "?"
		} else {
			c.Url += "&"
		}
		c.Url += string(buildQuery(params))
	}

	return c
}

func buildQuery(params interface{}) []byte {
	if params == nil {
		return nil
	}

	var buff bytes.Buffer
	buff.Grow(512)

	// map[string]string
	if data, ok := params.(map[string]string); ok {
		for k, v := range data {
			buff.WriteString(k + "=" + v + "&")
		}
		return buff.Bytes()[:buff.Len()-1]
	}

	// 反射判断类型
	tp := reflect.TypeOf(params)
	val := reflect.ValueOf(params)

	// 指针转一下
	if tp.Kind() == reflect.Pointer {
		tp = tp.Elem()
		val = val.Elem()
	}

	// 非结构体不通过
	if tp.Kind() != reflect.Struct {
		panic("参数类型不支持：" + tp.String())
	}

	for i := 0; i < tp.NumField(); i++ {
		k := tp.Field(i).Tag.Get("json")
		if k != "" {
			buff.WriteString(k + "=" + fmt.Sprintf("%v", val.Field(i).Interface()) + "&")
		}
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
	c.openLog = false
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
	if c.openLog {
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

// 解析结果。结构体指针 或 map[string]any指针
func (c *client) DoAndUnmarshal(to interface{}) error {
	c.Do()

	return json.Unmarshal(c.RespBody, to)
}
