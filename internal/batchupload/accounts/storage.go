package accounts

import (
	"WebApplication/config"
	"WebApplication/internal/logger"
	"WebApplication/internal/socket"
	"WebApplication/internal/storage/database"
	"WebApplication/internal/storage/database/model"
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/jmoiron/sqlx"
	"sync"
	"time"
)

type MetaData struct {
	File     *excelize.File
	FileName string
	Socket   *socket.Wrapper
	Cancel   *bool
	dbmd     *DbMetadata
	status   Status
	dberr    error
}
type Status struct {
	TotalRows int
	UnitRows  int
	StartTime time.Time
	EndTime   time.Time
	Stamp     time.Time
	xlserr    bool
}
type DbMetadata struct {
	Tx     *sqlx.Tx
	Stmt1  *sqlx.Stmt
	Stmt2  *sqlx.Stmt
	Stmt3  *sqlx.Stmt
	DoOnce *sync.Once
}

func (md *DbMetadata) StoreRecords(user model.User) (err error) {
	md.DoOnce.Do(func() {
		err = md.PrepareUpload()
		if err != nil {
			return
		}
	})
	err = database.StoreUser(md.Stmt1, user)
	if err != nil {
		logger.Println(err)
		return err
	}
	err = database.StoreAddress(md.Stmt2, user.Address)
	if err != nil {
		logger.Println(err)
		return err
	}
	err = database.StoreContact(md.Stmt3, user.Contact)
	return err
}
func (md *DbMetadata) PrepareUpload() (err error) {
	md.Tx, err = config.GetDB().Beginx()
	if err != nil {
		logger.Println(err)
		return
	}
	md.Stmt1, err = database.PrepareUser(md.Tx)
	if err != nil {
		logger.Println(err)
		return
	}
	md.Stmt2, err = database.PrepareAddress(md.Tx)
	if err != nil {
		logger.Println(err)
		return
	}
	md.Stmt3, err = database.PrepareContact(md.Tx)
	if err != nil {
		logger.Println(err)
		return
	}
	return nil
}
