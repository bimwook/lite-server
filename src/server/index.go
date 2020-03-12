package server

import (
	"io"
	"net/http"

	"../woo"
)

//Index 默认页
func Index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Server", ServerName)
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
}

//Now 服务器时间
func Now(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Server", ServerName)
	w.Write([]byte(woo.Now()))
}
