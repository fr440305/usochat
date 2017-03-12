package main

import (
	//"encoding/json"
	"bytes"
	"fmt"
	"log"
	//"io"
	//"io/ioutil"
	"github.com/gorilla/websocket"
	"net/http"
	//"strings"
)

type Client struct {
	msgs []string
	stat int8
}

func newClient() *Client {
	return new(Client)
}

type Center struct {
	web_clients []*Client
}

func newCenter() *Center {
	return new(Center)
}

var dialogs []string = *new([]string)

func WebSocketServFunc(w http.ResponseWriter, r *http.Request) {
	var upgdr = websocket.Upgrader{}
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
			return
		}
		if len(msg_cx) != 0 {
			dialogs = append(dialogs, (*bytes.NewBuffer(msg_cx)).String())
			fmt.Println(dialogs)
		}
		err = musubi.WriteMessage(msg_type, msg_cx)
		if err != nil {
			log.Print("Write Msg Err: ", err)
			return
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
