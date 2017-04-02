// USE MANY MANY TINY THEARDS TO SEND OR RECEIVE EVENTS.
// USE SLICE EVERYWHERE - DO NOT USE container/list.
// 不光要思考架构，还要思考架构的迭代与演化。
// 要牢记：软件是长出来的。

package main

import "fmt"
import "html"
import "github.com/gorilla/websocket"

//type Node maps to a client.
type Node struct {
	msg_from_center chan Msg
	c_ptr           *Center         // a pointer to center.
	conn            *websocket.Conn //connent client to node
	iden            string          // the identification for node.
}

func (N *Node) listenToUser(ifexit chan<- bool) {
	//listener
	//This goroutine receive msgs in the form of JSON from client.
	var err error
	var msg_cx []byte      // the byte array from user.
	var str_msg_cx string  // the conversion for byte array.
	var msg_to_center *Msg // the message that needs to send to center.
	for {
		//the code will be blocked here(conn.ReadMessage():
		//but don't worry, becase it's in the go statment.
		_, msg_cx, err = N.conn.ReadMessage()
		str_msg_cx = string(msg_cx[:])
		if err != nil {
			//the client was closed.
			msg_to_center = newMsg(N)
			msg_to_center.setDescription("user-logout")
			N.c_ptr.msg_queue <- *msg_to_center
			//FIXME - do not remove my self here.
			//make center to do this.
			N.c_ptr.removeNode(N)
			ifexit <- true
			fmt.Println("-close-client-")
			return
		}
		//check the content that client sent,
		fmt.Println(
			"received msg from client:\n\t",
			str_msg_cx,
			"\n\t",
			html.EscapeString(str_msg_cx),
		)
		msg_to_center = newMsg(N)
		//TODO - check the error:
		msg_to_center.parseJSON(str_msg_cx)
		fmt.Println("Node.handleUser", "msgtocenter", msg_to_center.description)
		//and push it to center.
		fmt.Println("Node.handleUser", *msg_to_center)
		N.c_ptr.msg_queue <- *msg_to_center
	}
}

func (N *Node) listenToCenter() {
	//responser
	//fetch the msg from center, and send it to client.
	var json_to_user string
	for {
		select {
		case msg := <-N.msg_from_center:
			msg.source_node = N
			json_to_user = msg.toJSON()
			fmt.Println("Node.handleCenter", json_to_user)
			N.conn.WriteMessage(
				websocket.TextMessage,
				[]byte(json_to_user),
			)
		}
	}
}

//use go statment to call this func
func (N *Node) run(ifexit chan<- bool) {
	//var err error
	var if_listener_exit = make(chan bool)
	go N.listenToUser(if_listener_exit)
	go N.listenToCenter()
	fmt.Println("node::Run()")
	select {
	case <-if_listener_exit:
		//if the listener exit, then the whole node will exit.
		ifexit <- true
		fmt.Println("node::run() -close-node-")
		return
	}
}
