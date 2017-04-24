//umsg.go
//This source file defined type Msg.

package main

import "html"
import "encoding/json"

const SET_DESCRIPTION byte = 1
const SET_CONTENT byte = 2

type Msg struct {
	usor    *Usor
	room    *Room
	meta    *Msg
	summary string
	content []string
}

func newMsg(source_usor *Node) *Msg {
	return nil
}

func (M *Msg) clone() *Msg {
}

func (M *Msg) set(set_what byte, description string, content []string) *Msg {
}

func (M *Msg) parseJSON(json_raw string) error {
}

func (M *Msg) toJSON() string {
}

func (M Msg) Error() string {
}
