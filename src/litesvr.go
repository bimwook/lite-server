package main

import (
	"fmt"
	"io"
	"muen"
	"net/http"
	"strconv"

	"./woo"
)

const servername = "Lite-Server/1.1"

func strToInt(s string, def int) int {
	ret, err := strconv.Atoi(s)
	if err != nil {
		return def
	}
	return ret
}

func buildMailSvr() {
	woo.ResetMailCache()
	chSave := make(chan *woo.MailItem)
	chRemove := make(chan string)
	http.HandleFunc("/mail/save.do", func(w http.ResponseWriter, r *http.Request) {
		method := r.Method
		if method == "POST" {
			var item woo.MailItem
			rowid := muen.NewKey()
			item.Rowid = rowid
			item.Module = r.FormValue("module")
			item.Sender = r.FormValue("sender")
			item.Receiver = r.FormValue("receiver")
			item.Data = r.FormValue("data")
			chSave <- &item
			w.Header().Set("Server", servername)
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
			w.Header().Set("Server", servername)
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			io.WriteString(w, woo.MailPeek(module, receiver))
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
			w.Header().Set("Server", servername)
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			rowid, data := woo.MailReceive(module, receiver)
			io.WriteString(w, data)
			if rowid != "" {
				chRemove <- rowid
			}
		} else {
			io.WriteString(w, "Method: POST;\r\n")
			io.WriteString(w, "Parametes: module, receiver;\r\n")
		}
	})

	go func() {
		for {
			select {
			case item := <-chSave:
				{
					woo.MailSave(item.Rowid, item.Module, item.Sender, item.Receiver, item.Data)
				}
			case rowid := <-chRemove:
				{
					woo.MailRemove(rowid)
				}
			}
		}
	}()
}

func buildCenterSvr() {
	chSave := make(chan *woo.CenterItem)
	http.HandleFunc("/center/save.do", func(w http.ResponseWriter, r *http.Request) {
		method := r.Method
		if method == "POST" {
			var item woo.CenterItem
			item.Name = r.FormValue("name")
			item.Dbase = r.FormValue("dbase")
			item.Remark = r.FormValue("remark")
			item.Data = r.FormValue("data")
			chSave <- &item
			w.Header().Set("Server", servername)
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
			hcount := woo.CenterHash(name, dbase, data)
			w.Header().Set("Server", servername)
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
			remark, data := woo.CenterLoad(name, dbase, rowid)
			w.Header().Set("Server", servername)
			io.WriteString(w, "# "+remark+"\r\n")
			io.WriteString(w, data)
		} else {
			io.WriteString(w, "//Method: POST;\r\n")
			io.WriteString(w, "//Parametes: name, dbase, rowid;\r\n")
		}
	})

	go func() {
		for {
			select {
			case item := <-chSave:
				{
					woo.CenterSave(item.Name, item.Dbase, item.Remark, item.Data)
				}
			}
		}
	}()
}

func main() {
	startAt := muen.Now()
	woo.ResetMain()
	http.HandleFunc("/now.do", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Server", servername)
		w.Write([]byte(muen.Now()))
	})
	buildMailSvr()
	buildCenterSvr()
	http.HandleFunc("/about.do", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Server", servername)
		io.WriteString(w, woo.AboutMe())
		io.WriteString(w, "  Start-At: "+startAt)
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Server", servername)
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		io.WriteString(w, "  Lite-Server is working...\r\n")
		io.WriteString(w, "--------------------------------------\r\n")
		io.WriteString(w, "  - now.do\r\n")
		io.WriteString(w, "  - center\r\n")
		io.WriteString(w, "    - save.do\r\n")
		io.WriteString(w, "    - hash.do\r\n")
		io.WriteString(w, "    - load.do\r\n")
		io.WriteString(w, "  - mail\r\n")
		io.WriteString(w, "    - save.do\r\n")
		io.WriteString(w, "    - peek.do\r\n")
		io.WriteString(w, "    - receive.do\r\n")
		io.WriteString(w, "  - about.do\r\n")
	})

	fmt.Println(" ")
	fmt.Println(" Lite-Server: http://127.0.0.1:11100")
	fmt.Println("------------------------------------------------------------------")
	fmt.Println(" - EMAIL: woo@omuen.com")
	fmt.Println(" - SERVER-UUID: " + woo.GetServerSerial())
	fmt.Println("------------------------------------------------------------------")
	fmt.Println(" Start-At: " + startAt)
	muen.Send("<!>\r\n")
	muen.Sendln("[Lite-Server] Start at " + startAt)
	http.ListenAndServe(":11100", nil)
}
