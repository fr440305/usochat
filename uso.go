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

var Uso_connpool = []Uconn{}

func connpool_add(pool []Uconn, uc Uconn) {
}

func connpool_del(pool []Uconn, uc *Uconn) {
}

var Uso_roompool = []Room{}

func roompool_getNameList() []string {
	list := []string{}
	for _, r := range Uso_roompool {
		list = append(list, r.Name)
	}
	return list
}

func roompool_add(r Room) {
}

func roompool_del(r *Room) {
}

func roompool_getRoomByName(name string) *Room {
	for i, r := range Uso_roompool {
		if r.Name == name {
			return &Uso_roompool[i]
		}
	}
	return nil
}

var Uso_websocket_upgrader = &websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

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

type Uconn struct {
	//if you are not in a Room, you do not need a name
	Name   string
	Quit   bool
	wsconn *websocket.Conn
}

func (c *Uconn) Read() (string, []string) {
	fmt.Println("Read - here")
	_, bar, err := c.wsconn.ReadMessage()
	fmt.Println("Uconn.Read string(bar) == ", string(bar))
	if err != nil {
		c.Quit = true
		fmt.Println("Uconn.Read ty, co == ", UsoMsg_Die, nil)
		return UsoMsg_Die, nil
	}
	msg := NewMessage(bar)
	if len(msg) == 0 {
		fmt.Println("Uconn.Read ty, co == ", "", nil)
		return "", nil
	} else {
		fmt.Println("Uconn.Read ty, co == ", msg[0], msg[1:])
		return msg[0], msg[1:]
	}
}

func (c *Uconn) Write(msg_cont ...string) {
	c.wsconn.WriteMessage(websocket.TextMessage, Message(msg_cont).ToJbar())

}

type Room struct {
	Name  string
	Hist  []string
	Conns []*Uconn
}

// public.
func (r Room) ServeMember(uc *Uconn) {
	// jobs <-- for all ur.Conns.R
	// if ur.jobs is closed (ur.IsActive == false), then open it and run the runner
	// for a := range uc.R { ur.jobs <- a }; close and quit
}

type Hall []*Uconn

func (h Hall) horn(msg_cont ...string) {
	for _, uc := range []*Uconn(h) {
		uc.Write(msg_cont...)
	}
}

func (h Hall) addAndHorn(uc *Uconn) {
	h = []*Uconn(append([]*Uconn(h), uc))
	h.horn(HallHorn_UsoToHall)
}

func (h Hall) delAndHorn(uc *Uconn) {
}

func (h Hall) handleUsoReq_ToRoom(uc *Uconn, co []string) {
	if co == nil {
		uc.Write(HallResp_Error, "Room name is required")
		return
	}
	room := roompool_getRoomByName(co[0])
	if room == nil {
		uc.Write(HallResp_Error, "No such room")
		return
	}
	h.delAndHorn(uc)
	room.ServeMember(uc)
	// check if this room is empty
	if len(room.Conns) == 0 {
		roompool_del(room)
	}
	if uc.Quit {
		return
	} else {
		h.addAndHorn(uc)
	}
}

func (h Hall) handleUsoReq_AddRoom(uc *Uconn, co []string) {
	if co == nil {
		uc.Write(HallResp_Error, "Room name is required")
		return
	}
	if roompool_getRoomByName(co[0]) != nil {
		uc.Write(HallResp_Error, "This room exists")
		return
	}
	room := Room{
		Name:  co[0],
		Hist:  []string{},
		Conns: []*Uconn{},
	}
	roompool_add(room)
	h.horn(HallHorn_UsoAddRoom, co[0])
	room.ServeMember(uc)
	if len(room.Conns) == 0 {
		roompool_del(&room)
	}
	if uc.Quit {
		return
	} else {
		h.addAndHorn(uc)
	}
}

func (h Hall) ServeGuest(uc *Uconn) {
	h.addAndHorn(uc)
	uc.Write(append([]string{HallResp_Rooms}, roompool_getNameList()...)...)
	for {
		ty, co := uc.Read()
		if ty == UsoMsg_Die {
			fmt.Println("Hall.ServeGuest die")
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

var Uso_hall = Hall{}

// export
func ServeWs(w http.ResponseWriter, r *http.Request) {
	conn, err := Uso_websocket_upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("ServeWs upgrade error")
		return
	}
	uso_conn := Uconn{
		Name:   "",
		Quit:   false,
		wsconn: conn,
	}
	Uso_connpool = append(Uso_connpool, uso_conn)
	fmt.Println("ServeWs Uso_connpool == ", Uso_connpool)
	Uso_hall.ServeGuest(&uso_conn)
	connpool_del(Uso_connpool, &uso_conn)
}
