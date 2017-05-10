package main

import "fmt"

var if_log bool = true

func _ulogSet(iflog bool) {
	if_log = iflog
}

func _ustd(G ...interface{}) {
	_ulog("@std@", G...)
}

func _uerr(G ...interface{}) {
	_ulog("@err@", G...)
}

func _ulog(G ...interface{}) {
	if if_log == true {
		fmt.Println(G...)
	}
}

func _usor(usor *Usor) {
}

func _usorArr(usor_arr []*Usor) {
}

func _uroom(room *Room) {

}

//I have not determined what is the duty of this type yet.
type Debugger struct {
}
