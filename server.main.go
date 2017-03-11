package main

import (
	//"encoding/json"
	"fmt"
	"log"
	//"io"
	//"io/ioutil"
	"github.com/gorilla/websocket"
	"net/http"
	//"strings"
)

var dialogs []string

var upgdr = websocket.Upgrader{}

func WebSocketServFunc(w http.ResponseWriter, r *http.Request) {
	musubi, err := upgdr.Upgrade(w, r, nil)
	if err != nil {
		log.Print("Upgrade Err: ", err)
		return
	}
	defer musubi.Close()
	for {
		msg_type, msg_cx, err := musubi.ReadMessage()
		if err != nil {
			log.Print("Read Msg Err: ", err)
			break
		}
		err = musubi.WriteMessage(msg_type, msg_cx)
		if err != nil {
			log.Print("Write Msg Err: ", err)
			break
		}
	}
}

func ServeHome() func(http.ResponseWriter, *http.Request) {
	return (func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.Error(w, "File Not Found", 404)
			return
		}
		if r.Method != "GET" {
			http.Error(w, "Bad Access Method", 405)
			return
		}
		http.ServeFile(w, r, "index.html")
	})
}

func ServeFile(filename string) func(http.ResponseWriter, *http.Request) {
	return (func(w http.ResponseWriter, r *http.Request) {
	})
}

func main() {
	/* File Server */
	http.HandleFunc("/", ServeHome())
	http.HandleFunc("/index.html", ServeFile("index.html"))
	http.HandleFunc("/app.js", ServeFile("app.js"))
	http.HandleFunc("/api.js", ServeFile("api.js"))
	http.HandleFunc("/ws", WebSocketServFunc)

	log.Fatal(http.ListenAndServe(":8888", nil))
	fmt.Println("vim-go")
}
