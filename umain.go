package main

import "syscall"

func main() {
	_ulog("_main", "http://127.0.0.1:9999")
	newCenter(syscall.Getpid()).Run()
}
