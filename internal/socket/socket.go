package socket

import (
	"WebApplication/internal/logger"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
)

type (
	KeyModel struct {
		Uuid string
	}
	Wrapper struct {
		Conn *websocket.Conn
		Data chan Data
	}
	Trans = map[int]*Wrapper
)

type Data struct {
	MsgType int
	Bt      []byte
	Err     error
}

var (
	Map      = map[KeyModel]Trans{}
	upgrader = websocket.Upgrader{
		HandshakeTimeout:  0,
		ReadBufferSize:    1024,
		WriteBufferSize:   1024,
		WriteBufferPool:   nil,
		Subprotocols:      nil,
		Error:             nil,
		CheckOrigin:       nil,
		EnableCompression: false,
	}
)

func GetSocket(model KeyModel, kind int, req *http.Request, res *echo.Response) (*Wrapper, error, bool) {
	conn, err := upgrader.Upgrade(res.Writer, req, nil)
	if err != nil {
		logger.Error(err)
		return nil, err, false
	}
	ts, ok := Map[model]
	if !ok { //means user want to initiate first transaction
		//create transaction map for user
		ts = make(Trans)
		wrap := &Wrapper{Conn: conn, Data: make(chan Data)}
		ts[kind] = wrap
		Map[model] = ts
		return wrap, nil, true
	}
	wrap, ok := ts[kind]
	if !ok {
		wrap = &Wrapper{
			Conn: conn,
			Data: make(chan Data),
		}
		ts[kind] = wrap
		return wrap, err, true
	}
	return wrap, nil, false
}

func (sw *Wrapper) Read() {
	for {
		msgType, bt, err := sw.Conn.ReadMessage()
		if err != nil {
			break
		}
		data := Data{
			MsgType: msgType,
			Bt:      bt,
			Err:     err,
		}
		sw.Data <- data
	}
	log.Println("exiting from read function...")
}

func Delete(userKey KeyModel, kind int) {
	trans, ok := Map[userKey]
	if !ok {
		return
	}
	wrap, ok := trans[kind]
	if !ok {
		return
	}
	_ = wrap.Conn.Close()
	delete(trans, kind)
	if len(trans) == 0 {
		delete(Map, userKey)
	}
}

//func PrintSockets(group *sync.WaitGroup){
//	defer group.Done()
//	for {
//		for _, v := range Map {
//			for _, vv := range v {
//				log.Println(v, "  and ", vv)
//				time.Sleep(5*time.Second)
//			}
//			log.Println(v, "  and  outside ")
//			time.Sleep(3*time.Second)
//		}
//	}
//}
