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


### HTTP请求

`HTTPGet(url string, params map[string]string) *client` GET请求
`HTTPPostJSON(url string, params any) *client` POST请求，json编码
`HTTPPostForm(url string, params map[string]string) *client` POST请求，x-www-form-urlencoded编码

`type client`



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