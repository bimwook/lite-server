package main

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	"./services"
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
	services.ResetMailCache()
	chSave := make(chan *services.MailItem)
	chRemove := make(chan string)
	http.HandleFunc("/mail/save.do", func(w http.ResponseWriter, r *http.Request) {
		method := r.Method
		if method == "POST" {
			var item services.MailItem
			rowid := woo.NewSerial()
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
			io.WriteString(w, services.MailPeek(module, receiver))
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
			rowid, data := services.MailReceive(module, receiver)
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
			if woo.IsDone {
				chExit <- woo.IsDone
				break
			}
			time.Sleep(10)
		}
	}()
	go func() {
		for !woo.IsDone {
			select {
			case item := <-chSave:
				{
					services.MailSave(item.Rowid, item.Module, item.Sender, item.Receiver, item.Data)
				}
			case rowid := <-chRemove:
				{
					services.MailRemove(rowid)
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

func buildCenterSvr() {
	chSave := make(chan *services.CenterItem)
	http.HandleFunc("/center/save.do", func(w http.ResponseWriter, r *http.Request) {
		method := r.Method
		if method == "POST" {
			var item services.CenterItem
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
			hcount := services.CenterHash(name, dbase, data)
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
			remark, data := services.CenterLoad(name, dbase, rowid)
			w.Header().Set("Server", servername)
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
			if woo.IsDone {
				chExit <- woo.IsDone
				break
			}
			time.Sleep(10)
		}
	}()
	go func() {
		for !woo.IsDone {
			select {
			case item := <-chSave:
				{
					services.CenterSave(item.Name, item.Dbase, item.Remark, item.Data)
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

func main() {
	startAt := woo.Now()
	woo.ResetMain()
	http.Handle("/www/", http.FileServer(http.Dir(".")))
	http.HandleFunc("/now.do", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Server", servername)
		w.Write([]byte(woo.Now()))
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
		io.WriteString(w, "  - www\r\n")
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

	listening := woo.GetValue("server.listen-on")
	if listening == "" {
		woo.SetValue("server.listen-on", "11100")
	}
	port := strToInt(listening, 11100)

	fmt.Println(" ")
	fmt.Println(" Lite-Server: http://127.0.0.1:" + strconv.Itoa(port))
	fmt.Println("------------------------------------------------------------------")
	fmt.Println(" - EMAIL: woo@omuen.com")
	fmt.Println(" - SERVER-UUID: " + woo.GetServerSerial())
	fmt.Println("------------------------------------------------------------------")
	fmt.Println(" Start-At: " + startAt)

	listener, _ := net.Listen("tcp", ":"+strconv.Itoa(port))
	go func() {
		http.Serve(listener, nil)
	}()

	chExit := make(chan os.Signal, 1)
	signal.Notify(chExit, os.Interrupt, os.Kill)
	<-chExit

	woo.IsDone = true
	listener.Close()
	fmt.Println("Waiting for 1s to Close...")
	time.Sleep(time.Duration(1) * time.Second)
}
