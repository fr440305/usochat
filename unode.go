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
		return newErrMsg("U == nil")
	}
	if U.room != nil {
		//cannot join
		return newErrMsg("A room-usor can not join in another room.")
	}
	if room_name == "" {
		return newErrMsg("room name cannot be empty.")
	}
	_ulog("@std@ Usor.join Requesting the room: ", room_name)
	U.room = U.eden.ReqRoom(room_name)
	U.room.AddUsor(U)
	U.eden.RmUsor(U)
	_ulog("@std@ Usor.join {room, room_name, usors}={", U.room, U.room.name, U.room.usors, "}")
	_ulog("@std@ Usor.join Join Successful.")
	return nil
}

func (U *Usor) exitroom(if_rm_room string) error {
	var err error
	if U == nil {
		return newErrMsg("Usor.exitroom - U == nil.")
	}
	if if_rm_room != "rm" && if_rm_room != "rsv" {
		return newErrMsg("Usor.exitroom - parametric != {rm, rsv}")
	}
	if U.room == nil {
		return newErrMsg("A eden-usor cannot exit room because it was already exited.")
	}
	U.room.RmUsor(U)
	if if_rm_room == "rm" {
		err = U.room.SelfKill()
	}
	U.room = nil
	U.eden.AddUsor(U)
	_ulog("@std@ Usor.exitroom exit successful.")
	return err
}

func (U *Usor) setname(name string) error {
	var err error
	if name == "" {
		name = "YouKe:Passerby"
	}
	if U == nil {
		return newErrMsg("U == nil")
	}
	if U.room != nil { //room-usor
		err = U.room.RenameUsor(U, name)
	} else { //eden-usor
		U.name = name
	}
	if err != nil {
		return err
	} else {
		U.writeMsg(newMsg("~~name", [][]string{[]string{name}}))
		return nil
	}
}

func (U *Usor) say(dialog string) error {
	if U == nil {
		return newErrMsg("U == nil")
	}
	if U.room == nil {
		return newErrMsg("A eden-usor cannot say.")
	}
	if dialog == "" {
		return newErrMsg("dialog cannot be empty")
	}
	U.room.AddDialog(U.name, dialog)
	return nil
}

func (U *Usor) readMsg() *Msg {
	var msg *Msg
	_, barjson, err := U.conn.ReadMessage()
	if err != nil {
		//going-away
		U.conn.Close()
		return newMsg("gone", [][]string{[]string{}})
	} else {
		//here may appears error
		msg = newBarMsg(barjson)
		if msg == nil {
			return newErrMsg("unexpected input json from client.")
		}
	}
	return msg

}

func (U *Usor) writeMsg(msg *Msg) error {
	var err = U.conn.WriteMessage(websocket.TextMessage, msg.barjsonify())
	return err
}

func (U *Usor) handleClient() {
	var client_msg *Msg
	var err error
	if U == nil {
		_ulog("@err@ Usor.handleClient U == nil")
		return
	}
	_ulog("@std@ Usor.handleClient")
	//var msg *Msg
	for {
		client_msg = U.readMsg()
		_ulog("\n@std@ Usor.handleClient", client_msg.Summary)
		switch client_msg.Summary {
		case "gone":
			err = U.exitroom("rsv")
			return //the only return of client-handler.
		case "setname":
			err = U.setname(client_msg.Content[0][0])
		case "join":
			err = U.join(client_msg.Content[0][0])
		case "exitroom":
			err = U.exitroom(client_msg.Content[0][0])
		case "say":
			err = U.say(client_msg.Content[0][0])
		case "error":
			err = client_msg
		}
		if err != nil {
			_ulog("@err@ Usor.handleClient", err.Error())
			U.writeMsg(newErrMsg(err.Error()))
		}
	}
}

func (U *Usor) OnRun() {
	U.handleClient()
}

func (U *Usor) OnEden(room_name_list []string) {
	_ulog("@std@ Usor.OnEden The name-list is:", room_name_list)
	U.writeMsg(newMsg("-)eden", [][]string{room_name_list}))
}

func (U *Usor) OnRoom(room_name string, chist [][]string) {
	_ulog("@std@ Usor.OnRoom The chist=chat-history is:", chist)
	var msg_cnt = [][]string{[]string{room_name}}
	msg_cnt = append(msg_cnt, chist...)
	U.writeMsg(newMsg("-)room", msg_cnt))
}

