package main

import "github.com/fr440305/uso"
import "net/http"

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "test.html")
	})
	http.HandleFunc("/uso/conn", uso.ServeWs)
	http.ListenAndServe(":9999", nil)
}
