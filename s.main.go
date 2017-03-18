package main

import "fmt"
import "net/http"
import "strconv"
import "github.com/gorilla/websocket"

type Center struct {
	num_onliner int
}

func newCenter() *Center {
	return new(Center)
}

func (c *Center) GetOnliner() []byte {
	return []byte(strconv.Itoa(c.num_onliner))
}

func (c *Center) AddOnliner() {
	c.num_onliner += 1
}

func main() {
	fmt.Println("http://127.0.0.1:9999")
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	var center = newCenter()
	http.Handle("/", http.FileServer(http.Dir(".")))
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		//create the connection:
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		fmt.Println("/ws - accessing....")
		center.AddOnliner()
		//Polling:
		for {
			//Read message from the client:
			//code will be blocked here until it received msg:
			msg_type, msg_cx, err := conn.ReadMessage()
			if err != nil {
				fmt.Println("Fatal - conn-readmsg")
				return
			}
			fmt.Println(msg_type, msg_cx)
			//Write the messages to the client:
			err = conn.WriteMessage(websocket.TextMessage, center.GetOnliner())
			if err != nil {
				fmt.Println("Fatal - conn--response")
				return
			}
			fmt.Println("conn--response....!")
		}
	})
	http.ListenAndServe(":9999", nil)
}
