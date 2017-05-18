package main

func main() {
	_ulog("_main", "http://127.0.0.1:9999")
	newUserver().ListenAndServe(":9999")
}
