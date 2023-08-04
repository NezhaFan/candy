### Golang工具包
用到泛型，版本 >= 1.18


### 类型转换
`BytesToString([]byte) string` bytes转字符串
`StringToBytes(string) []byte` 字符串转bytes
`FloatToString(float64) string` 浮点数转字符串
`Round(float64, int) float64` 将float64保留x位小数
`Sha256(string) string` SHA256加密
`HmacSha256(s string, key []byte) string` SHA256加密
`Base64([]byte) string` base64
`IPToInt(string) int` IP 点分十进制 => 十进制
`IntToIP(int) string` IP 十进制 => 点分十进制
`RuneLength(string) int` 字符串长度
`IsEmpty(comparable) bool` 是否为空
`Max(Integer, ...Integer)` 多项中的最大值
`Min(Integer, ...Integer)` 多项中的最小值
`Div(Integer, Integer, int)` 做除法，保留x位小数


### 时间操作
> 对time.Time再封装为Timer，并非要取代time.Time，主要用于日期和月份的操作。

`SetDefaultLocation(string)` 设置Timer的默认时区（默认本地时区）
`type Timer`
- `NowTimer() Timer` 当前时间
- `StringTimer(string) Timer` 从字符串格式化
- `UnixTimer(int64) Timer` 从时间戳格式化
- `Add(time.Duration) Timer` 同`time.Time.Add()`
- `AddDay(i int) Timer` 加天数
- `AddMonth(i int) Timer` 加月数
- `AddYear(i int) Timer` 加年数
- `DayStart() Timer` 当日开始
- `DayEnd() Timer` 当日结束 (当日最后1秒，非第二日0时)
- `MonthStart() Timer` 当月开始
- `MonthEnd() Timer` 当月结束 (当月最后1秒，非第二月0时)
- `String() string` 转化回时间字符串
- `Time() time.Time` 转化回时间
- `Unix() int64` 转化回时间戳

```go
// 链式操作

now := candy.NowTimer()

// 下个月的时间范围
next := now.AddMonth(1)
fmt.Println(next.MonthStart().String())
fmt.Println(next.MonthEnd().String())

// 今天的时间范围
fmt.Println(now.DayStart().String())
fmt.Println(now.DayEnd().String())

```


### timex 包
```go

// 生成Timex，参数格式化int64或string的给定时间，其他类型返回当前时间
func New(interface{}) Timex

// 事件操作
func (Timex).Add(time.Duration) Timex
func (Timex).AddDay(int) Timex
func (Timex).AddMonth(int) Timex
func (Timex).AddYear(int) Timex
// 当日
func (Timex).DayEnd() Timex
func (Timex).DayStart() Timex
// 当月
func (Timex).MonthEnd() Timex
func (Timex).MonthStart() Timex
// 当周 (周一到周日)
func (Timex).WeekEnd() Timex
func (Timex).WeekStart() Timex

// 最终转化
func (Timex).Time() time.Time
func (Timex).String() string
func (Timex).Unix() int64
```


### HTTP请求

```go

// 设置日志
func SetLogger(w io.Writer)

```

`HTTPPostJSON(url string, params interface{}) *client` POST请求，json编码
`HTTPPostForm(url string, params map[string]string) *client` POST请求，x-www-form-urlencoded编码
`HTTPGet(url string, params map[string]string) *client` GET请求
`SetHTTPLogger(io.Writer)` 设置日志输出
`SetHTTPTimeout(timeout time.Duration)` 设置默认超时时间 (默认10秒)

`type client`
- `AddHeader(key, value string)` 添加header
- `SetTimeout(time.Duration)` 设置此请求超时时间 
- `CloseLog()` 此请求关闭日志
- `Do() ([]byte, error)` 发出请求
- `DoAndUnmarshal(to interface{}) error` 发出请求，并解析json结果到

```go
// 链式操作
var req struct {
  Id   int    `json:"id"`
  Name string `json:"name"`
}
req.Id = 100
req.Name = "Alice"

var resp struct {
  Status int    `json:"status"`
  Msg    string `json:"msg"`
  Data   interface{}    `json:"data"`
}

// 简单使用
candy.HTTPPostJSON("https://kikia.cc/api/test", &req).DoAndUnmarshal(&resp)
fmt.Println("返回结果：", resp.Status, resp.Msg, resp.Data)

// 参数 也可以用 map[string]string 。 但是禁用map[string]interface{}.
// 返回 也可以用 map[string]interface{} 接收
req2 := map[string]string{"id": "100", "name": "Alice"}
var resp2 map[string]interface{}
candy.HTTPPostJSON("https://kikia.cc/api/test", &req2).DoAndUnmarshal(&resp2)
fmt.Println("返回结果：", resp2)

// 灵活运用
client := candy.HTTPGet("https://kikia.cc/api/test", map[string]string{"Id": "123"})
// 更多的超时时间，关闭日志记录
client.SetTimeout(time.Second*3).CloseLog().AddHeader("Auth", "Bearer xxx")
// 自己解析 []byte
b, err := client.Do()
if err != nil {
  fmt.Println("ERROR:", err)
} else {
  fmt.Println("返回结果", string(b))
}
```