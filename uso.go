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

func GetHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	/* for testing */
	fmt.Println("<!--get")
	fmt.Println("  Method = ", r.Method)
	fmt.Println("  RawQuery = ", r.URL.RawQuery)
	for d, r := range r.Form {
		fmt.Println("  Form[", d, "] |-> ", r)
		if d == "conversation" {
			fmt.Println("  conversation:")
			/* TODO - pretend XSS attack, pay attention to security */
			resp_json, err := json.Marshal(dialogs)
			if err != nil {
				fmt.Println("http-get-Fatal!!-json.Marshal")
				return
			}
			w.Header().Set("Content-Type", "text/json")
			w.WriteHeader(http.StatusOK)
			w.Write(resp_json)
			for index, diaelm := range dialogs {
				fmt.Println("    dialog[", index, "] = ", diaelm)
			}
		}
	}
	fmt.Println("-->")
}

func PostHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	/* for testing */
	fmt.Println("<!--post")
	if r.Method == "POST" {
		fmt.Println("  Method = ", r.Method, ", normal")
	} else {
		fmt.Println("  Method = ", r.Method, ", ALERT!")
	}
	for k, v := range r.PostForm {
		fmt.Println("  PostForm[", k, "] |-> ", v)
		if k == "dialog" {
			dialogs = append(dialogs, strings.Join(v, ""))
		}
	}
	fmt.Println("-->")
}

type WebPg struct {
}

func (wpg WebPg) ServeHome() func(http.ResponseWriter, *http.Request) {
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

func (wpg WebPg) ServeFile(filename string) func(http.ResponseWriter, *http.Request) {
	return (func(w http.ResponseWriter, r *http.Request) {
	})
}

func main() {
	/* File Server */
	http.HandleFunc("/", WebPg.ServeHome())
	http.HandleFunc("/index.html", WebPg.ServeFile("index.html"))
	http.HandleFunc("/app.js", WebPg.ServeFile("app.js"))
	http.HandleFunc("/api.js", WebPg.ServeFile("api.js"))
	http.HandleFunc("/get", GetHandler)
	http.HandleFunc("/post", PostHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
	fmt.Println("vim-go")
}
