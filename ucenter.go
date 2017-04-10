//FileName: ucenter.go
//Description: This file defined a type called Center.

package main

import "net/http"
import "github.com/gorilla/websocket"
import "strconv"

//The duty of the following type, `Center`, is to handle the messages comes from Nodes.
//It receives messages in msg_queue, which are come from Nodes, and then checks them.
//Then, Center will create two catagories of message: response message and boardcast message.
//Finally, Center will send the response message to the Node which send the message to Center.msg_queue,
//and send the boardcast message to all the nodes.

type Center struct {
	msg_queue chan Msg           //Nodes will push their messages here.
	dialogs   []*Msg             //Chatting history
	nodes     []*Node            //All the nodes.
	upgrader  websocket.Upgrader //Constant.
}

//The following constructor creates a new Center instance and returns it.
func newCenter() *Center {
	_ulog("_newCenter")
	return &Center{
		msg_queue: make(chan Msg),
		dialogs:   []*Msg{},
		//dialogs:   *new([]*Msg),
		nodes: []*Node{},
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
	}
}

//If a new user opens the website, a new http.Request and a new http.ResponseWriter will also be created.
//The following method(it is also a constructor. Only center can create node.) extracts the
//ResponseWriter(w), and the Request(r).
//Then it will create a node that `maps` to that user, and push it into C.nodes.
//Finallly, it will return that node.
func (C *Center) newNode(w http.ResponseWriter, r *http.Request) *Node {
	var err error
	var res = new(Node)
	_ulog("\n\n\nCenter.newNode", "A new node goes in!")
	res.msg_from_center = make(chan Msg)
	res.c_ptr = C
	res.conn, err = C.upgrader.Upgrade(w, r, nil)
	if err != nil {
		_ulog("fatal - Center.newNode", "cannot create websocket.conn")
	}
	//TODO - add a node iden allocation.
	C.nodes = append(C.nodes, res)
	return res
}

//The following method removes the useless node from Center.nodes.
//If the node cannot be found, it returns a error.
func (C *Center) removeNode(rm_node *Node) error {
	//TODO - handle the error //
	var i = 0
	var node_ptr *Node = nil
	for i, node_ptr = range C.nodes {
		if node_ptr == rm_node {
			break
		}
	}
	C.nodes = append(C.nodes[:i], C.nodes[i+1:]...)
	return nil
}

//This method send message to all the nodes.
func (C *Center) boardcast(boardcast_msg Msg) error {
	//boardcast_msg.description = string(append([]byte(boardcast_msg.description), '-', '*'))
	for _, N := range C.nodes {
		N.msg_from_center <- boardcast_msg
	}
	return nil //TODO - handle the error//
}

//The follwoing method returns the number of people online.
//For example, if there are three people online, then it will return (3, "3").
func (C *Center) getOnliner() (int, string) {
	return len(C.nodes), string(strconv.Itoa(len(C.nodes)))
}

//This method is the major method of Center type.
//It extracts the message in Center.msg_queue, and check it.
//Then it creates a response message and a boardcast message.
//Then it sends the response message to the Node which the oringinal message comes from.
//Then if sends the boardcast message to all the nodes.
//Notice: Use go statment to call this function.
func (C *Center) handleNodes() {
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
			//_ulog("Center.handleNodes", "---", msg.source_node)
			//_ulog("Center.handleNodes", "---", msg.description)
			_ulog("Center.handleNodes", "receive this Msg from node:", receive_msg.toJSON())
			//check:
			rec_msg_desp = receive_msg.description
			if rec_msg_desp == "login" {
				for _, prev_msg := range C.dialogs {
					_ulog("Center.handleNodes", chat_hist)
					chat_hist = append(chat_hist, prev_msg.content[:]...)
				}
				_, string_onliner = C.getOnliner()
				response_msg = receive_msg.msgCopy('0')
				response_msg.setContent(chat_hist)
				boardcast_msg = receive_msg.msgCopy('*')
				boardcast_msg.setContent([]string{string_onliner})
			} else if rec_msg_desp == "logout" {
				//remove this node:
				C.removeNode(receive_msg.source_node)
				_, string_onliner = C.getOnliner()
				response_msg = receive_msg.msgCopy('0')
				response_msg.setContent([]string{"tara"})
				boardcast_msg = receive_msg.msgCopy('*')
				boardcast_msg.setContent([]string{string_onliner})
			} else if rec_msg_desp == "msg-text" {
				//save the message into Center.dialogs
				_ulog("Center.handleNodes", "saves this msg to chattinghist slice.")
				C.dialogs = append(C.dialogs, receive_msg.msgCopy(' '))
				response_msg = receive_msg.msgCopy('0')
				response_msg.setContent([]string{"send successful"}) //should be chatting hist.
				_ulog("Center.handleNodes", "__DEBUG_RESPMSG__", response_msg.toJSON)
				boardcast_msg = receive_msg.msgCopy('*')
			} else if rec_msg_desp == "msg-pic" {
				//picture.
				_ulog("Center.handleNodes", "received a picture.")
				response_msg = receive_msg.msgCopy('0')
				response_msg.setContent([]string{"send successful"}) //should be chatting hist.
				//_ulog("Center.handleNodes", "__DEBUG_RESPMSG__", response_msg.toJSON)
				boardcast_msg = receive_msg.msgCopy('*')
			} else {
				//TODO - handle the error//
			}
			//send them back:
			// always true:
			if response_msg != nil {
				_ulog("Center.handleNodes", "__DEBUG_RECEIVEMSG__", receive_msg.source_node)
				receive_msg.source_node.msg_from_center <- *response_msg
				_ulog("Center.handleNodes", "has responsed this Msg to the node:", response_msg.toJSON())
			}
			if boardcast_msg != nil {
				C.boardcast(*boardcast_msg)
				_ulog("Center.handleNodes", "has boardcasten this Msg to all the nodes:", boardcast_msg.toJSON())
			}
		}
	}
}
