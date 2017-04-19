//umsg.go
//This source file defined type Msg.

package main

import "html"
import "encoding/json"

type Msg struct {
	usor        *Usor
	room        *Room
	description string
	content     []string
}

func newMsg(source_usor *Node) *Msg {
	var res = new(Msg)
	res.usor = source_usor
	res.room = source_node.room
	res.description = ""
	res.content = []string{}
	//_ulog("newMsg", res)
	return res
}

//new_msg_type = '0' | '*' | ' '
//'0' for response message, '*' for boardcast message.
// ' ' for original message.
func (M *Msg) msgCopy(new_msg_type byte) *Msg {
	var new_description string
	if new_msg_type == ' ' {
		new_description = M.description
	} else {
		new_description = string(append([]byte(M.description), '-', new_msg_type))
	}
	return &Msg{
		source_node: M.source_node,
		description: new_description,
		content:     M.content[:],
	}
}

func (M *Msg) setDescription(description string) *Msg {
	M.description = html.EscapeString(description)
	return M
}

func (M *Msg) setContent(content []string) *Msg {
	for i, str := range content {
		content[i] = html.EscapeString(str)
	}
	M.content = content
	return M
}

//Pay attention to the probobaly-appear errors.
func (M *Msg) parseJSON(json_raw string) error {
	var user_msg struct {
		SouceNode   string   `json:"source_node"`
		Description string   `json:"description"`
		Content     []string `json:"content"`
	}
	json.Unmarshal([]byte(json_raw), &user_msg)
	M.setDescription(user_msg.Description)
	M.setContent(user_msg.Content)
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
