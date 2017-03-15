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

type Node struct {
	msgs   []string
	center Center
	w      http.ResponseWriter
	r      *http.Request
	stat   int8
}

func newNode(center Center, w http.ResponseWriter, r *http.Request) *Node {
	var res = new(Node)
	res.center = center
	res.w = w
	res.r = r
	return res
}

type Center struct {
	nodes []*Node
}

func newCenter() *Center {
	return new(Center)
}
func (c *Center) AddNode(new_node *Node) error {
	/* TODO - considerate possble error */
	c.nodes = append(c.nodes, new_node)
	return nil
}

func (c *Center) Boardcast(msg string) {
}

func (c *Center) Loop() {
}

func WebSocketServFunc(center Center, w http.ResponseWriter, r *http.Request) {
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
	var center = newCenter()
	go center.Loop()
	http.HandleFunc("/", ServeHome())
	http.HandleFunc("/index.html", ServeFile("index.html"))
	http.HandleFunc("/app.js", ServeFile("app.js"))
	http.HandleFunc("/api.js", ServeFile("api.js"))

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		var node = newNode(center, w, r)
		center.Add(node)
	})

	log.Fatal(http.ListenAndServe(":8888", nil))
	fmt.Println("vim-go")
}
