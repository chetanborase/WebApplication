package storage

import (
	"WebApplication/internal/common"
	"WebApplication/internal/logger"
	"WebApplication/internal/storage/cache"
	"WebApplication/internal/storage/database"
	dbModel "WebApplication/internal/storage/database/model"
	"WebApplication/internal/storage/model"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"log"
	"net/url"
	"time"
)

func GetUserSet(values url.Values) (searchResult dbModel.SearchResult, err error) {
	before := time.Now()
	reqValue := model.MapUrlToStruct(values)
	defer func() { searchResult.Time = fmt.Sprint(time.Since(before)) }()
	key, err := common.CreateSimpleHash(reqValue)
	if err != nil {
		return dbModel.SearchResult{}, err
	}
	value, found := cache.Get(key, cache.DB2)
	if found {
		if searchResult, ok := value.(dbModel.SearchResult); ok {
			return searchResult, nil
		} else {
			logger.Log.Println("failed to unbind")
		}
	}

	users, err := database.GetUserRecords(reqValue)

	if err != nil {
		logger.Error(err)
		return searchResult, err
	}
	searchResult = searchResponse(reqValue, users)
	//saveCache(key string,value inter)
	cache.Set(key, searchResult, cache.DB2)
	return searchResult, nil
}
func searchResponse(values model.UrlRequest, users []dbModel.User) dbModel.SearchResult {

	isLastPage := false
	rowsCount, _ := database.GetTotalRecords()
	totalPages := rowsCount / values.Limit
	if values.Page >= totalPages {
		isLastPage = true
	}
	length := len(users)
	maxId := 0
	if length >= 1 {
		maxId = *users[length-1].Id
	}
	return dbModel.SearchResult{
		Users:      users,
		PageNo:     values.Page,
		TotalPages: totalPages,
		IsLastPage: isLastPage,
		LastID:     maxId,
	}
}

func GetUser(key string) (dbModel.User, error) {
	value, found := cache.Get(key, cache.DB1)
	if found && value != nil {
		logger.Log.Println("cache hit")
		u := value.(dbModel.User)
		return u, nil
	}
	user, err := database.GetUser(key)
	if err != nil {
		logger.Println(err)
	}
	cache.Set(key, user, cache.DB1)
	return user, err
}

func InsertUser(user model.User) error {
	u := dbModel.MapToDbModel(user)
	id := uuid.NewString()
	u.UserUUID = &id
	u.AddressUserUUID = &id
	u.ContactUserUUID = &id
	err := database.InsertUser(u)
	if err != nil {
		return err
	}
	cache.Set(user.UserUUID, u, cache.DB1)
	return nil
}

func UpdateUser(user model.User) error {
	u := dbModel.MapToDbModel(user)
	//val, found := cache.Get(*u.UserUUID, cache.DB1)
	//if found {
	//	userFromCache := val.(database.User)
	//	if u.UpdatedAt.Before(userFromCache.UpdatedAt) {
	//		return errors.New("data has been changed from somewhere else, Please Request Again")
	//	}
	//} else {
	/*
		err := database.LockUserAccounts()
		if err != nil {
			log.Fatal(err)
		}*/
	if u.UserUUID==nil {
		return errors.New("user uuid is null")
	}
	dbUser, err := database.GetUser(*u.UserUUID)
	if err != nil {
		log.Println(err)
		return err
	}
	if !dbUser.UpdatedAt.Equal(u.UpdatedAt) {
		return errors.New("data has been changed from somewhere else, Please Request Again")
	}
	/*
		defer database.UnLockUserAccounts()
	*/
	err = database.UpdateUser(u)
	if err != nil {
		return err
	}
	cache.Set(user.UserUUID, u, cache.DB1)
	return nil
}

func DeleteUser(uuid string) error {
	err := database.DeleteUser(uuid)
	if err != nil {
		return err
	}
	cache.Delete(uuid, cache.DB1)
	return nil
}
