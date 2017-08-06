package uso // import "github.com/fr440305/uso"

import "github.com/gorilla/websocket"
import "net/http"
import "encoding/json"
import "fmt"

const (
	// MsgType:

	// um : 由 uso-conn 发出的消息
	UsoMsg_ToRoom   = "um:inr"
	UsoMsg_AddRoom  = "um:++r"
	UsoMsg_ExitRoom = "um:exr"
	UsoMsg_Say      = "um:say"
	UsoMsg_Die      = "um:die"

	HallHorn_UsoToHall   = "hh:uinh"
	HallHorn_UsoExitHall = "hh:uexh"
	HallHorn_UsoAddRoom  = "hh:u++r"
	HallHorn_DelRoom     = "hh:r--"
	HallHorn_Error       = "hh:err"

	HallResp_Rooms = "hr:rooms"
	HallResp_Error = "hr:err"
)

type Message []string

// string -> []interface{} // tokenlize
// []interface{} -> UsoMsg // done
// use json.Unmarshal
func NewMessage(bar []byte) Message {
	strarr := &[]string{}
	err := json.Unmarshal(bar, strarr)
	if err != nil {
		fmt.Println(err.Error())
		return Message{""}
	}
	fmt.Println("NewMessage", string(bar))
	return Message(*strarr)
}

// return []byte(`["Type", "Cont-1", "Cont-2", ... , "Cont-n"]`)
// use json.Marshal
func (um Message) ToJbar() []byte {
	bar, err := json.Marshal(um)
	if err != nil {
		fmt.Println("Message.ToJbar error", err.Error())
		return []byte("")
	} else {
		return bar
	}
}

type Conn struct {
	//if you are not in a Room, you do not need a name
	Name   string
	Quit   bool
	wsconn *websocket.Conn
}

func (c *Conn) Read() (string, []string) {
	fmt.Println("Read - here")
	_, bar, err := c.wsconn.ReadMessage()
	fmt.Println("Conn.Read string(bar) == ", string(bar))
	if err != nil {
		c.Quit = true
		fmt.Println("Conn.Read ty, co == ", UsoMsg_Die, nil)
		return UsoMsg_Die, nil
	}
	msg := NewMessage(bar)
	if len(msg) == 0 {
		fmt.Println("Conn.Read ty, co == ", "", nil)
		return "", nil
	} else {
		fmt.Println("Conn.Read ty, co == ", msg[0], msg[1:])
		return msg[0], msg[1:]
	}
}

func (c *Conn) Write(msg_type string, msg_cont ...string) {
	msg := append([]string{msg_type}, msg_cont...)
	c.wsconn.WriteMessage(websocket.TextMessage, Message(msg).ToJbar())

}

type Hall []*Conn

func (h Hall) horn(msg_type string, msg_cont ...string) {
	for _, uc := range []*Conn(h) {
		uc.Write(msg_type, msg_cont...)
	}
}

func (h *Hall) addAndHorn(uc *Conn) {
	*h = append(*h, uc)
	h.horn(HallHorn_UsoToHall)
}

func (h *Hall) delAndHorn(uc *Conn) {
	for i := 0; i < len(*h); i++ {
		if (*h)[i] == uc {
			*h = append((*h)[:i], (*h)[i+1:]...)
			return
		}
	}
}

func (h *Hall) handleUsoReq_ToRoom(uc *Conn, co []string) {
	if co == nil {
		uc.Write(HallResp_Error, "Room name is required")
		return
	}
}

func (h *Hall) handleUsoReq_AddRoom(uc *Conn, co []string) {
	if co == nil {
		uc.Write(HallResp_Error, "Room name is required")
		return
	}
}

func (h *Hall) ServeGuest(uc *Conn) {
	h.addAndHorn(uc)
	fmt.Println(h)
	uc.Write(HallResp_Rooms, "waiting..")
	for {
		ty, co := uc.Read()
		if ty == UsoMsg_Die {
			fmt.Println("Hall.ServeGuest die")
			h.delAndHorn(uc)
			return
		}
		fmt.Println("Hall.ServeGuest ty, co == ", ty, co)
		switch ty {
		case UsoMsg_ToRoom:
			//h.handleUsoReq_ToRoom(uc, co)
		case UsoMsg_AddRoom:
			//h.handleUsoReq_AddRoom(uc, co)
		default:
			uc.Write(HallResp_Error, "Invalid message")
		}
	}
}

type Room struct {
	Name  string
	Hist  []string
	Conns []*Conn
}

func (r Room) ServeMember(uc *Conn) {
	// jobs <-- for all ur.Conns.R
	// if ur.jobs is closed (ur.IsActive == false), then open it and run the runner
	// for a := range uc.R { ur.jobs <- a }; close and quit
}

var Uso_hall = &Hall{} // -> var uso_Hall = ...

var uso_WebsocketUpgrader = &websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

//============
//--connpool--
//============

var Uso_connpool = []*Conn{} // -> var uso_ConnPool = ConnPool{}

func connpool_add(w http.ResponseWriter, r *http.Request) *Conn {
	ws_conn, err := uso_WebsocketUpgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("ServeWs upgrade error")
		return
	}
	conn := Conn{
		Name:   "",
		Quit:   false,
		wsconn: ws_conn,
	}
	Uso_connpool = append(Uso_connpool, &conn)
	return conn
}

func connpool_del(c *Conn) {
}

//============
//--roompool--
//============

var Uso_roompool = []*Room{} // -> var uso_RoomPool = RoomPool{}

func roompool_add(name string) *Room {
}

func roompool_del(r *Room) {
}

// export
func ServeWs(w http.ResponseWriter, r *http.Request) {
	conn := connpool_add(w, r)
	Uso_hall.ServeGuest(&conn)
	connpool_del(conn)
}
