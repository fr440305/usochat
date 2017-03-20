//BUG - After the client exit, the number in center will not reduce.
//TODO - Add a onclose event in frontend.

// USE MANY MANY TINY THEARDS TO SEND OR RECEIVE EVENTS.

package main

import "fmt"
import "net/http"
import "strconv"
import "github.com/gorilla/websocket"

type Center struct {
	event_queue chan map[string]string
	upgrader    websocket.Upgrader //Constant
	num_onliner int
}

func newCenter() *Center {
	var res = new(Center)
	res.event_queue = make(chan map[string]string)
	res.upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	return res
}

//return the number of people online:
func (c *Center) GetOnliner() []byte {
	//convert c.num_online, which is int, to a byte array:
	return []byte(strconv.Itoa(c.num_onliner))
}

func (c *Center) AddOnliner(increment int) {
	if c.num_onliner < increment {
		c.num_onliner = 0
	} else {
		c.num_onliner += increment
	}
}

func (c *Center) Listen() {
}

func (c *Center) ServeWs(w http.ResponseWriter, r *http.Request) {
	//Initialization:
	//create the connection:
	conn, err := c.upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	fmt.Println("/ws - accessing....")
	c.AddOnliner(1)
	//Polling:
	for {
		//Read message from the client:
		//code will be blocked here until it received msg:
		msg_type, msg_cx, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("Fatal - conn--readmsg")
			return
		}
		fmt.Println(msg_type, msg_cx)
		//Write the messages(how many people onlines) to the client:
		err = conn.WriteMessage(websocket.TextMessage, c.GetOnliner())
		if err != nil {
			fmt.Println("Fatal - conn--response")
			return
		}
		fmt.Println("conn--response....!")
	}
}

func main() {
	fmt.Println("http://127.0.0.1:9999")
	var center = newCenter()
	//To serve the webpages to the client:
	http.Handle("/", http.FileServer(http.Dir(".")))
	//To handle the websocket request:
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		//center.newNode().serveWs(w, r)
		center.ServeWs(w, r)
	})
	http.ListenAndServe(":9999", nil)
}
