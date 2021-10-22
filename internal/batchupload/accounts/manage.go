package accounts

import (
	"WebApplication/internal/logger"
	"WebApplication/internal/socket"
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/gorilla/websocket"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"
)

func HandleUpload(wrap *socket.Wrapper, wg *sync.WaitGroup) {
	defer wg.Done()
	var cancelled bool
	md := MetaData{
		Socket: wrap,
		Cancel: &cancelled,
	}
	for {
		data := <-wrap.Data
		if string(data.Bt) == "cancel" {
			cancelled = true
			break
		}
		if string(data.Bt) == "close" {
			break
		}
		if data.MsgType == websocket.BinaryMessage {
			err := md.Socket.Conn.WriteMessage(websocket.TextMessage, []byte("Processing file, please wait!"))
			fileName := getFileName()
			err = ioutil.WriteFile(fileName, data.Bt, os.ModePerm)
			if err != nil {
				logger.Println(err)
			}
			file, err := excelize.OpenFile(fileName)
			if err != nil {
				logger.Println(err)
			}
			md.File = file
			go md.StartUpload()
		}
	}
}

func getFileName() (filename string) {
	fileName := strconv.Itoa(int(time.Now().UnixNano())) + "_req"
	filename = filepath.Join("public", fileName+".xlsx")
	return
}
