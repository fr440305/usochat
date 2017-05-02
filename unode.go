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
	conn *websocket.Conn //client <--conn--> node
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
				if strjson[0] == '{' {
					//json-form
				} else {
					//raw-form
				}
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

func (U *Usor) Run() {
	U.handleClient()
}

func (U *Usor) BeSent(msg *Msg) *Msg {
	return msg
}

type UsorList []*Usor

func (UL *UsorList) add(usor *Usor) *Usor {
	*UL = append(*UL, usor)
	return usor
}

func (UL *UsorList) rm(usor *Usor) *Usor {
	return nil
}

func (UL UsorList) boardcast(msg *Msg) *Msg {
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

func (R *Room) Run() {
}

type RoomList []*Room

func (RL *RoomList) add(room *Room) *Room {
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
	return nil
}

func (RL RoomList) usorAmount() int64 {
	return 0
}

func (RL RoomList) roomAmount() int64 {
	return int64(len(RL))
}

type Eden struct {
	center *Center
	guests *UsorList
}

func (E *Eden) AddUsor(usor *Usor) {
	E.guests.add(usor)
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

func (C *Center) newRoom(name string) *Room {
	var room = new(Room)
	room.name = name
	room.qmsg = new(MsgList)
	room.chist = [][]string{}
	room.usors = new(UsorList)
	room.center = C
	_ulog("@dat@", "Center.newRoom", C.rooms)
	return room
}

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
