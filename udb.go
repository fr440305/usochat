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

//I have not determined what is the duty of this type yet.
type Debugger struct {
}
