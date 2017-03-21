//BUG - After the client exit, the number in center will not reduce.
//TODO - Add a onclose event in frontend.

// USE MANY MANY TINY THEARDS TO SEND OR RECEIVE EVENTS.

package main

import "fmt"
import "net/http"
import "github.com/gorilla/websocket"

//import "strconv"

type Msg map[string][]byte

//type Node maps to a client.
type Node struct {
	c_ptr *Center         // a pointer to center.
	conn  *websocket.Conn //connent client to node
}

//use go statment to call this func
func (n *Node) Run(ifexit chan<- bool) {
	var err error
	defer func() { ifexit <- true }()
	fmt.Println("node::Run()")
	err = n.conn.WriteMessage(websocket.TextMessage, []byte{'h', 'e', 'l', 'l', 'o', ' ', 'c', 'l', 'i', 'e', 'n', 't'})
	if err != nil {
		fmt.Println("fatal - node::Run() - cannot write msg to cilent")
	}
	for {
	}
}

type Center struct {
	msg_queue   chan Msg
	nodes       []*Node
	upgrader    websocket.Upgrader //Constant
	num_onliner int                //just for test.
}

func newCenter() *Center {
	fmt.Println("newCenter()")
	var res = new(Center)
	res.msg_queue = make(chan Msg)
	res.nodes = *new([]*Node)
	fmt.Println(res.nodes)
	res.upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	return res
}

func (c *Center) newNode(w http.ResponseWriter, r *http.Request) *Node {
	var err error
	fmt.Println("center::newNode()")
	var res = new(Node)
	res.c_ptr = c
	res.conn, err = c.upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("fatal - center::newNode - when creating the websocket.conn")
	}
	c.nodes = append(c.nodes, res)
	fmt.Println("online: ", c.GetOnliner())
	return res
}

//listen and handle the msg.
//use go statment to call this func.
func (c *Center) Run() {
	for {
		select {
		//case msg := <-c.msg_queue:
		//...//
		default:
		}
	}
}

//return the number of people online:
func (c *Center) GetOnliner() int {
	fmt.Println("center::GetOnliner()", c.nodes)
	return len(c.nodes)
}

func main() {
	fmt.Println("http://127.0.0.1:9999")
	var center = newCenter()
	go center.Run()
	//To provide the webpages to the client:
	http.Handle("/", http.FileServer(http.Dir(".")))
	//To handle the websocket request:
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		var if_node_exit = make(chan bool)
		go center.newNode(w, r).Run(if_node_exit)
		select {
		case <-if_node_exit:
			fmt.Println("A node exit.")
			return
		}
	})
	http.ListenAndServe(":9999", nil)
}
