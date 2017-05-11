package main

import "fmt"

func _ulog(G ...interface{}) {
	fmt.Println(G...)
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
