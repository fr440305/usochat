//umsg.go
//This source file defined type Msg.

package main

import "html"
import "encoding/json"

const SET_DESCRIPTION byte = 1
const SET_CONTENT byte = 2

type Msg struct {
	usor        *Usor
	room        *Room
	description string
	content     []string
}

func newMsg(source_usor *Node) *Msg {
	return nil
}

func (M *Msg) clone() *Msg {
}

func (M *Msg) set(set_what byte, description string, content []string) *Msg {
}

//Pay attention to the probobaly-appear errors.
func (M *Msg) parseJSON(json_raw string) error {
	var user_msg struct {
		SouceNode   string   `json:"source_node"`
		Description string   `json:"description"`
		Content     []string `json:"content"`
	}
	json.Unmarshal([]byte(json_raw), &user_msg)
	M.set(SET_DESCRIPTION|SET_CONTENT, user_msg.Description, user_msg.Content)
	return nil
}

//TODO - check error
//This method transforms the Msg::M to JSON string.
func (M *Msg) toJSON() string {
	var res []byte
	var err error
	var source_node_iden string
	if M.source_node == nil {
		source_node_iden = ""
	} else {
		source_node_iden = M.source_node.idString()
	}
	var user_msg = struct {
		SouceNode   string   `json:"source_node"`
		Description string   `json:"description"`
		Content     []string `json:"content"`
	}{source_node_iden, M.description, M.content}
	res, err = json.Marshal(user_msg)
	if err != nil {
		//TODO - error handler goes here...
	}
	return string(res)
	//return `{"content":["toJSON","toJSON"]}`
}

func (M Msg) Error() string {
	if M.description == "error" && M.content != nil && len(M.content) != 0 {
		return M.content[0]
	}
	return ""
}
