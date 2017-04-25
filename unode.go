//unode.go
//This go source file defined three go type: Usor, Room, and Center.
//Author(s): __HUO_YU__

package main

import "github.com/gorilla/websocket"
import "net/http"
import "strconv"

//type Usor maps to a client.
type Usor struct {
	msgchan chan Msg
	room    *Room
	conn    *websocket.Conn //connent client to node
	nid     int64           //node id
}

func newUsor(upgdr websocket.Upgrader, w http.ResponseWriter, r *http.Request) *Usor {

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
	msg_queue chan Msg
	msg_hist  []Msg
	members   []Usor
	center    *Center
}

func newRoom() {
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

type Center struct {
	rooms       []Room
	ws_upgrader websocket.Upgrader //const
}

func newCenter() *Center {
	return nil
}

func (C *Center) run() {
}
