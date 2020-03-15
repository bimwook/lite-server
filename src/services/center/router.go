package center

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"../../server"
	"../../woo"
)

//Start 初始化数据服务
func Start() {
	center := oCenter{}
	chSave := make(chan *Item)
	http.HandleFunc("/center/save.do", func(w http.ResponseWriter, r *http.Request) {
		method := r.Method
		if method == "POST" {
			var item Item
			item.Rowid = woo.NewSerial()
			item.Name = r.FormValue("name")
			item.Dbase = r.FormValue("dbase")
			item.Remark = r.FormValue("remark")
			item.Data = r.FormValue("data")
			chSave <- &item
			w.Header().Set("Server", server.ServerName)
			io.WriteString(w, `{"ret":true,"rowid":"`+item.Rowid+`"}`)
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
			hcount := center.Hash(name, dbase, data)
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
			item := center.Load(name, dbase, rowid)
			w.Header().Set("Server", server.ServerName)
			io.WriteString(w, "# "+item.Remark+"\r\n")
			io.WriteString(w, item.Data)
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
		fmt.Println("[Service] Center: Closing...")
	}()
	go func() {
		for !server.IsTerminated {
			select {
			case item := <-chSave:
				{
					center.Save(item)
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
