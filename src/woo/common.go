package woo

import (
	"crypto/sha256"
	"encoding/base64"
	"math/rand"
	"strconv"
	"time"
)

//StrToInt 字符串转整形
func StrToInt(s string, def int) int {
	ret, err := strconv.Atoi(s)
	if err != nil {
		return def
	}
	return ret
}

//Now 当前日期及时间
func Now() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

//NewSerial 返回一个随机字符串
func NewSerial() string {
	ret := time.Now().Format("20060102")
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 40; i++ {
		ret = ret + strconv.Itoa(r.Intn(10))
	}
	return ret
}

//Sha256 Base64-Hash
func Sha256(data string) string {
	hm := sha256.New()
	hm.Write([]byte(data))
	return base64.StdEncoding.EncodeToString(hm.Sum(nil))
}
