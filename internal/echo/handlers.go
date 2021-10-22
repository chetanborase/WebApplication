package echo

import (
	"WebApplication/internal/batchupload/accounts"
	"WebApplication/internal/socket"
	"WebApplication/internal/storage"
	"WebApplication/internal/storage/model"
	"github.com/labstack/echo/v4"
	"net/http"
	"strings"
	"sync"
)

// CreateAccount godoc
// @Summary Create New Account
// @Description Accepts User Data in json format
// @Accept  json
// @Produce  json
// @Param account body model.User true "User Data"
// @Success 200 {string} string "successfully created"
// @Failure 400,404 {string} string	"error"
// @failure 422 {string} string	"error"
// @Failure 500 {string} string	"error"
// @response default {string} string "http response"
// @Router /accounts/ [post]
func CreateAccount(ctx echo.Context) error {
	u := model.User{}
	err := ctx.Bind(&u)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}
	err = storage.InsertUser(u)
	if err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") {
			return ctx.JSON(http.StatusUnprocessableEntity, err.Error())
		} else {
			return ctx.JSON(http.StatusInternalServerError, err.Error())
		}
	}
	return ctx.JSON(http.StatusOK, "Successfully Created")
}

// UpdateAccount godoc
// @Summary Update User Details
// @Description Accept User Details in json format and updates details in database
// @Accept  json
// @Produce  json
// @Param account body model.User true "User Data"
// @Success 200 {string} string "Updated Successfully"
// @Failure 400,404 {string} string	"error"
// @failure 422 {string} string	"User Not Exists"
// @Failure 500 {string} string	"error"
// @response default {string} string "http response"
// @Router /accounts/ [put]
func UpdateAccount(ctx echo.Context) error {
	u := model.User{}
	err := ctx.Bind(&u)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}
	err = storage.UpdateUser(u)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return ctx.JSON(http.StatusUnprocessableEntity, "user not exist")
		} else {
			return ctx.JSON(http.StatusInternalServerError, err.Error())
		}
	}
	return ctx.JSON(http.StatusOK, "Updated Successfully")
}

// DeleteAccount godoc
// @Summary Delete User Account
// @Description Accept User uuid as a query parameter and deletes account
// @Accept  json
// @Produce  json
// @Param id path string true "User UUID"
// @Success 200 {string} string "Deleted Successfully"
// @Failure 400,404 {string} string	"error"
// @failure 422 {string} string	"user not exist"
// @Failure 500 {string} string	"error"
// @response default {string} string "http response"
// @Router /accounts/{id} [delete]
func DeleteAccount(ctx echo.Context) error {
	uuid := ctx.Param("id")
	err := storage.DeleteUser(uuid)
	if err != nil {
		if strings.Contains(err.Error(), "user not exist") {
			return ctx.JSON(http.StatusUnprocessableEntity, err.Error())
		} else {
			return ctx.JSON(http.StatusInternalServerError, err.Error())
		}
	}
	return ctx.JSON(http.StatusOK, "Deleted Successfully")
}

// ListAccounts godoc
// @Summary Get Bunch of Users details
// @Description Apply search and filter on database records and fetch result
// @Accept  plain
// @Produce  json
// @Param key query string false "key"
// @Param city query []string false "city" collectionFormat(multi)
// @Success 200 dbModel.SearchResult "Search Result"
// @Failure 400,404 {string} string	"error"
// @Failure 500 {string} string	"error"
// @response default {string} string "http response"
// @Router /accounts/search/ [get]
// format http://localhost:8080/accounts/search/?key=ab&city=sanglii
func ListAccounts(ctx echo.Context) error {
	//ipAddr := ctx.RealIP()
	//logger.Println(ipAddr)
	//log.Fatal("")
	urlValue := ctx.QueryParams()
	//fmt.Println(urlValue["key"])
	//fmt.Println(urlValue["city"])
	//return nil
	searchResult, err := storage.GetUserSet(urlValue)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}
	return ctx.JSONPretty(http.StatusOK, searchResult, "\n")
}

// ShowAccount godoc
// @Summary Show Account Details
// @Description Accepts User Data in json format
// @Accept  plain
// @Produce  json
// @Param id path string true "User UUID"
// @Success 200 {array} model.User ""
// @Failure 400,404 {string} string "User Not Found"
// @Failure 500 {string} string	"error"
// @Router /accounts/{id} [get]
func ShowAccount(ctx echo.Context) error {
	uuid := ctx.Param("id")
	user, err := storage.GetUser(uuid)
	//err= errors.New("no rows in result set")
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return ctx.JSON(http.StatusNotFound, "User Not Found")
		} else {
			return ctx.JSON(http.StatusInternalServerError, err.Error())
		}
	}
	return ctx.JSONPretty(http.StatusOK, user, "\n")
}

// AccountsBatchUpload godoc
// @Summary Upload Multiple User Accounts
// @Description Accept excel file in WebSocket than Read Validates and Stores To Database
// @tag.name  batchupload
// @Accept  json
// @Produce  json
// @Param id query string true "User UUID"
// @Success 200 {string} string	"ok"
// @Failure 400,404 {string} string	"error"
// @failure 422 {string} string	"error"
// @Failure 500 {string} string	"error"
// @response default {string} string "http response"
// @Router /batch/accounts/{id} [get]
func AccountsBatchUpload(ctx echo.Context) error {
	uuid := ctx.Param("id")
	userKey := socket.KeyModel{Uuid: uuid}
	wrap, err, isNewSocket := socket.GetSocket(userKey, socket.Accounts, ctx.Request(), ctx.Response())
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}
	go wrap.Read()
	if !isNewSocket {
		return nil
	}
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go accounts.HandleUpload(wrap, wg)
	defer socket.Delete(userKey, socket.Accounts)
	wg.Wait()
	return nil
}
