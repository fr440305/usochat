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

func _usor(usor *Usor) string {
	return ""
}

func _uroom(room *Room) string {
	return ""
}

//I have not determined what is the duty of this type yet.
type Debugger struct {
}
