package main

import "syscall"

func main() {
	_ulog("_main", "http://127.0.0.1:9999")
	var u_center = newCenter(syscall.Getpid())
	u_center.run()
}
