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

type GetHandler struct {
}

func (GH GetHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

type PostHandler struct {
}

func (PH PostHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

func main() {
	/* File Server */
	//http.HandleFunc("/", MapFile(""))
	//http.HandleFunc("/index.html", MapFile("index.html"))
	//http.HandleFunc("/app.js", MapFile("app.js"))
	//http.HandleFunc("/api.js", MapFile("api.js"))
	//http.HandleFunc("/get", GetHandler)
	//http.HandleFunc("/post", PostHandler)

	go func() {
		log.Fatal(http.ListenAndServe(":8081", new(GetHandler)))
	}()
	go func() {
		log.Fatal(http.ListenAndServe(":8082", new(PostHandler)))
	}()
	log.Fatal(http.ListenAndServe(":8080", http.FileServer(http.Dir("."))))
	fmt.Println("vim-go")
}
