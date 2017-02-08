package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

func MapFile(file_name string) func(http.ResponseWriter, *http.Request) {
	if file_name == "" {
		return (func(res http.ResponseWriter, req *http.Request) {
			if req.URL.Path == "/" {
				fmt.Println("right") /* for index.html */
				/* use Library::io to get the address of file */
				file_content, err := ioutil.ReadFile("index.html")
				if err != nil {
					/* Write 404 into ResponseWriter */
					res.WriteHeader(http.StatusNotFound)
					return
				}
				/* Write file_content into ResponseWrite */
				res.WriteHeader(http.StatusOK)
				io.WriteString(res, string(file_content))
				/* for testing */
				fmt.Println("Route: index.html --> ", req.URL.Path)
			} else {
				fmt.Println("wrong") /* 404 */
				res.WriteHeader(http.StatusNotFound)
				return
			}
		})
	} else {
		return (func(res http.ResponseWriter, req *http.Request) {
			/* use Library::io to get the address of file */
			file_content, err := ioutil.ReadFile(file_name)
			if err != nil {
				/* Write 404 into ResponseWriter */
				res.WriteHeader(http.StatusNotFound)
				return
			}
			/* Write file_content into ResponseWrite */
			res.WriteHeader(http.StatusOK)
			io.WriteString(res, string(file_content))
			/* for testing */
			fmt.Println("Route:", file_name, "--> ", req.URL.Path)
		})
	}
}

func main() {
	http.HandleFunc("/", MapFile(""))
	http.HandleFunc("/index.html", MapFile("index.html"))
	http.HandleFunc("/app.js", MapFile("app.js"))
	http.HandleFunc("/api.js", MapFile("api.js"))
	http.HandleFunc("/get", func(w http.ResponseWriter, r *http.Request) {
		/* for testing */
		/* TODO - handle query */
		fmt.Println("<!--get")
		fmt.Println("  Method = ", r.Method)
		fmt.Println("  RawQuery = ", r.URL.RawQuery)
		fmt.Println("-->")
	})
	http.HandleFunc("/post", func(w http.ResponseWriter, r *http.Request) {
		/* for testing */
		fmt.Println("<!--post")
		if r.Method == "POST" {
			fmt.Println("  Method = ", r.Method, ", normal")
		} else {
			fmt.Println("  Method = ", r.Method, ", ALERT!")
		}
		fmt.Println("-->")
	})
	http.ListenAndServe(":8080", nil)
	fmt.Println("vim-go")
}
