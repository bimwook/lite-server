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

	"./server"
	"./services"
	"./session"
	"./woo"
)

func main() {
	startAt := woo.Now()
	server.ResetMain()
	session.Reset()
	http.HandleFunc("/", server.Index)

	http.Handle("/www/", http.FileServer(http.Dir(".")))
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

	listener, _ := net.Listen("tcp", ":"+strconv.Itoa(port))
	go func() {
		http.Serve(listener, nil)
	}()

	chExit := make(chan os.Signal, 1)
	signal.Notify(chExit, os.Interrupt, os.Kill)
	<-chExit

	server.IsTerminated = true
	listener.Close()
	fmt.Println("Waiting for 1s to Close...")
	time.Sleep(time.Duration(1) * time.Second)
}
