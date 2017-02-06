package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/foo", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.URL.Path)
	})
	http.ListenAndServe(":8080", http.FileServer(http.Dir(".")))
	fmt.Println("vim-go")
}
