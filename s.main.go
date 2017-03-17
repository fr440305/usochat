package main

import "fmt"
import "net/http"
import "github.com/gorilla/websocket"

func main() {
	fmt.Println("http://127.0.0.1:9999")

	http.Handle("/", http.FileServer(http.Dir(".")))
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("/ws - accessing....")
	})
	http.ListenAndServe(":9999", nil)
}
