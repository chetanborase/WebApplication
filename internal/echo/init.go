package echo

import (
	_ "WebApplication/docs"
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func Init() {
	e := echo.New()
	e.Static("/", "public")

	bulkUpload := e.Group("batch/")
	bulkUpload.GET("accounts/:id", AccountsBatchUpload)

	g := e.Group("/accounts/")
	g.GET(":id", ShowAccount)
	g.POST("", CreateAccount)
	g.PUT("", UpdateAccount)
	g.DELETE(":id ", DeleteAccount)
	g.GET("search/", ListAccounts)
	e.GET("swagger/*", echoSwagger.WrapHandler)
	//go startSwagger()
	e.Logger.Fatal(e.Start(echoAddress()))
}
