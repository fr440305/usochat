// USE MANY MANY TINY THEARDS TO SEND OR RECEIVE EVENTS.
// USE SLICE EVERYWHERE - DO NOT USE container/list.
// 不光要思考架构，还要思考架构的迭代与演化。
// 要牢记：软件是长出来的。

//CODE_COMPLETE:
// --all TODOs & FIXMEs
// - documentation: on business logic.

package main

import "net/http"
import "github.com/gorilla/websocket"
import "strconv"

type Center struct {
	msg_queue chan Msg
	dialogs   []*Msg // chatting history
	nodes     []*Node
	upgrader  websocket.Upgrader //Constant
}

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

//This method removes the useless node from center.nodes.
//If the node cannot be found, it returns a error.
func (C *Center) removeNode(rm_node *Node) error {
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
	return nil //TODO//
}

//listen and handle the msg.
//use go statment to call this func.
func (C *Center) handleNodes() {
	var receive_msg Msg
	var response_msg *Msg
	var boardcast_msg *Msg
	var rec_msg_desp string
	var chat_hist []string
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
			if rec_msg_desp == "user-login" {
				for _, prev_msg := range C.dialogs {
					_ulog("Center.handleNodes", chat_hist)
					chat_hist = append(chat_hist, prev_msg.content[:]...)
				}
				response_msg = receive_msg.msgCopy('0')
				response_msg.setContent(chat_hist)
				boardcast_msg = receive_msg.msgCopy('*')
				boardcast_msg.setContent([]string{strconv.Itoa(C.getOnliner())})
			} else if rec_msg_desp == "user-logout" {
				//TODO - remove this node.
				response_msg = receive_msg.msgCopy('0')
				response_msg.setContent([]string{"tara"})
				boardcast_msg = receive_msg.msgCopy('*')
				boardcast_msg.setContent([]string{strconv.Itoa(C.getOnliner())})
			} else if rec_msg_desp == "user-msg-text" {
				//save the message into Center.dialogs
				_ulog("Center.handleNodes", "saves this msg to chattinghist slice.")
				C.dialogs = append(C.dialogs, receive_msg.msgCopy(' '))
				response_msg = receive_msg.msgCopy('0')
				response_msg.setContent([]string{"send successful"}) //should be chatting hist.
				_ulog("Center.handleNodes", "__DEBUG_RESPMSG__", response_msg.toJSON)
				boardcast_msg = receive_msg.msgCopy('*')
			} else {
				//error
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

//return the number of people online:
func (C *Center) getOnliner() int {
	return len(C.nodes)
}
