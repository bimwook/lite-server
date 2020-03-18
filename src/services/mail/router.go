package mail

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"../../server"
	"../../woo"
)

//Actions 消息服务
var Actions IMail

//Start 初始化邮件服务
func Start() {
	mail := oMail{maindb: "./dbase/mail.db"}

	mail.Start()
	Actions = mail.GetActions()
	chSave := make(chan *Item)
	chRemove := make(chan string)
	http.HandleFunc("/mail/save.do", func(w http.ResponseWriter, r *http.Request) {
		method := r.Method
		if method == "POST" {
			var item Item
			rowid := woo.NewSerial()
			item.Rowid = rowid
			item.Module = r.FormValue("module")
			item.Sender = r.FormValue("sender")
			item.Receiver = r.FormValue("receiver")
			item.Data = r.FormValue("data")
			chSave <- &item
			w.Header().Set("Server", server.ServerName)
			io.WriteString(w, `{"ret": true, "rowid":"`+rowid+`"}`)
		} else {
			io.WriteString(w, "Method: POST;\r\n")
			io.WriteString(w, "Parametes: module, sender, receiver, data;\r\n")
		}
	})
	http.HandleFunc("/mail/peek.do", func(w http.ResponseWriter, r *http.Request) {
		method := r.Method
		if method == "POST" {
			module := r.FormValue("module")
			receiver := r.FormValue("receiver")
			w.Header().Set("Server", server.ServerName)
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			io.WriteString(w, mail.Peek(module, receiver))
		} else {
			io.WriteString(w, "Method: POST;\r\n")
			io.WriteString(w, "Parametes: module, receiver;\r\n")
		}
	})
	http.HandleFunc("/mail/receive.do", func(w http.ResponseWriter, r *http.Request) {
		method := r.Method
		if method == "POST" {
			module := r.FormValue("module")
			receiver := r.FormValue("receiver")
			w.Header().Set("Server", server.ServerName)
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			rowid, data := mail.Receive(module, receiver)
			io.WriteString(w, data)
			if rowid != "" {
				chRemove <- rowid
			}
		} else {
			io.WriteString(w, "Method: POST;\r\n")
			io.WriteString(w, "Parametes: module, receiver;\r\n")
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
		fmt.Println("[Service] Mail: Closing...")
	}()
	go func() {
		for !server.IsTerminated {
			select {
			case item := <-chSave:
				{
					mail.Save(item)
				}
			case rowid := <-chRemove:
				{
					mail.Remove(rowid)
				}
			case <-chExit:
				{
					break
				}
			}
		}
		fmt.Println("[Service] Mail: Closed.")
	}()
}
