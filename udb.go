package main

import "fmt"

var if_log bool = true

func _ulogSet(iflog bool) {
	if_log = iflog
}

func _ulog(G ...interface{}) {
	if if_log == true {
		fmt.Println(G...)
	}
}

func _usor(usor *Usor) {
	if usor == nil {
		return
	}
	var pub_usor = struct {
		Nid      uint64
		R        *Room
		MsgQueue chan *Msg
	}{usor.nid, usor.room, usor.msg_queue}
	fmt.Println(pub_usor)
}

func _usorArr(usor_arr []*Usor) {
	for _, each_usor := range usor_arr {
		_usor(each_usor)
	}
}

func _uroom(room *Room) {

}

//I have not determined what is the duty of this type yet.
type Debugger struct {
}
