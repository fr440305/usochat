//unode.go
//This go source file defined three go type: Usor, Room, and Center.
//Author(s): __HUO_YU__

package main

import "github.com/gorilla/websocket"
import "net/http"

//import "strconv"

//type Usor maps to a client.
type Usor struct {
	nid       uint64 //node id
	msg_queue chan *Msg
	eden      *Eden
	room      *Room
	conn      *websocket.Conn //client <--conn--> node
}

func (U *Usor) newMsg() *Msg {
	return nil
}

func (U *Usor) handleClient() {
	_ulog("Usor.handleClient")
	var msgtype int
	var barjson []byte //bar = byte array
	var strjson string
	var err error
	//var msg *Msg
	for {
		msgtype, barjson, err = U.conn.ReadMessage()
		_ulog("Usor.HandleClient.middle")
		if err != nil {
			//Gone.
			_ulog("@err@", "Usor.handleClient", err.Error())
			U.conn.Close()
			return
		} else {
			if msgtype == websocket.TextMessage {
				strjson = string(barjson)
				_ulog("@std@", "Usor.handleClient type=", msgtype, strjson)
			} else if msgtype == websocket.BinaryMessage {
				_ulog("@std@", "Usor.handleClient type=", msgtype, barjson)
			} else if msgtype == websocket.CloseMessage {
				_ulog("@std@", "Usor.handleClient type=", msgtype, strjson)
			} else { //Unexpected Message.
				_ulog("@err@", "Usor.handleClient type=", msgtype, strjson)
			}
		}
	}
}

func (U *Usor) handleRoom() {
	var msg *Msg
	select {
	case msg = <-U.msg_queue:
		//if it is a logout-ok msg, then return.
		_ulog("@std@", "Usor.handleRoom", msg.jsonify())
	}
}

func (U *Usor) Run() {
	go U.handleClient()
	U.handleRoom()
}

type Room struct {
	rid       uint64
	name      string
	msg_queue chan *Msg
	msg_hist  []*Msg
	members   []*Usor
	center    *Center
}

func (R *Room) addUsor(usor *Usor) error {
	return nil
}

func (R *Room) rmUsor(rm_usr *Usor) error {
	return nil
}

func (R *Room) boardcast(bcmsg Msg) error {
	return nil
}

func (R *Room) handleCenter() {
}

func (R *Room) handleUsors() {
}

func (R *Room) Run() {
}

type Eden struct {
	center    *Center
	msg_queue chan *Msg
	members   []*Usor
}

func (E *Eden) AddUsor(usor *Usor) *Usor {
	E.members = append(E.members, usor)
	return usor
}

func (E *Eden) Run() {
}

//The main server
type Center struct {
	pid         int //process id
	msg_queue   chan *Msg
	eden        *Eden
	rooms       []*Room
	usors       []*Usor
	ws_upgrader websocket.Upgrader //const
}

func newCenter(pid int) *Center {
	var center = new(Center)
	center.pid = pid
	center.msg_queue = make(chan *Msg)
	center.newEden()
	_ulog("@pid@", pid)
	return center
}

func (C *Center) newEden() *Eden {
	var eden = new(Eden)
	eden.members = []*Usor{}
	eden.center = C
	go eden.Run()
	C.eden = eden
	_ulog("@std@", "Center.newEden")
	return eden
}

func (C *Center) newRoom(name string) *Room {
	var room = new(Room)
	room.rid = C.validRoomId()
	room.name = name
	room.msg_queue = make(chan *Msg)
	room.msg_hist = []*Msg{}
	room.members = []*Usor{}
	room.center = C
	C.rooms = append(C.rooms, room)
	_ulog("@dat@", "Center.newRoom", C.rooms)
	go room.Run()
	return room
}

func (C *Center) newUsor(w http.ResponseWriter, r *http.Request) *Usor {
	var usor = new(Usor)
	var err error
	usor.nid = C.validUsorId()
	usor.msg_queue = make(chan *Msg)
	usor.eden = C.eden
	usor.room = nil
	usor.conn, err = C.ws_upgrader.Upgrade(w, r, w.Header())
	if err != nil {
		_ulog("@err@", "Center.newUsor", err.Error())
		return nil
	}
	C.usors = append(C.usors, usor)
	_ulog("@std@", "Center.newUsor")
	_usorArr(C.usors)
	_ulog("Center.newUsor", "tail")
	go usor.Run()
	return usor
}

func (C Center) validRoomId() uint64 {
	return 0
}

func (C Center) validUsorId() uint64 {
	return 0
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
	http.ListenAndServe(":9999", nil) //go func(){}
	C.handleRooms()
}
