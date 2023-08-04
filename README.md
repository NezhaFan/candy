### Golang工具包
用到泛型，版本 >= 1.18


### 通用转换 candy包
```go
BytesToString([]byte) string // bytes转字符串
StringToBytes(string) []byte // 字符串转bytes
FloatToString(float64) string // 浮点数转字符串
Round(float64, int) float64 // 将float64保留x位小数
Sha256(string) string // SHA256加密
HmacSha256(s string, key []byte) string // SHA256加密
Base64([]byte) string // base64
IPToInt(string) int // IP 点分十进制 => 十进制
IntToIP(int) string // IP 十进制 => 点分十进制
RuneLength(string) int // 字符串长度
IsEmpty(comparable) bool // 是否为空
Max(Integer, ...Integer) // 多项中的最大值
Min(Integer, ...Integer) // 多项中的最小值
Div(Integer, Integer, int) // 做除法，保留x位小数
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
