package main

import (
	"net/http"
)

func ServeHome() func(http.ResponseWriter, *http.Request) {
	return (func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.Error(w, "File Not Found", 404)
			return
		}
		if r.Method != "GET" {
			http.Error(w, "Bad Access Method", 405)
			return
		}
		http.ServeFile(w, r, "index.html")
	})
}

func ServeFile(filename string) func(http.ResponseWriter, *http.Request) {
	return (func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filename)
	})
}
