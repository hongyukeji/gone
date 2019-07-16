package _conv

import "unsafe"

//[]byte转string
func BytesToStr(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

//[]byte转换为字符串
func BytesMustStr(value []byte, defaults ...string) string {
	if len(value) > 0 {
		return BytesToStr(value)
	}
	if len(defaults) > 0 {
		return defaults[0]
	}
	return ""
}
