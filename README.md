### Golang工具包
版本 >= 1.18

### 类型转换
`BytesToString([]byte) string` bytes转字符串
`StringToBytes(string) []byte` 字符串转bytes
`FloatToString(float64) string` 浮点数转字符串
`Round(n float64, i int) float64` n保留i位小数
`Sha256(string) string` SHA256加密
`HmacSha256(s string, key []byte) string` SHA256加密
`Base64([]byte) string` base64
`IPToInt(string) int` IP 点分十进制 => 十进制
`IntToIP(int) string` IP 十进制 => 点分十进制
`RuneLength(string) int` 字符串长度
`IsEmpty(comparable) bool` 是否为空
`Max(Integer, Integer)` 两者最大值
`Max(Integer, Integer)` 两者最大值
