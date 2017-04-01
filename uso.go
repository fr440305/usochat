// USE MANY MANY TINY THEARDS TO SEND OR RECEIVE EVENTS.
// USE SLICE EVERYWHERE - DO NOT USE container/list.
// 不光要思考架构，还要思考架构的迭代与演化。
// 要牢记：软件是长出来的。

//CODE_COMPLETE:
// --all TODOs & FIXMEs
// - documentation: on business logic.

package main

import "fmt"
import "html"
import "net/http"
import "github.com/gorilla/websocket"
import "encoding/json"
import "strconv"

type Msg struct {
	source_node *Node // != nil => from node. ==nil => from center.
	description string
	content     []string
}

func newMsg(source_node *Node) *Msg {
	var res = new(Msg)
	res.source_node = source_node
	res.description = ""
	res.content = []string{}
	fmt.Println("newMsg", res)
	return res
}

func (M *Msg) msgCopy() *Msg {
	return &Msg{
		source_node: M.source_node,
		description: M.description,
		content:     M.content[:],
	}
}

func (M *Msg) setDescription(description string) *Msg {
	M.description = html.EscapeString(description)
	return M
}

func (M *Msg) setContent(content []string) *Msg {
	for i, str := range content {
		content[i] = html.EscapeString(str)
		fmt.Println("Msg.setContent", content[i])
	}
	M.content = content
	fmt.Println("Msg.setContent", content)
	return M
}

//Pay attention to the probobaly-appear errors.
func (M *Msg) parseJSON(json_raw string) error {
	var user_msg struct {
		SouceNode   string   `json:"source_node"`
		Description string   `json:"description"`
		Content     []string `json:"content"`
	}
	json.Unmarshal([]byte(json_raw), &user_msg)
	M.setDescription(user_msg.Description)
	M.setContent(user_msg.Content)
	fmt.Println("Msg.parseJSON", user_msg)
	fmt.Println("Msg.parseJSOn - end.")
	return nil
}

//TODO - check error
//This method transforms the Msg::M to JSON string.
func (M *Msg) toJSON() string {
	var res []byte
	var err error
	fmt.Println("Msg.toJSON", "begin")
	var user_msg = struct {
		SouceNode   string   `json:"source_node"`
		Description string   `json:"description"`
		Content     []string `json:"content"`
	}{M.source_node.iden, M.description, M.content}
	fmt.Println("Msg.toJSON", user_msg)
	res, err = json.Marshal(user_msg)
	if err != nil {
		//TODO - error handler goes here...
	}
	fmt.Println("Msg.toJSON", user_msg)
	fmt.Println("Msg.toJSOn - end.", string(res))
	return string(res)
	//return `{"content":["toJSON","toJSON"]}`
}

func (M Msg) Error() string {
	if M.description == "error" && M.content != nil && len(M.content) != 0 {
		return M.content[0]
	}
	return ""
}

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
			fmt.Println("-close-client-")
			N.c_ptr.removeNode(N)
			ifexit <- true
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

type Center struct {
	msg_queue chan Msg
	dialogs   []*Msg // chatting history
	nodes     []*Node
	upgrader  websocket.Upgrader //Constant
}

func newCenter() *Center {
	fmt.Println("newCenter()")
	return &Center{
		msg_queue: make(chan Msg),
		dialogs:   *new([]*Msg),
		nodes:     *new([]*Node),
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
	}
}

func (C *Center) newNode(w http.ResponseWriter, r *http.Request) *Node {
	var err error
	var res = new(Node)
	fmt.Println("center::newNode()")
	res.msg_from_center = make(chan Msg) //string
	res.c_ptr = C
	res.conn, err = C.upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(
			"fatal - center::newNode -",
			"when creating the websocket.conn",
		)
	}
	C.nodes = append(C.nodes, res)
	fmt.Println("online: ", C.getOnliner())
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
		select {
		case receive_msg = <-C.msg_queue:
			//fmt.Println("Center.handleNodes", "---", msg.source_node)
			//fmt.Println("Center.handleNodes", "---", msg.description)
			//fmt.Println("Center.handleNodes", "---", msg.content)
			//check:
			rec_msg_desp = receive_msg.description
			if rec_msg_desp == "user-login" {
				for _, prev_msg := range C.dialogs {
					fmt.Println("Center.handleNodes", chat_hist)
					chat_hist = append(chat_hist, prev_msg.content[:]...)
				}
				response_msg = receive_msg.msgCopy()
				response_msg.setDescription(string(append([]byte(receive_msg.description), '-', '0')))
				response_msg.setContent(chat_hist) //should be chatting hist.
				boardcast_msg = receive_msg.msgCopy()
				boardcast_msg.setDescription(string(append([]byte(receive_msg.description), '-', '*')))
				boardcast_msg.setContent([]string{strconv.Itoa(C.getOnliner())})
			} else if rec_msg_desp == "user-logout" {
				//no msg-0. only has msg-*.
				boardcast_msg = receive_msg.msgCopy()
				boardcast_msg.setDescription(string(append([]byte(receive_msg.description), '-', '*')))
				boardcast_msg.setContent([]string{strconv.Itoa(C.getOnliner())})
			} else if rec_msg_desp == "user-msg-text" {
				//save the message into Center.dialogs
				C.dialogs = append(C.dialogs, receive_msg.msgCopy())
				response_msg = receive_msg.msgCopy()
				response_msg.source_node = nil
				response_msg.setDescription(string(append([]byte(receive_msg.description), '-', '0')))
				response_msg.setContent([]string{"send successful"}) //should be chatting hist.
				boardcast_msg = receive_msg.msgCopy()
				boardcast_msg.setDescription(string(append([]byte(receive_msg.description), '-', '*')))
			} else {
				//error
			}
			//send them back:
			if response_msg != nil {
				receive_msg.source_node.msg_from_center <- *response_msg
			}
			if boardcast_msg != nil {
				C.boardcast(*boardcast_msg)
			}
			fmt.Println("Center.handleNodes - dialogs:", C.dialogs)
			fmt.Println("Center.handleNodes", "888")
		}
	}
}

//return the number of people online:
func (C *Center) getOnliner() int {
	return len(C.nodes)
}

func main() {
	fmt.Println("http://127.0.0.1:9999")
	var center = newCenter()
	go center.handleNodes()
	//To provide the webpages to the client:
	http.Handle("/", http.FileServer(http.Dir("./frontend")))
	//To handle the websocket request:
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		var if_node_exit = make(chan bool)
		go center.newNode(w, r).run(if_node_exit)
		select {
		case <-if_node_exit:
			fmt.Println("A node exit.")
			return
		}
	})
	http.ListenAndServe(":9999", nil)
}
