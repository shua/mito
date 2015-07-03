package main

import (
	"flag"
	"log"
	"net/http"
	"text/template"
)

var addr = flag.String("port", ":8080", "http port")
var mainTempl = template.Must(template.ParseFiles("mito.html"))

func serveMain(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "Not found", 404)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not supported", 405)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	mainTempl.Execute(w, r.Host)
}

func serveWs(w http.ResponseWriter, r *http.Request) {
    if r.Method != "GET" {
        http.Error(w, "Method not supported", 405)
        return
    }
    ws, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Println(err)
        return
    }
    c := &connection{ send: make(chan []byte, 256), ws:ws }
    l.join <- c
    go c.writePump()
    c.readPump()
}

func main() {
	flag.Parse()
        go l.run()
	http.HandleFunc("/", serveMain)
        http.HandleFunc("/ws", serveWs)
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
