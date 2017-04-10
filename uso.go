// USE MANY MANY TINY THEARDS TO SEND OR RECEIVE EVENTS.
// USE SLICE EVERYWHERE - DO NOT USE container/list.
// 不光要思考架构，还要思考架构的迭代与演化。
// 要牢记：软件是长出来的。

//CODE_COMPLETE:
// --all TODOs & FIXMEs
// - documentation: on business logic.

package main

import "net/http"
import "syscall"

func main() {
	var center = newCenter()
	var pid = syscall.Getpid()
	_ulogSet(true)
	_ulog("_main", "http://127.0.0.1:9999")
	_ulog("@pm", "pid = ", pid)
	go center.handleNodes()
	//To provide the webpages to the client:
	http.Handle("/", http.FileServer(http.Dir("./frontend")))
	//To handle the websocket request:
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		var if_node_exit = make(chan bool)
		go center.newNode(w, r).run(if_node_exit)
		select {
		case <-if_node_exit:
			_ulog("_main", "A node exit.")
			return
		}
	})
	http.ListenAndServe(":9999", nil)
}
