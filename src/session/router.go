package session

import (
	"io"
	"net/http"
	"strconv"

	"../server"
)

//Actions 管理
var Actions ISession

//Start 初始化邮件服务
func Start() {
	var session oSession
	session.Start()
	Actions = session.GetActions()
	http.HandleFunc("/session/checkin.do", func(w http.ResponseWriter, r *http.Request) {
		method := r.Method
		if method == "POST" {
			//var item Item
			name := r.FormValue("name")
			secret := r.FormValue("secret")
			token, ok := session.CheckIn(name, secret)
			w.Header().Set("Server", server.ServerName)
			io.WriteString(w, `{"ret": `+strconv.FormatBool(ok)+`, "token":"`+token+`"}`)
		} else {
			io.WriteString(w, "Method: POST;\r\n")
			io.WriteString(w, "Parametes: name, secret;\r\n")
		}
	})
	http.HandleFunc("/session/checkout.do", func(w http.ResponseWriter, r *http.Request) {
		method := r.Method
		if method == "POST" {
			name := r.FormValue("name")
			token := r.FormValue("token")
			ok := session.CheckOut(name, token)
			w.Header().Set("Server", server.ServerName)
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")

			io.WriteString(w, `{"ret": `+strconv.FormatBool(ok)+`}`)
		} else {
			io.WriteString(w, "Method: POST;\r\n")
			io.WriteString(w, "Parametes: name, token;\r\n")
		}
	})
	http.HandleFunc("/session/check.do", func(w http.ResponseWriter, r *http.Request) {
		method := r.Method
		if method == "POST" {
			name := r.FormValue("name")
			token := r.FormValue("token")
			w.Header().Set("Server", server.ServerName)
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			ok := session.Check(name, token)
			io.WriteString(w, `{"ret": `+strconv.FormatBool(ok)+`}`)
		} else {
			io.WriteString(w, "Method: POST;\r\n")
			io.WriteString(w, "Parametes: name, token;\r\n")
		}
	})
}
