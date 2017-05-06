//unode.go
//This go source file defined three go type: Usor, Room, and Center.
//Author(s): __HUO_YU__

package main

import "github.com/gorilla/websocket"
import "net/http"

//import "strconv"

//type Usor maps to a client.
type Usor struct {
	name string
	qmsg *MsgList
	eden *Eden
	room *Room
	conn *websocket.Conn //client <--conn--> usor
}

func (U *Usor) join(room_name string) error {
	if U == nil {
		_ulog("@err@ Usor.join U == nil")
		return newErrMsg("U == nil")
	}
	if U.room != nil {
		//cannot join
		return newErrMsg("A room-usor can not join in another room.")
	}
	_ulog("@std@ Usor.join Requesting the room: ", room_name)
	var res_room = U.eden.ReqRoom(room_name)
	U.room = res_room
	res_room.AddUsor(U)
	_ulog("@std@ Usor.join Room={", &U.room, U.room.name, U.room.usors, "}")
	_ulog("@std@ Usor.join Join Successful.")
	return nil //good`
}

func (U *Usor) handleClient() {
	if U == nil {
		_ulog("@err@ Usor.handleClient U == nil")
		return
	}
	_ulog("@std@ Usor.handleClient")
	//var msg *Msg
	for {
		msgtype, barjson, err := U.conn.ReadMessage()
		if err != nil {
			//Gone.
			_ulog("@err@ Usor.handleClient", err.Error())
			U.conn.Close()
			return
		} else {
			if msgtype == websocket.TextMessage {
				var client_msg = newBarMsg(barjson)
				_ulog("@std@ Usor.handleClient type=", msgtype, string(barjson))
				if client_msg.Summary == "join" {
					_ulog("@std Usor.handleClient join")
					U.join(client_msg.Content[0][0])
				}
			} else {
				_ulog("@std@ Usor.handleClient type=", msgtype, barjson)
			}
		}
	}
}

func (U *Usor) Run() {
	U.handleClient()
}

func (U *Usor) OnEden(room_name_list []string) {
	_ulog("@std@ Usor.OnEden The name-list is:", room_name_list)
	U.conn.WriteMessage(
		websocket.TextMessage,
		newMsg("room-name-list", [][]string{room_name_list}).barjsonify(),
	)
}

func (U *Usor) OnRoom(chist [][]string) {
	_ulog("@std@ Usor.OnRoom The chist=chat-history is:", chist)
	U.conn.WriteMessage(
		websocket.TextMessage,
		newMsg("chist", [][]string(chist)).barjsonify(),
	)
}

func (U *Usor) OnBoardcasted(msg *Msg) {
	U.conn.WriteMessage(
		websocket.TextMessage,
		msg.barjsonify(),
	)
}

type UsorList []*Usor

func (UL *UsorList) add(usor *Usor) *Usor {
	if UL == nil || usor == nil {
		_ulog("@err@ UsorList.add U||usor == nil")
	}
	*UL = append(*UL, usor)
	return usor
}

func (UL *UsorList) rm(usor *Usor) *Usor {
	return nil
}

func (UL UsorList) boardcast(msg *Msg) *Msg {
	if msg == nil {
		_ulog("@err@ UsorList.boardcast msg == nil")
		return nil
	}
	for _, u := range UL {
		u.OnBoardcasted(msg)
	}
	return msg
}

func (UL UsorList) usorAmount() int64 {
	return int64(len(UL))
}

type Room struct {
	name   string
	qmsg   *MsgList
	chist  [][]string //chat history
	usors  *UsorList
	center *Center
}

func (R *Room) handleCenter() {
}

func (R *Room) handleUsors() {
}

func (R Room) OnKilled() {
}

func (R Room) OnSaid(usor_name string, dialog string) {
}

func (R *Room) AddUsor(usor *Usor) *Usor {
	if R == nil || usor == nil {
		_ulog("@err@ Room.AddUsor R||usor == nil")
		return nil
	}
	R.usors.add(usor)
	usor.OnRoom(R.chist)
	//boardcast
	R.usors.boardcast(newMsg("join", [][]string{[]string{usor.name}}))
	return usor
}

func (R *Room) Run() {
}

type RoomList []*Room

func (RL *RoomList) add(room *Room) *Room {
	if RL == nil || room == nil {
		_ulog("@err@ RoomList.add RL||room == nil")
		return nil
	}
	*RL = append(*RL, room)
	return room
}

func (RL *RoomList) rm(room *Room) *Room {
	return nil
}

func (RL RoomList) boardcast(msg *Msg) *Msg {
	return nil
}

func (RL RoomList) lookup(room_name string) *Room {
	for _, r := range RL {
		if r != nil && r.name == room_name {
			return r
		}
	}
	//no room
	return nil
}

func (RL RoomList) usorAmount() int64 {
	return 0
}

func (RL RoomList) roomAmount() int64 {
	return int64(len(RL))
}

func (RL RoomList) list() []string {
	var res = []string{}
	for _, r := range RL {
		if r != nil {
			res = append(res, r.name)
		}
	}
	return res
}

type Eden struct {
	center *Center
	guests *UsorList
}

func (E Eden) ReqRoom(room_name string) *Room {
	var res_room = E.center.rooms.lookup(room_name)
	if res_room != nil {
		return res_room
	} else {
		return E.center.NewRoom(room_name)
	}
}

func (E *Eden) AddUsor(usor *Usor) *Usor {
	if E == nil || usor == nil {
		_ulog("@err@ Eden.AddUsor E||usor == nil pointer.")
		return nil
	}
	E.guests.add(usor)
	usor.OnEden(E.center.RoomNameList())
	return usor
}

//The main server
type Center struct {
	eden        *Eden
	rooms       *RoomList
	ws_upgrader websocket.Upgrader //const
}

func newCenter(pid int) *Center {
	var center = new(Center)
	center.eden = center.newEden()
	center.rooms = new(RoomList)
	_ulog("@pid@", pid)
	return center
}

func (C *Center) newEden() *Eden {
	var eden = new(Eden)
	eden.center = C
	eden.guests = new(UsorList)
	//C.eden = eden
	_ulog("@std@", "Center.newEden")
	return eden
}

//will be called by Eden.ReqRoom
func (C *Center) NewRoom(name string) *Room {
	var room = new(Room)
	_ulog("@std@ Center.newRoom name=", name)
	room.name = name
	room.qmsg = new(MsgList)
	room.chist = [][]string{}
	room.usors = new(UsorList)
	room.center = C
	C.rooms.add(room)
	_ulog("@std@ Center.newRoom newRoom={", room.name, room.chist, room.usors, "}")
	_ulog("@std@", "Center.newRoom rooms=", C.rooms)
	return room
}

//will be called iff a new client opened.
func (C *Center) newUsor(w http.ResponseWriter, r *http.Request) *Usor {
	var usor = new(Usor)
	var err error
	usor.qmsg = new(MsgList)
	usor.eden = C.eden
	usor.room = nil
	usor.conn, err = C.ws_upgrader.Upgrade(w, r, w.Header())
	if err != nil {
		_ulog("@err@", "Center.newUsor", err.Error())
		return nil
	}
	//C.eden.add this usor.
	go usor.Run()
	return usor
}

func (C *Center) RoomNameList() []string {
	return C.rooms.list()
}

func (C *Center) handleRooms() error {
	//including eden.
	return nil
}

func (C *Center) Run() {
	http.Handle("/", http.FileServer(http.Dir("frontend")))
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		_ulog("@std@", "Center.run()", "/ws")
		C.eden.AddUsor(C.newUsor(w, r))
	})
	http.ListenAndServe(":9999", nil)
}
