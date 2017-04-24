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
	nid     int64           // the identification for node.
}

//The following function returns the string form of the node id.
func (U *Usor) Get(get_what string) (int64, string) {
	//TODO//
	return ""
}

func (U *Usor) handleUser(ifexit chan<- bool) {
}

//responser
//fetch the msg from center, and send it to client.
func (U *Usor) handleRoom() {
}

func (U *Usor) Run() {
}

type Room struct {
	rid       uint64
	msg_queue chan Msg
	msg_hist  []Msg
	members   []Usor
	center    Center
}

func (R *Room) newNode() {
}

func (R *Room) removeUsor(rm_usr *Usor) error {
}

func (R *Room) boardcast(bcmsg Msg) error {
}

func (R *Room) Get(get_what string) (int64, string) {
}

func (R *Room) Run() {
}

func (R *Room) PushMsg(m Msg) error {
}

type Center struct {
	rooms       []Room
	ws_upgrader Websocket.Upgrader //const
}

func (C *Center) NewRoom() *Room {
	return nil
}