func (U *Usor) OnBoardcasted(msg *Msg) {
	U.writeMsg(msg)
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
	for i, u := range *UL {
		if usor == u {
			*UL = append((*UL)[:i], (*UL)[i+1:]...) //rm
			return usor
		}
	}
	// no this usor
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

func (UL UsorList) list() []string {
	var res = []string{}
	for _, usor := range UL {
		res = append(res, usor.name)
	}
	return res
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

func (R *Room) SelfKill() error {
	if R == nil {
		return newErrMsg("Room.SelfKill R == nil")
	}
	if len(*(R.usors)) == 0 {
		R.center.RmRoom(R)
		return nil
	} else {
		return newErrMsg("Room.SelfKill It still has usor(s) in this room.")
	}
}

func (R *Room) AddDialog(usor_name string, dialog string) {
	if R == nil {
		_ulog("@err@ Room.AddDialog R == nil")
	}
	R.chist = append(R.chist, []string{usor_name, dialog})
	R.usors.boardcast(newMsg("++dialog", [][]string{[]string{usor_name, dialog}}))
}

func (R *Room) AddUsor(usor *Usor) *Usor {
	if R == nil || usor == nil {
		_ulog("@err@ Room.AddUsor R||usor == nil")
		return nil
	}
	R.usors.add(usor)
	usor.OnRoom(R.name, R.chist)
	//boardcast
	R.usors.boardcast(newMsg("++usor", [][]string{[]string{usor.name}, R.usors.list()}))
	return usor
}

func (R *Room) RenameUsor(usor *Usor, new_name string) error {
	var ori_name = usor.name
	if R == nil || usor == nil {
		return newErrMsg("R||usor == nil")
	}
	if new_name == "" {
		new_name = "YouKe:Passerby"
	}
	usor.name = new_name
	R.usors.boardcast(newMsg(
		"~~usor",
		[][]string{[]string{ori_name, new_name}, R.usors.list()},
	))
	return nil
}

func (R *Room) RmUsor(usor *Usor) *Usor {
	if R == nil || usor == nil {
		_ulog("@err@ Room.AddUsor R||usor == nil")
		return nil
	}
	R.usors.rm(usor)
	R.usors.boardcast(newMsg("--usor", [][]string{[]string{usor.name}, R.usors.list()}))
	return usor
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
	if RL == nil || room == nil {
		_ulog("@err@ RoomList.rm RL||room == nil")
		return nil
	}
	for i, r := range *RL {
		if r == room {
			*RL = append((*RL)[:i], (*RL)[i+1:]...)
			return r
		}
	}
	//no this room
	return nil
}

func (RL RoomList) boardcast(msg *Msg) *Msg {
	if msg == nil {
		_ulog("@err@ RoomList.boardcast msg == nil")
		return nil
	}
	for _, r := range RL {
		r.usors.boardcast(msg)
	}
	return msg
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
	var res_room = E.center.LookupRoom(room_name)
	if res_room != nil {
		return res_room
	} else {
		E.guests.boardcast(newMsg("++room", [][]string{[]string{room_name}}))
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

func (E *Eden) RmUsor(usor *Usor) *Usor {
	if E == nil || usor == nil {
		_ulog("@err@ Eden.RmUsor E||usor == nil pointer.")
		return nil
	}
	E.guests.rm(usor)
	return usor
}

func (E Eden) OnRmRoom(room_list []string) {
	E.guests.boardcast(newMsg("--room", [][]string{room_list}))
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

func (C *Center) LookupRoom(name string) *Room {
	if C == nil {
		_ulog("@err@ Center.LookupRoom C == nil")
		return nil
	}
	return C.rooms.lookup(name)
}

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

func (C *Center) RmRoom(room *Room) {
	C.rooms.rm(room)
	C.eden.OnRmRoom(C.rooms.list())
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
	go usor.OnRun()
	return usor
}

func (C *Center) RoomNameList() []string {
	return C.rooms.list()
}

func (C *Center) Run() {
	//too verbose ?
	var u_serve = func(w http.ResponseWriter, r *http.Request, fn string) {
		if r.Method == "GET" {
			w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
			w.Header().Set("Pragma", "no-cache")
			w.Header().Set("Expires", "0")
			http.ServeFile(w, r, "./frontend/"+fn)
		} else {
			http.Error(w, "Bad Request", 404)
		}
	}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			u_serve(w, r, "index.html")
		} else {
			http.Error(w, "Bad Request", 404)
		}
	})
	http.HandleFunc("/uclient.js", func(w http.ResponseWriter, r *http.Request) {
		u_serve(w, r, "uclient.js")
	})
	http.HandleFunc("/ui.js", func(w http.ResponseWriter, r *http.Request) {
		u_serve(w, r, "ui.js")
	})
	http.HandleFunc("/ustyle.css", func(w http.ResponseWriter, r *http.Request) {
		u_serve(w, r, "ustyle.css")
	})
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		_ulog("@std@", "Center.run()", "/ws")
		C.eden.AddUsor(C.newUsor(w, r))
	})
	http.ListenAndServe(":9999", nil)
}
