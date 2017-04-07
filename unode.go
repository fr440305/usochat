// USE MANY MANY TINY THEARDS TO SEND OR RECEIVE EVENTS.
// USE SLICE EVERYWHERE - DO NOT USE container/list.
// 不光要思考架构，还要思考架构的迭代与演化。
// 要牢记：软件是长出来的。

package main

import "github.com/gorilla/websocket"

//type Node maps to a client.
type Node struct {
	msg_from_center chan Msg
	c_ptr           *Center         // a pointer to center.
	conn            *websocket.Conn //connent client to node
	nid             int64           // the identification for node.
}

//The following function returns the string form of the node id.
func (N *Node) idString() string {
	//TODO//
	return ""
}

//The following function will be called in a go statment because it is a theard.
//It extracts the JSON string message form the user and
//handle this message. It will send the message to center if nessesary.
func (N *Node) handleUser(ifexit chan<- bool) {
	var err error
	var msg_cx []byte      // the byte array from user.
	var str_msg_cx string  // the conversion for byte array. Will be a JSON string.
	var msg_to_center *Msg // the message that needs to send to center.
	for {
		//the code will be blocked here(conn.ReadMessage():
		//but don't worry, becase it's in the go statment.
		_, msg_cx, err = N.conn.ReadMessage()
		str_msg_cx = string(msg_cx[:])
		if err != nil {
			//the client was closed.
			_ulog("\n\n\nNode.handleUser", "A user has been leaving!")
			msg_to_center = newMsg(N)
			msg_to_center.setDescription("logout")
			N.c_ptr.msg_queue <- *msg_to_center
			ifexit <- true
			_ulog("Node.handleUser", "exits")
			return
		}
		//check the content that client sent,
		_ulog("\nNode.handleUser", "received JSON:", str_msg_cx)
		msg_to_center = newMsg(N)
		//TODO - check the error//
		msg_to_center.parseJSON(str_msg_cx)
		//and push it to center.
		_ulog("Node.handleUser", "send this msg to center:", msg_to_center.toJSON())
		N.c_ptr.msg_queue <- *msg_to_center
	}
}

//responser
//fetch the msg from center, and send it to client.
func (N *Node) handleCenter() {
	var msg Msg             // The message received from center.
	var json_to_user string // The JSON string that needs to be sent to user.
	for {
		select {
		case msg = <-N.msg_from_center:
			_ulog("Node.handleCenter", "receives this Msg from center:", msg.toJSON())
			if msg.description == "logout-0" {
				_ulog("Node.handleCenter", "exits.")
				return
			} else {
				msg.source_node = N
				json_to_user = msg.toJSON()
				_ulog("Node.handleCenter", "Send this json to user:", json_to_user)
				N.conn.WriteMessage(websocket.TextMessage, []byte(json_to_user))
			}
		}
	}
}

//use go statment to call this func
func (N *Node) run(ifexit chan<- bool) {
	//var err error
	var if_listener_exit = make(chan bool)
	go N.handleUser(if_listener_exit)
	go N.handleCenter()
	_ulog("Node.run")
	select {
	case <-if_listener_exit:
		ifexit <- true
		_ulog("Node.run", "exits")
		return
	}
}
