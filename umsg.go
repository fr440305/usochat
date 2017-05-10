//umsg.go
//This source file defined type Msg.

package main

//import "html"
import "encoding/json"

const SET_DESCRIPTION byte = 1
const SET_CONTENT byte = 2

//The instance of Msg can only be constructed by Usor.
type Msg struct {
	Summary string
	Content [][]string
}

func newMsg(summary string, content [][]string) *Msg {
	return &Msg{
		Summary: summary,
		Content: content,
	}
}

func newBarMsg(barjson []byte) *Msg {
	var resmsg = new(Msg)
	//...
	err := json.Unmarshal(barjson, resmsg)
	if err != nil {
		_ulog("@err@ Msg.newBarJson", err.Error())
		return nil
	}
	return resmsg
}

func newErrMsg(errinfo string) *Msg {
	return newMsg("error", [][]string{[]string{errinfo}})
}

func (M *Msg) clone() *Msg {
	return nil
}

func (M *Msg) barjsonify() []byte {
	res, err := json.Marshal(M)
	if err != nil {
		_ulog("@err@ Msg.barjsonify", err.Error())
		return []byte{}
	}
	return res
}

func (M Msg) Error() string {
	if M.Summary == "error" {
		return M.Content[0][0]
	} else {
		return "" //not an error
	}
}

type MsgList []*Msg
