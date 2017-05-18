package main

import "net/http"

func main() {
	_ulog("_main", "http://127.0.0.1:9999")
	http.ListenAndServe(":9999", newUserver().Mux())
}
