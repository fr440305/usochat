package uso // import "github.com/fr440305/uso"

import "strconv"
import "github.com/gorilla/websocket"
import "net/http"
import "encoding/json"
import "fmt"

// MsgType:
const (
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
	_, bar, err := c.wsconn.ReadMessage()
	if err != nil {
		c.Quit = true
		return UsoMsg_Die, nil
	}
	msg := NewMessage(bar)
	if len(msg) == 0 {
		return "", nil
	} else {
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
			h.horn(HallHorn_UsoExitHall)
			return
		}
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
		case UsoMsg_AddRoom:
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
}

//========
//--vars--
//========

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
		return nil
	}
	conn := Conn{
		Name:   "",
		Quit:   false,
		wsconn: ws_conn,
	}
	Uso_connpool = append(Uso_connpool, &conn)
	return &conn
}

func connpool_del(c *Conn) {
	for i := 0; i < len(Uso_connpool); i++ {
		if Uso_connpool[i] == c {
			Uso_connpool = append(Uso_connpool[:i], Uso_connpool[i+1:]...)
		}
	}
}

func connpool_lentostr() string {
	return strconv.Itoa(len(Uso_connpool))
}

//============
//--roompool--
//============

var Uso_roompool = []*Room{} // -> var uso_RoomPool = RoomPool{}

func roompool_add(name string) *Room {
	return nil
}

func roompool_del(r *Room) {
}

func roompool_getRoomByName(name string) *Room {
	return nil
}

func roompool_namelist() []string {
	return nil
}

// export
func ServeWs(w http.ResponseWriter, r *http.Request) {
	conn := connpool_add(w, r)
	if conn == nil {
		return
	}
	Uso_hall.ServeGuest(conn)
	connpool_del(conn)
}
