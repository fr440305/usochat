package main

import (
	"encoding/json"
	"fmt"
	"log"
	//"io"
	//"io/ioutil"
	"net/http"
	"strings"
)

var dialogs []string

func main() {
	/* File Server */
	http.HandleFunc("/", ServeHome())
	http.HandleFunc("/index.html", ServeFile("index.html"))
	http.HandleFunc("/app.js", ServeFile("app.js"))
	http.HandleFunc("/api.js", ServeFile("api.js"))
	http.HandleFunc("/get", GetHandler)
	http.HandleFunc("/post", PostHandler)

	log.Fatal(http.ListenAndServe(":8888", nil))
	fmt.Println("vim-go")
}
