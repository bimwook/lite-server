package center

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"../../server"
)

//BuildSvr 初始化数据服务
func BuildSvr() {
	chSave := make(chan *Item)
	http.HandleFunc("/center/save.do", func(w http.ResponseWriter, r *http.Request) {
		method := r.Method
		if method == "POST" {
			var item Item
			item.Name = r.FormValue("name")
			item.Dbase = r.FormValue("dbase")
			item.Remark = r.FormValue("remark")
			item.Data = r.FormValue("data")
			chSave <- &item
			w.Header().Set("Server", server.ServerName)
			io.WriteString(w, `{"ret":true,"rowid":""}`)
		} else {
			io.WriteString(w, "//Method: POST;\r\n")
			io.WriteString(w, "//Parametes: name, dbase, remark, data;\r\n")
		}
	})
	http.HandleFunc("/center/hash.do", func(w http.ResponseWriter, r *http.Request) {
		method := r.Method
		if method == "POST" {
			name := r.FormValue("name")
			dbase := r.FormValue("dbase")
			data := r.FormValue("data")
			hcount := Hash(name, dbase, data)
			w.Header().Set("Server", server.ServerName)
			io.WriteString(w, `{"ret":true,"count":`+hcount+`}`)
		} else {
			io.WriteString(w, "//Method: POST;\r\n")
			io.WriteString(w, "//Parametes: name, dbase, data;\r\n")
		}
	})
	http.HandleFunc("/center/load.do", func(w http.ResponseWriter, r *http.Request) {
		method := r.Method
		if method == "POST" {
			name := r.FormValue("name")
			dbase := r.FormValue("dbase")
			rowid := r.FormValue("rowid")
			remark, data := Load(name, dbase, rowid)
			w.Header().Set("Server", server.ServerName)
			io.WriteString(w, "# "+remark+"\r\n")
			io.WriteString(w, data)
		} else {
			io.WriteString(w, "//Method: POST;\r\n")
			io.WriteString(w, "//Parametes: name, dbase, rowid;\r\n")
		}
	})

	chExit := make(chan bool)
	go func() {
		for {
			if server.IsTerminated {
				chExit <- server.IsTerminated
				break
			}
			time.Sleep(10)
		}
	}()
	go func() {
		for !server.IsTerminated {
			select {
			case item := <-chSave:
				{
					Save(item.Name, item.Dbase, item.Remark, item.Data)
				}
			case <-chExit:
				{
					break
				}
			}
		}
		fmt.Println("[Service] Center: Closed.")
	}()
}
