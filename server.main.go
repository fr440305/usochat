package main

import (
	"encoding/json"
	"fmt"
	"log"
	//"io"
	//"io/ioutil"
	"github.com/gorilla/websocket"
	"net/http"
	"strings"
)

var dialogs []string

func WebSocketServFunc(w http.ResponseWriter, r *http.Request) {
	musubi, err := websocket.Upgrader(w, r, nil)
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
