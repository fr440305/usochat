//unode.go
//This go source file defined three go type: Usor, Room, and Center.
//Author(s): __HUO_YU__

package main

import "github.com/gorilla/websocket"
import "net/http"

//import "strconv"

//type Usor maps to a client.
type Usor struct {
	msgchan chan Msg
	room    *Room
	conn    *websocket.Conn //connent client to node
	nid     int64           //node id
}

func newUsor(upgdr websocket.Upgrader, w http.ResponseWriter, r *http.Request) *Usor {
	return nil
}

//eg - ("id")-->(0, "0");
func (U *Usor) get(get_what string) (int64, string) {
	//TODO//
	return 0, ""
}

func (U *Usor) run() {
}

type Room struct {
	rid       uint64
	name      string
	msg_queue chan Msg
	msg_hist  []Msg
	members   []Usor
	center    *Center
}

func newRoom() *Room {
	return nil
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
	msg_q       chan Msg
	rooms       []Room
	ws_upgrader websocket.Upgrader //const
}

func newCenter(pid int) *Center {
	var center = new(Center)
	center.pid = pid
	center.msg_q = make(chan Msg)
	center.rooms = []Room{}
	return center
}

func (C *Center) handleRooms() {
}

func (C *Center) run() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_ulog("Center.run()", "ennn...")
	})
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		//C.rooms[0].newUsor(C.validUsorId()).run()
	})
	http.ListenAndServe(":9999", nil)
}
