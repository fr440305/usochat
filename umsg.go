//umsg.go
//This source file defined type Msg.

package main

//import "html"
//import "encoding/json"

const SET_DESCRIPTION byte = 1
const SET_CONTENT byte = 2

//The instance of Msg can only be constructed by Usor.
type Msg struct {
	usor    *Usor
	eden    *Eden
	room    *Room
	center  *Center
	summary string
	content [][]string
}

func newMsg(usor *Usor, eden *Eden, room *Room, center *Center, summary string, content [][]string) *Msg {
	return &Msg{
		usor:    usor,
		eden:    eden,
		room:    room,
		center:  center,
		summary: summary,
		content: content,
	}
}

func (M *Msg) clone() *Msg {
	return nil
}

func (M *Msg) parse(json_raw string) error {
	return nil
}

func (M *Msg) jsonify() string {
	return ""
}

func (M Msg) Error() string {
	if M.summary == "error" {
		return M.content[0][0]
	} else {
		return "" //not an error
	}
}

type MsgList []*Msg
