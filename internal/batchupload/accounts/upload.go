package accounts

import (
	"WebApplication/config"
	"WebApplication/internal/common"
	"WebApplication/internal/logger"
	"WebApplication/internal/socket"
	"WebApplication/internal/storage/database/validate"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

var errCellName = "W"
var temp time.Time

func (md *MetaData) StartUpload() {
	///assume there is only one sheet
	sheets := md.File.GetSheetMap()
	_ = md.Socket.Conn.WriteMessage(websocket.TextMessage, []byte("generating metadata..."))
	///loop for every sheets in given file
	for _, sheetName := range sheets {
		err := md.uploadSheet(sheetName)
		if err != nil {
			logger.Println(err)
		}
	}
}

func (md *MetaData) uploadSheet(sheetName string) (err error) {
	md.dbmd = &DbMetadata{DoOnce: &sync.Once{}}
	_ = md.Socket.Conn.WriteMessage(websocket.TextMessage, []byte("validating file..."))

	rows := md.File.GetRows(sheetName)
	GenUploadMetadata(md, rows)
	if md.status.TotalRows <= 1 {
		return errors.New("empty sheet")
	}
	err = validate.SheetHeader(rows[0])
	if err != nil {
		_ = md.Socket.Conn.WriteMessage(websocket.TextMessage, []byte(err.Error()))
		logger.Println()
		return err
	}
	_ = md.Socket.Conn.WriteMessage(websocket.TextMessage, []byte("Hang on we are preparing database connection"))
	md.status.StartTime = time.Now()
	for i := 1; i < md.status.TotalRows; i++ {
		///check if cancelled by user
		if *md.Cancel {
			_ = md.Socket.Conn.WriteMessage(websocket.TextMessage, []byte("Cancelled Successfully"))
			///init has error to rollback
			md.dberr = errors.New("cancelled by user")
			return nil
		}
		///Conn.WriteMessage Status to User
		if i%md.status.UnitRows == 0 {
			md.status.SendUploadStatus(i, md.Socket)
		}
		///Validate record and populate to struct
		u, errEncountered := validate.SheetData(rows[i])
		if errEncountered == nil {
			if md.dberr == nil {
				md.dberr = md.dbmd.StoreRecords(u)
			}
			continue
		}
		if !md.status.xlserr {
			md.status.xlserr = true
		}
		///Write Error to xlsx, if Encountered while SheetData validating
		md.File.SetCellValue(sheetName, errCellName+strconv.Itoa(i+1), errEncountered.Error()) //index out of index
	}
	md.status.EndTime = time.Now()
	_ = md.Socket.Conn.WriteMessage(websocket.TextMessage, []byte("Upload Completed "+md.status.EndTime.Sub(md.status.StartTime).String()))
	defer func() {
		err = finalizeUpload(md)
	}()
	return err
}

func finalizeUpload(md *MetaData) error {
	///Conn.WriteMessage completion response to client
	err := common.FinalizeTransaction(md.dbmd.Tx, md.dberr)
	if err != nil {
		logger.Error(err)
	}
	err = md.dbmd.Stmt1.Close()
	if err != nil {
		logger.Error(err)
	}
	err = md.dbmd.Stmt2.Close()
	if err != nil {
		logger.Error(err)
	}
	err = md.dbmd.Stmt3.Close()
	if err != nil {
		logger.Error(err)
	}
	if md.status.xlserr == true {
		resDir, resFile := filepath.Split(md.FileName)
		resFile = strings.ReplaceAll(resFile, "_req", "_res")
		resfilepath := filepath.Join(resDir, resFile)
		err := md.File.SaveAs(resfilepath)
		if err != nil {
			logger.Println(err)
		}
		file, err := os.Open(resfilepath)
		if err != nil {
			logger.Println(err)
		}
		bytes, err := ioutil.ReadAll(file)
		err = md.Socket.Conn.WriteMessage(websocket.BinaryMessage, bytes)
		if err != nil {
			logger.Println(err)
		}
	}
	StoreActivityToDb(md)
	return nil
}
func (upst *Status) SendUploadStatus(rowindex int, socket *socket.Wrapper) {
	temp = time.Now()
	estimateTime := int(time.Since(upst.Stamp).Seconds()) * (upst.TotalRows - rowindex) / upst.UnitRows
	t, _ := time.ParseDuration(fmt.Sprintf("%v%s", estimateTime, "s"))
	status := fmt.Sprint("uploading Row No. ", rowindex, " Completed ", (rowindex*100)/upst.TotalRows, "% , Estimate Completion Time =", t)
	upst.Stamp = time.Now()
	fmt.Println(time.Since(temp))
	_ = socket.Conn.WriteMessage(websocket.TextMessage, []byte(status))
}
func GenUploadMetadata(md *MetaData, rows [][]string) {
	md.status = Status{}
	md.status.TotalRows = len(rows)
	md.status.UnitRows = md.status.TotalRows / 100
}
func StoreActivityToDb(md *MetaData) {
	insertActivity := fmt.Sprint(`insert into Tbl_UploadHistory ( 
                               FileName, StartTime, 
                               EndTime, TotalTime) 
                               value (?,?,?,?)`)
	totalTime := md.status.EndTime.Sub(md.status.StartTime).Seconds()
	_, err := config.GetDB().Exec(insertActivity,
		md.FileName, md.status.StartTime,
		md.status.EndTime, totalTime)
	if err != nil {
		logger.Println(err)
	}
}
