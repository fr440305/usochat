package main

import "fmt"
import "net/http"
import "github.com/fr440305/uso"

func main() {
	fmt.Println("_main", "http://127.0.0.1:9999/app.uso")
	http.ListenAndServe(":9999", uso.NewUserver().Mux())
}
