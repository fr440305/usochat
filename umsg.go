//umsg.go
//This source file defined type Msg.

package main

//import "html"
//import "encoding/json"

const SET_DESCRIPTION byte = 1
const SET_CONTENT byte = 2

type Msg struct {
	usor    *Usor
	room    *Room
	center  *Center
	summary string
	content [][]string
}

func newMsg(source_usor *Usor, source_room *Room) *Msg {
	return nil
}

func (M *Msg) clone() *Msg {
	return nil
}

func (M *Msg) toRoom(dest *Room) {
}

func (M *Msg) toUsor(dest *Usor) {
}

func (M *Msg) toCenter(dest *Center) {
}

func (M *Msg) set(set_what byte, description string, content []string) *Msg {
	return nil
}

func (M *Msg) parseJSON(json_raw string) error {
	return nil
}

func (M *Msg) jsonify() string {
	return ""
}

func (M Msg) Error() string {
	return ""
}
