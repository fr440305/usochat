package uso // import "github.com/fr440305/uso"

import "github.com/gorilla/websocket"
import "net/http"
import "encoding/json"
import "fmt"

const (
	// ["u->h"]
	MsgType_uso_to_hall = "u->h"

	// ["u->r", room-name]
	MsgType_uso_to_room = "u->r"

	// ["u++r", room-name]
	MsgType_uso_new_room = "u++r"

	// ["uerr", err-situation]
	MsgType_uso_err = "uerr"

	// ["uh->"]
	MsgType_uso_exit_hall = "uh->"
	MsgType_uso_exit_room = "ur->"
	MsgType_uso_quit      = "u->x"
	// ...
)

var Uso_connpool = []Usoconn{}
var Uso_Roompool = []Room{}

func getRoomByName(name string, roompool []Room) *Room {
	for i, r := range Uso_Roompool {
		if r.Name == name {
			return &Uso_Roompool[i]
		}
	}
	return nil
}

var Uso_websocket_upgrader = &websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Message []string

// string -> []interface{} // tokenlize
// []interface{} -> UsoMsg // done
// use json.Unmarshal
func NewMessage(bar []byte) Message {
	strarr := &[]string{}
	err := json.Unmarshal(bar, strarr)
	if err != nil {
		fmt.Println(err.Error())
		return Message{""}
	}
	fmt.Println("NewMessage", string(bar))
	return Message(*strarr)
}

func (um Message) Type() string {
	if len([]string(um)) == 0 {
		return ""
	} else {
		fmt.Println("Type", um[0])
		return um[0]
	}
}

func (um Message) Cont() []string {
	if len([]string(um)) <= 1 {
		return nil
	} else {
		return um[1:]
	}
}

// return []byte(`["Type", "Cont-1", "Cont-2", ... , "Cont-n"]`)
// use json.Marshal
func (um Message) ToJbar() []byte {
	bar, err := json.Marshal(um)
	if err != nil {
		return []byte("[" + MsgType_uso_err + "]")
	} else {
		return bar
	}
}

type Usoconn struct {
	//if you are not in a Room, you do not need a name
	Name   string
	Quit   bool
	R      chan Message
	W      chan Message
	wsconn *websocket.Conn
}

// public.
func (uc Usoconn) Run() {
	// if err in wsconn then close R and close W.
	go func(rch chan Message) {
		//wsconn -> uc.R
		for {
			_, cont, err := uc.wsconn.ReadMessage()
			if err != nil {
				uc.Quit = true
				fmt.Println("Usoconn Run uso conn exit.")
				rch <- Message{MsgType_uso_quit}
				close(rch)
				return
			} else {
				fmt.Println("Usoconn Run", string(cont))
				rch <- NewMessage(cont)
			}
		}
	}(uc.R)
	go func(wch chan Message) {
		//uc.W -> wsconn
		for {
			msg := <-wch
			err := uc.wsconn.WriteMessage(
				websocket.TextMessage,
				msg.ToJbar(),
			)
			if err != nil {
				uc.Quit = true
				close(wch)
				return
			}
		}
	}(uc.W)
}

type Room struct {
	Name  string
	Hist  []string
	Conns []*Usoconn
	jobs  chan Message // if no one, runner will close it and quit.
}

func (r Room) run() {
	// open the job clannel
	// for a := range jobs {...}; close and quit
}

// public.
func (r Room) ServeMember(uc *Usoconn) {
	// jobs <-- for all ur.Conns.R
	// if ur.jobs is closed (ur.IsActive == false), then open it and run the runner
	// for a := range uc.R { ur.jobs <- a }; close and quit
}

type Hall struct {
	Guests    []*Usoconn
	IsRunning bool
	jobs      chan Message
}

func (h Hall) Run() {
	if h.IsRunning == false {
		go func() {
			h.IsRunning = true
			for msg := range h.jobs {
				// uso_to_hall | uso_to_room | uso_exit_hall |
				fmt.Println("Hall Run", msg)
			}
		}()
	}
}

func (h Hall) ServeGuest(uc *Usoconn) {
	h.Guests = append(h.Guests, uc)
	h.jobs <- Message{MsgType_uso_to_hall}
	//read uc.R
	for msg := range uc.R {
		// toroom | newroom
		fmt.Println("Hall ServeGuest", msg, msg.Type(), "type|cont", msg.Cont())
		switch msg.Type() {
		case MsgType_uso_to_room:
			fmt.Println("Hall ServeGuest to-room")

			room := getRoomByName(msg.Cont()[0], Uso_Roompool)
			if room == nil {
				uc.W <- Message{MsgType_uso_err, "no such room"}
			} else {
				room.ServeMember(uc)
				//then, they leave.
				// but not sure if us still active.
				h.Guests = append(h.Guests, uc)
				h.jobs <- Message{MsgType_uso_to_hall}
			}
		case MsgType_uso_new_room:
		default:
		}
	}
	// uc.R closed:
	h.jobs <- Message{MsgType_uso_exit_hall}
	//remove
}

var hall = Hall{
	Guests:    []*Usoconn{},
	IsRunning: false,
	jobs:      make(chan Message),
}

// export
func Run() {
	fmt.Println("Run")
	go hall.Run()
}

// export
func ServeWs(w http.ResponseWriter, r *http.Request) {
	// Uso_online += 1
	fmt.Println("new usor addin")
	if hall.IsRunning == false {
		fmt.Println("ServeWs Run Hall Run")
		Run()
	}
	conn, err := Uso_websocket_upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("ServeWs upgrade error")
		return
	}
	uso_conn := Usoconn{
		Name:   "",
		Quit:   false,
		R:      make(chan Message),
		W:      make(chan Message),
		wsconn: conn,
	}
	Uso_connpool = append(Uso_connpool, uso_conn)
	uso_conn.Run()
	hall.ServeGuest(&uso_conn)
}
