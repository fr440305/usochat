package uso // import "github.com/fr440305/uso"

import "github.com/gorilla/websocket"
import "net/http"
import "fmt"

const (
	MsgType_uso_login = "ul"
	MsgType_uso_err   = "ue"
	MsgType_uso_quit  = "uq"
	// ...
)

var Uso_connpool = []Usoconn{}
var Uso_Roompool = []Room{}
var Uso_websocket_upgrader = &websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Message []string

// string -> []interface{} // tokenlize
// []interface{} -> UsoMsg // done
func NewMessage(bar []byte) Message {
	return Message{string(bar)} //dummy
}

func (um Message) Type() string {
	if len([]string(um)) == 0 {
		return ""
	} else {
		return um[0]
	}
}

func (um Message) Cont() []string {
	return []string{}
}

// return []byte(`["Type", "Cont-1", "Cont-2", ... , "Cont-n"]`)
func (um Message) ToJbar() []byte {
	if len([]string(um)) == 0 {
		return []byte{}
	} else {
		return []byte(um[0])
	}
}

type Usoconn struct {
	//if you are not in a Room, you do not need a name
	Name   string
	R      chan Message
	W      chan Message
	wsconn *websocket.Conn
}

// public.
func (uc Usoconn) Run() {
	// if err in wsconn then close R and close W.
	go func(rch chan Message) {
		for {
			_, cont, err := uc.wsconn.ReadMessage()
			if err != nil {
				fmt.Println("uso conn exit.")
				rch <- Message{MsgType_uso_quit}
				close(rch)
				return
			} else {
				fmt.Println(string(cont))
				rch <- NewMessage(cont)
			}
		}
	}(uc.R)
	go func(wch chan Message) {
		for {
			msg := <-wch
			err := uc.wsconn.WriteMessage(
				websocket.TextMessage,
				msg.ToJbar(),
			)
			if err != nil {
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
	h.IsRunning = true
	for msg := range h.jobs {
		fmt.Println("Hall Run", msg)
	}
}

func (h Hall) ServeGuest(uc *Usoconn) {
	//read uc.R
	for msg := range uc.R {
		fmt.Println("Hall ServeGuest", msg)
		h.jobs <- msg
	}
	// uc.R closed:
	h.jobs <- Message{MsgType_uso_quit}
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
		R:      make(chan Message),
		W:      make(chan Message),
		wsconn: conn,
	}
	Uso_connpool = append(Uso_connpool, uso_conn)
	uso_conn.Run()
	hall.ServeGuest(&uso_conn)
}
