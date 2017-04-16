//FileName:
//	unode.go
//FileDescription:
//	This go source file defined three go type: Usor, Room, and Center.
//Author(s):
//	__HUO_YU__

package main

import "github.com/gorilla/websocket"
import "net/http"
import "strconv"

//type Usor maps to a client.
type Usor struct {
	msgchan chan Msg
	room    *Room
	conn    *websocket.Conn //connent client to node
	nid     int64           // the identification for node.
}

//The following function returns the string form of the node id.
func (U *Usor) idString() string {
	//TODO//
	return ""
}

//The following function will be called in a go statment because it is a theard.
//It extracts the JSON string message form the user and
//handle this message. It will send the message to center if nessesary.
func (U *Usor) handleUser(ifexit chan<- bool) {
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
			_ulog("\n\n\nUsor.handleUser", "A user has been leaving!")
			msg_to_center = newMsg(N)
			msg_to_center.setDescription("logout")
			N.c_ptr.msg_queue <- *msg_to_center
			ifexit <- true
			_ulog("Usor.handleUser", "exits")
			return
		}
		//check the content that client sent,
		_ulog("\nUsor.handleUser", "received JSON:", str_msg_cx)
		msg_to_center = newMsg(N)
		//TODO - check the error//
		msg_to_center.parseJSON(str_msg_cx)
		//and push it to center.
		_ulog("Usor.handleUser", "send this msg to center:", msg_to_center.toJSON())
		N.c_ptr.msg_queue <- *msg_to_center
	}
}

//responser
//fetch the msg from center, and send it to client.
func (U *Usor) handleRoom() {
	var msg Msg             // The message received from center.
	var json_to_user string // The JSON string that needs to be sent to user.
	for {
		select {
		case msg = <-N.msg_from_center:
			_ulog("Usor.handleCenter", "receives this Msg from center:", msg.toJSON())
			if msg.description == "logout-0" {
				_ulog("Usor.handleCenter", "exits.")
				return
			} else {
				msg.source_node = N
				json_to_user = msg.toJSON()
				_ulog("Usor.handleCenter", "Send this json to user:", json_to_user)
				N.conn.WriteMessage(websocket.TextMessage, []byte(json_to_user))
			}
		}
	}
}

//use go statment to call this func
func (U *Usor) run(ifexit chan<- bool) {
	//var err error
	var if_listener_exit = make(chan bool)
	go N.handleUser(if_listener_exit)
	go N.handleCenter()
	_ulog("Usor.run")
	select {
	case <-if_listener_exit:
		ifexit <- true
		_ulog("Usor.run", "exits")
		return
	}
}

//FileName: uroom.go
//Description: This file defined a type called Center.

type Room struct {
	rid       uint64
	msg_queue chan Msg           //Usors will push their messages here.
	msg_hist  []Msg              //history messages.
	nodes     []Usor             //All the nodes.
	upgrader  websocket.Upgrader //Constant.
}

//The following method removes the useless node from Center.nodes.
//If the node cannot be found, it returns a error.
func (R *Room) removeUsor(rm_node *Usor) error {
	//TODO - handle the error //
	var i = 0
	var node_ptr *Usor = nil
	for i, node_ptr = range C.nodes {
		if node_ptr == rm_node {
			break
		}
	}
	C.nodes = append(C.nodes[:i], C.nodes[i+1:]...)
	return nil
}

//This method send message to all the nodes.
func (R *Room) boardcast(boardcast_msg Msg) error {
	//boardcast_msg.description = string(append([]byte(boardcast_msg.description), '-', '*'))
	for _, N := range C.nodes {
		N.msg_from_center <- boardcast_msg
	}
	return nil //TODO - handle the error//
}

//The follwoing method returns the number of people in this room.
//For example, if there are three people online, then it will return (3, "3").
func (R *Room) getNumUser() (int, string) {
	return len(C.nodes), string(strconv.Itoa(len(C.nodes)))
}

//This method is the major method of Center type.
//It extracts the message in Center.msg_queue, and check it.
//Then it creates a response message and a boardcast message.
//Then it sends the response message to the Usor which the oringinal message comes from.
//Then if sends the boardcast message to all the nodes.
//Notice: Use go statment to call this function.
func (R *Room) handleUsors() {
	var receive_msg Msg
	var response_msg *Msg
	var boardcast_msg *Msg
	var rec_msg_desp string
	var chat_hist []string
	var string_onliner string
	for {
		//initialize
		response_msg = nil
		boardcast_msg = nil
		rec_msg_desp = ""
		chat_hist = []string{}
		select {
		case receive_msg = <-C.msg_queue:
			//_ulog("Center.handleUsors", "---", msg.source_node)
			//_ulog("Center.handleUsors", "---", msg.description)
			_ulog("Center.handleUsors", "receive this Msg from node:", receive_msg.toJSON())
			//check:
			rec_msg_desp = receive_msg.description
			if rec_msg_desp == "login" {
				for _, prev_msg := range C.dialogs {
					_ulog("Center.handleUsors", chat_hist)
					chat_hist = append(chat_hist, prev_msg.content[:]...)
				}
				_, string_onliner = C.getOnliner()
				response_msg = receive_msg.msgCopy('0')
				response_msg.setContent(chat_hist)
				boardcast_msg = receive_msg.msgCopy('*')
				boardcast_msg.setContent([]string{string_onliner})
			} else if rec_msg_desp == "logout" {
				//remove this node:
				C.removeUsor(receive_msg.source_node)
				_, string_onliner = C.getOnliner()
				response_msg = receive_msg.msgCopy('0')
				response_msg.setContent([]string{"tara"})
				boardcast_msg = receive_msg.msgCopy('*')
				boardcast_msg.setContent([]string{string_onliner})
			} else if rec_msg_desp == "msg-text" {
				//save the message into Center.dialogs
				_ulog("Center.handleUsors", "saves this msg to chattinghist slice.")
				C.dialogs = append(C.dialogs, receive_msg.msgCopy(' '))
				response_msg = receive_msg.msgCopy('0')
				response_msg.setContent([]string{"send successful"}) //should be chatting hist.
				_ulog("Center.handleUsors", "__DEBUG_RESPMSG__", response_msg.toJSON)
				boardcast_msg = receive_msg.msgCopy('*')
			} else if rec_msg_desp == "msg-pic" {
				//picture.
				_ulog("Center.handleUsors", "received a picture.")
				response_msg = receive_msg.msgCopy('0')
				response_msg.setContent([]string{"send successful"}) //should be chatting hist.
				//_ulog("Center.handleUsors", "__DEBUG_RESPMSG__", response_msg.toJSON)
				boardcast_msg = receive_msg.msgCopy('*')
			} else {
				//TODO - handle the error//
			}
			//send them back:
			// always true:
			if response_msg != nil {
				_ulog("Center.handleUsors", "__DEBUG_RECEIVEMSG__", receive_msg.source_node)
				receive_msg.source_node.msg_from_center <- *response_msg
				_ulog("Center.handleUsors", "has responsed this Msg to the node:", response_msg.toJSON())
			}
			if boardcast_msg != nil {
				C.boardcast(*boardcast_msg)
				_ulog("Center.handleUsors", "has boardcasten this Msg to all the nodes:", boardcast_msg.toJSON())
			}
		}
	}
}

type Center struct {
}

func (C *Center) newRoom() *Room {
	return nil
}

func (C *Center) newNode() *Node {
	return nil
}
