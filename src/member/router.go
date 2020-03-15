package member

import (
	"io"
	"net/http"
	"strconv"

	"../server"
)

//Actions 成员管理
var Actions IMember

//Start 初始化邮件服务
func Start() {
	member := oMember{maindb: "./dbase/member.db"}
	member.Start()
	Actions := member.GetActions()
	http.HandleFunc("/member/new.do", func(w http.ResponseWriter, r *http.Request) {
		method := r.Method
		if method == "POST" {
			var item Item
			item.Level = 10
			item.Name = r.FormValue("name")
			item.Secret = r.FormValue("secret")
			data, ok := Actions.New(item)
			w.Header().Set("Server", server.ServerName)
			io.WriteString(w, `{"ret": `+strconv.FormatBool(ok)+`, "data":"`+data+`"}`)
		} else {
			io.WriteString(w, "Method: POST;\r\n")
			io.WriteString(w, "Parametes: name, secret;\r\n")
		}
	})
	http.HandleFunc("/member/renew.do", func(w http.ResponseWriter, r *http.Request) {
		method := r.Method
		if method == "POST" {
			//var item Item
			name := r.FormValue("name")
			secret := r.FormValue("secret")
			data, ok := Actions.Renew(name, secret)
			w.Header().Set("Server", server.ServerName)
			io.WriteString(w, `{"ret": `+strconv.FormatBool(ok)+`, "secret":"`+data+`"}`)
		} else {
			io.WriteString(w, "Method: POST;\r\n")
			io.WriteString(w, "Parametes: name, secret;\r\n")
		}
	})
	http.HandleFunc("/member/remove.do", func(w http.ResponseWriter, r *http.Request) {
		method := r.Method
		if method == "POST" {
			//var item Item
			name := r.FormValue("name")
			secret := r.FormValue("secret")
			data, ok := Actions.Remove(name, secret)
			w.Header().Set("Server", server.ServerName)
			io.WriteString(w, `{"ret": `+strconv.FormatBool(ok)+`, "secret":"`+data+`"}`)
		} else {
			io.WriteString(w, "Method: POST;\r\n")
			io.WriteString(w, "Parametes: name, secret;\r\n")
		}
	})
}
