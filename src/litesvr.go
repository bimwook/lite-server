package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	"./server"
	"./services"
	"./session"
	"./woo"
)

func main() {
	startAt := woo.Now()
	server.Start()
	session.Start()
	http.HandleFunc("/", server.Index)
	os.MkdirAll("./www", os.ModePerm)
	os.MkdirAll("./var", os.ModePerm)
	http.Handle("/www/", http.FileServer(http.Dir(".")))
	http.Handle("/var/", http.FileServer(http.Dir(".")))
	http.HandleFunc("/now.do", server.Now)
	services.Start()
	http.HandleFunc("/about.do", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Server", server.ServerName)
		io.WriteString(w, server.AboutMe())
		io.WriteString(w, "  Start-At: "+startAt)
	})

	listening := server.GetValue("server.listen-on")
	if listening == "" {
		server.SetValue("server.listen-on", "11100")
	}
	port := woo.StrToInt(listening, 11100)

	fmt.Println(" ")
	fmt.Println(" Lite-Server: http://127.0.0.1:" + strconv.Itoa(port))
	fmt.Println("------------------------------------------------------------------")
	fmt.Println(" - EMAIL: woo@omuen.com")
	fmt.Println(" - SERVER-UUID: " + server.GetServerSerial())
	fmt.Println("------------------------------------------------------------------")
	fmt.Println(" Start-At: " + startAt)
	fmt.Println(" ")

	chExit := make(chan os.Signal, 1)
	signal.Notify(chExit, os.Interrupt, os.Kill)
	go func() {
		err := http.ListenAndServe(":"+strconv.Itoa(port), nil)
		if err != nil {
			fmt.Println("启动失败：端口被占用")
			chExit <- os.Interrupt
		}
	}()
	<-chExit
	server.IsTerminated = true
	//fmt.Println("Waiting for 3s to Close...")
	time.Sleep(time.Duration(3) * time.Second)
}
