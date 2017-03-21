//BUG - After the client exit, the number in center will not reduce.
//TODO - Add a onclose event in frontend.

// USE MANY MANY TINY THEARDS TO SEND OR RECEIVE EVENTS.

package main

import "fmt"
import "net/http"
import "github.com/gorilla/websocket"

//import "strconv"

type Msg struct {
}

//type Node maps to a client.
type Node struct {
	c_ptr *Center         // a pointer to center.
	conn  *websocket.Conn //connent client to node
}

//use go statment to call this func
func (n *Node) Run(ifexit chan<- bool) {
	//var err error
	var if_listener_exit = make(chan bool)
	fmt.Println("node::Run()")
	go func() {
		//listener
		//var msg_type int
		var msg_cx []byte
		var err error
		for {
			//the code will be blocked here:
			//but don't worry, becase it's in the go statment.
			_, msg_cx, err = n.conn.ReadMessage()
			if err != nil {
				//the client was closed.
				fmt.Println("-close-client-")
				if_listener_exit <- true
				return
			}
			//check the content that client sent, and push it to center.
			if string(msg_cx[:]) == "_NEW_CLIENT_" {
				fmt.Println("-new-client-")
				//code for pushing goes here...
			} else {
				//other message...
				fmt.Println("received msg from client:", string(msg_cx[:]))
				//code for pushing goes here...
			}
		}
	}()
	go func() {
		//responser
		//fetch the msg from center, and send it to client.
	}()
	select {
	case <-if_listener_exit:
		//if the listener exit, then the whole node will exit.
		ifexit <- true
		fmt.Println("node::run() -close-node-")
		return
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
		fmt.Println(
			"fatal - center::newNode -",
			"when creating the websocket.conn",
		)
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