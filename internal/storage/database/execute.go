package database

import (
	"WebApplication/config"
	"WebApplication/internal/common"
	"WebApplication/internal/logger"
	"WebApplication/internal/storage/database/model"
	model2 "WebApplication/internal/storage/model"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"log"
	"time"
)

func PrepareUser(tx *sqlx.Tx) (stmt *sqlx.Stmt, err error) {
	return tx.Preparex(insertUser)
}
func PrepareAddress(tx *sqlx.Tx) (stmt *sqlx.Stmt, err error) {
	return tx.Preparex(insertAddress)
}
func PrepareContact(tx *sqlx.Tx) (stmt *sqlx.Stmt, err error) {
	return tx.Preparex(insertContact)
}
func StoreUser(stmt *sqlx.Stmt, user model.User) error {
	if stmt == nil {
		panic("Nil Pointer")
	}

	_, err := stmt.Exec(user.UserUUID, user.IdXref,
		user.FirstName, user.MiddleName,
		user.LastName, user.FullName,
		user.Email, user.UserName, user.Password,
		time.Now())
	if err != nil {
		logger.Error(err)
	}
	return err
}
func StoreAddress(stmt *sqlx.Stmt, addr model.Address) error {
	_, err := stmt.Exec(*addr.AddressUserUUID, *addr.Address1,
		*addr.Address2, *addr.Address3,
		*addr.Area, *addr.City, *addr.State,
		*addr.Country, *addr.PinCode)
	if err != nil {
		logger.Println(err)
	}
	return err
}
func StoreContact(stmt *sqlx.Stmt, contact model.Contact) error {
	_, err := stmt.Exec(*contact.ContactUserUUID, *contact.SocialMediaID,
		*contact.WebSite, *contact.DialCode,
		*contact.PhoneNumber, *contact.FullPhoneNumber)
	if err != nil {
		logger.Println(err)
	}
	return err
}
func UpdateUser(u model.User) (err error) {
	tx, err := config.GetDB().Beginx()
	if err != nil {
		return err
	}
	defer func() { _ = common.FinalizeTransaction(tx, err) }()
	stmt1, _ := tx.Preparex(updateUser)
	stmt2, _ := tx.Preparex(updateAddress)
	stmt3, _ := tx.Preparex(updateContact)
	defer func() { stmt1.Close(); stmt2.Close(); stmt3.Close() }()
	_, err = tx.Exec(lockTable)
	defer func() { _, _ = tx.Exec(unlockTable) }()
	if err != nil {
		return err
	}
	result, err := stmt1.Exec(u.IdXref, u.FirstName, u.MiddleName, u.LastName, u.FullName, u.Email, u.Password, time.Now().UTC(), u.UserUUID)

	if err != nil {
		return err
	}
	if result != nil {
		var i int64
		i, err = result.RowsAffected()
		if i == 0 {
			err = errors.New("records not found")
			return err
		}
	}
	_, err = stmt2.Exec(u.Address1, u.Address2, u.Address3, u.Area, u.City, u.State, u.Country, u.PinCode, u.UserUUID)
	if err != nil {
		return err
	}
	_, err = stmt3.Exec(u.SocialMediaID, u.WebSite,
		u.DialCode,
		u.PhoneNumber,
		u.FullPhoneNumber, u.UserUUID)

	if err != nil {
		log.Fatal("failed to unlock", err)
	}
	return err
}
func DeleteUser(uuid string) error {
	result, err := config.GetDB().Exec(deleteUser, time.Now(), uuid)
	if err != nil {
		logger.Println(err)
		return err
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("no account found with given user id")
	}
	return nil
}
func getResult(query string, args ...interface{}) ([]model.User, error) {
	//fmt.Print("Query  =  ", query, "  args   : ", args)
	row, err := config.GetDB().Queryx(query, args...)
	if err != nil {
		logger.Println(err)
		return nil, err
	}
	var users []model.User
	for row.Next() {
		u := model.User{}
		err = row.Scan(&u.Id, &u.UserUUID, &u.IdXref, &u.FirstName, &u.MiddleName, &u.LastName,
			&u.FullName, &u.Email, &u.UserName, &u.CreatedAt, &u.Address1, &u.Address2, &u.Address3,
			&u.Area, &u.City, &u.State, &u.Country, &u.PinCode,
			&u.SocialMediaID, &u.WebSite, &u.DialCode, &u.PhoneNumber, &u.FullPhoneNumber)
		if err != nil {
			logger.Println(err)
		}
		users = append(users, u)
	}
	return users, err
}
func GetUserRecords(values model2.UrlRequest) ([]model.User, error) {
	query, args := getSearchQuery(values)
	fmt.Println(query, args)
	//return nil, nil
	users, err := getResult(query, args...)

	if err != nil {
		return nil, err
	}
	return users, err
}
func GetUser(uuid string) (model.User, error) {
	result := config.GetDB().QueryRow(`select u.id,u.UserUUID,u.IdXref,u.FirstName,u.MiddleName,u.LastName,
       u.FullName,u.Email,u.UserName,u.CreatedAt,AI.Address1,AI.Address2,AI.Address3,
       AI.Area,AI.City,AI.State,AI.Country,AI.PinCode,
       CI.SocialMediaID,CI.WebSite,CI.DialCode,CI.PhoneNumber,CI.FullPhoneNumber
from Tbl_UserInfo as u
         inner join Tbl_AddressInfo AI
                    on u.UserUUID = AI.UserUUID
         inner join Tbl_ContactInfo CI on u.UserUUID = CI.UserUUID
where u.UserUUID = ? and DeletedAt is null`, uuid)
	u := model.User{}
	err := result.Scan(&u.Id, &u.UserUUID, &u.IdXref, &u.FirstName, &u.MiddleName, &u.LastName,
		&u.FullName, &u.Email, &u.UserName, &u.CreatedAt, &u.Address1, &u.Address2, &u.Address3,
		&u.Area, &u.City, &u.State, &u.Country, &u.PinCode,
		&u.SocialMediaID, &u.WebSite, &u.DialCode, &u.PhoneNumber, &u.FullPhoneNumber)
	return u, err
}

func InsertUser(u model.User) error {
	if u.UserUUID == nil {
		return errors.New("non-empty user uuid required")
	}
	tx, err := config.GetDB().Beginx()
	defer func() { _ = common.FinalizeTransaction(tx, err) }()
	if err != nil {
		return err
	}
	stmt1, err := PrepareUser(tx)
	if err != nil {
		return err
	}
	stmt2, err := PrepareContact(tx)
	if err != nil {
		return err
	}
	stmt3, err := PrepareAddress(tx)
	if err != nil {
		return err
	}
	err = StoreUser(stmt1, u)
	if err != nil {
		return err
	}
	err = StoreAddress(stmt3, u.Address)
	if err != nil {
		return err
	}
	err = StoreContact(stmt2, u.Contact)
	if err != nil {
		return err
	}
	return err
}

func LockUserAccounts() (err error) {
	_, err = config.GetDB().Exec(lockTable)
	return err
}
func UnLockUserAccounts() {
	_, err := config.GetDB().Exec(unlockTable)
	if err != nil {
		logger.Error(err)
	}
}
