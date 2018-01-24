package main

import "strconv"
import "github.com/fr440305/uso"
import "log"
import "net/http"
import "os"
import "fmt"

func help() {
	fmt.Println(`
Intro:
"usod" is a little web server providing instant messaging function.

Usage:
$ usod
$ usod <portnum>

`)
}

func main() {
	var portnum string
	switch len(os.Args) {
	case 1:
		portnum = "9999"
	case 2:
		portnum = os.Args[1]
		if _, err := strconv.Atoi(portnum); err != nil {
			log.Println("main::OsArgErr Portnum should be {0 .. 65535}")
			help()
			return
		}
	default:
		log.Println("main::OsArgErr Length of os.Args should be either one or two.")
		help()
		return
	}
	if _, err := os.Stat("test.html"); err != nil {
		log.Fatal("main::FileNotFound ", err)
	}
	fmt.Println("Listen on port", portnum, "...")
	fmt.Println("Visit http://127.0.0.1:"+portnum, "to see the result")


	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "test.html")
	})
	http.HandleFunc("/uso/conn", uso.ServeWs)
	http.ListenAndServe(":"+portnum, nil)
}
