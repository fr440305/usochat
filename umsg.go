// USE MANY MANY TINY THEARDS TO SEND OR RECEIVE EVENTS.
// USE SLICE EVERYWHERE - DO NOT USE container/list.
// 不光要思考架构，还要思考架构的迭代与演化。
// 要牢记：软件是长出来的。

//CODE_COMPLETE:
// --all TODOs & FIXMEs
// - documentation: on business logic.

package main

import "fmt"
import "html"
import "encoding/json"

type Msg struct {
	source_node *Node // != nil => from node. ==nil => from center.
	description string
	content     []string
}

func newMsg(source_node *Node) *Msg {
	var res = new(Msg)
	res.source_node = source_node
	res.description = ""
	res.content = []string{}
	fmt.Println("newMsg", res)
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
		fmt.Println("Msg.setContent", content[i])
	}
	M.content = content
	fmt.Println("Msg.setContent", content)
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
	fmt.Println("Msg.parseJSON", user_msg)
	fmt.Println("Msg.parseJSOn - end.")
	return nil
}

//TODO - check error
//This method transforms the Msg::M to JSON string.
func (M *Msg) toJSON() string {
	var res []byte
	var err error
	fmt.Println("Msg.toJSON", "begin")
	var user_msg = struct {
		SouceNode   string   `json:"source_node"`
		Description string   `json:"description"`
		Content     []string `json:"content"`
	}{M.source_node.iden, M.description, M.content}
	fmt.Println("Msg.toJSON", user_msg)
	res, err = json.Marshal(user_msg)
	if err != nil {
		//TODO - error handler goes here...
	}
	fmt.Println("Msg.toJSON", user_msg)
	fmt.Println("Msg.toJSOn - end.", string(res))
	return string(res)
	//return `{"content":["toJSON","toJSON"]}`
}

func (M Msg) Error() string {
	if M.description == "error" && M.content != nil && len(M.content) != 0 {
		return M.content[0]
	}
	return ""
}
