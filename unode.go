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
	room      *Room
	conn      *websocket.Conn //client <--conn--> node
}

func (U *Usor) newMsg() *Msg {
	return nil
}

//eg - ("id")-->(0, "0");
func (U *Usor) get(get_what string) (int64, string) {
	//TODO//
	return 0, ""
}

func (U *Usor) handleClient() {
	var msgtype int
	var barjson []byte //bar = byte array
	var strjson string
	var err error
	//var msg *Msg
	for {
		msgtype, barjson, err = U.conn.ReadMessage()
		if err != nil {
			_ulog("@err@", "Usor.handleClient", err.Error())
			return
		} else {
			if msgtype == websocket.TextMessage {
				strjson = string(barjson)
				_ulog("@std@", "Usor.handleClient", msgtype, strjson)
			} else if msgtype == websocket.BinaryMessage {
				_ulog("@std@", "Usor.handleClient", msgtype, barjson)
			} else if msgtype == websocket.CloseMessage {
				_ulog("@std@", "Usor.handleClient", msgtype, strjson)
			} else {
				_ulog("@std@", "Usor.handleClient", msgtype, strjson)
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

func (U *Usor) run() {
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

func (R *Room) rmUsor(rm_usr *Usor) error {
	return nil
}

func (R *Room) boardcast(bcmsg Msg) error {
	return nil
}

func (R *Room) get(get_what string) (int64, string) {
	return 0, ""
}

func (R *Room) run() {
}

func (R *Room) push(m Msg) error {
	return nil
}

//The main server
type Center struct {
	pid         int //process id
	msg_queue   chan *Msg
	rooms       []*Room
	usors       []*Usor
	ws_upgrader websocket.Upgrader //const
}

func newCenter(pid int) *Center {
	var center = new(Center)
	center.pid = pid
	center.msg_queue = make(chan *Msg)
	center.newRoom("Eden")
	_ulog("@pid@", pid)
	return center
}

func (C Center) validRoomId() uint64 {
	return 0
}

func (C Center) validUsorId() uint64 {
	return 0
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
	return room
}

func (C *Center) newUsor(room *Room, w http.ResponseWriter, r *http.Request) *Usor {
	var usor = new(Usor)
	var err error
	usor.nid = C.validUsorId()
	usor.msg_queue = make(chan *Msg)
	usor.room = room
	usor.conn, err = C.ws_upgrader.Upgrade(w, r, w.Header())
	if err != nil {
		_ulog("@err@", "Center.newUsor", err.Error())
		return nil
	}
	C.usors = append(C.usors, usor)
	_ulog("@std@", "Center.newUsor", C.usors)
	return usor
}

//return that how many time it sent.
func (C *Center) boardcast() int {
	return 0
}

func (C *Center) handleRooms() error {

	return nil
}

func (C *Center) run() {
	http.Handle("/", http.FileServer(http.Dir("frontend")))
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		_ulog("@std@", "Center.run()", "/ws")
		C.newUsor(C.rooms[0], w, r).run()
	})
	http.ListenAndServe(":9999", nil) //go func(){}
	C.handleRooms()
}
