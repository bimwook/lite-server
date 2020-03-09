package muen

import (
	"bufio"
	"io"
	"math/rand"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

//Error 错误
type Error struct {
	Code int
	Info string
}

//ShowError 显示一个错误
func (o Error) ShowError() string {
	return "[" + strconv.Itoa(o.Code) + "] " + o.Info
}

//Now 返回格式化的当前日期及时间
func Now() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

//Rndid 返回一个随机ID
func Rndid() string {
	return time.Now().Format("20060102150405")
}

//Send 发送调试信息
func Send(str string) bool {
	socket, err := net.Dial("udp4", "127.0.0.1:19800")
	if err != nil {
		return false
	}
	io.WriteString(socket, str)
	defer socket.Close()
	return true
}

//Sendln 发送带换行的调试信息
func Sendln(str string) bool {
	socket, err := net.Dial("udp4", "127.0.0.1:19800")
	if err != nil {
		return false
	}
	io.WriteString(socket, "["+Now()+"] "+str+"\r\n")
	defer socket.Close()
	return true
}

//NewKey 返回一个随机字符串
func NewKey() string {
	ret := time.Now().Format("20060102")
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 40; i++ {
		ret = ret + strconv.Itoa(r.Intn(10))
	}
	return ret
}

//SubString 截取字符串
func SubString(str string, begin, length int) string {
	s := []rune(str)
	size := len(s)
	if begin < 0 {
		begin = 0
	}
	if begin >= size {
		begin = size
	}
	end := begin + length
	if end > size {
		end = size
	}
	return string(s[begin:end])
}

//LoadMap 加载一个Map
func LoadMap(fn string) (ret map[string]string, ok bool) {
	ok = true
	file, e := os.Open(fn) // For read access.
	if e != nil {
		ok = false
		return nil, ok
	}
	ret = make(map[string]string)
	reader := bufio.NewReader(file)
	for {
		r, e := reader.ReadString('\n')
		line := strings.Replace(strings.Replace(r, "\r", "", -1), "\n", "", -1)
		if e != nil {
			break
		}
		p := strings.Index(line, "=")
		if p != -1 {
			key := SubString(line, 0, p)
			value := SubString(line, p+1, len(line))
			ret[key] = value
		}
	}
	defer file.Close()

	return ret, ok
}

//WebEncode 编码
func WebEncode(s string) string {
	ret := strings.Replace(s, "&", "&amp;", -1)
	ret = strings.Replace(ret, "<", "&lt;", -1)
	ret = strings.Replace(ret, ">", "&gt;", -1)
	return ret
}

//Root 结构
type Root struct {
	Name    string
	Version string
	Tag     int
}

//Now 当前日期及时间
func (o Root) Now() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

//Log 日志
func (o Root) Log(str string) string {
	return "[" + o.Now() + "] " + o.Name + ": " + str
}

//SetName 设置名称
func (o *Root) SetName(name string) bool {
	o.Name = name
	return true
}

//SetVersion 设置版本号
func (o *Root) SetVersion(version string) bool {
	o.Version = version
	return true
}
