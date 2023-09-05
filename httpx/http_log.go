package httpx

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"time"
)

var (
	// 日志输出
	httpLogger io.Writer = os.Stdout
)

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

func printLog(start time.Time, x *httpx) {
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
