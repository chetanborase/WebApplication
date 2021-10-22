package main

import (
	"WebApplication/env"
	"WebApplication/internal/echo"
	"WebApplication/internal/storage/cache"
)

func init() {
	env.LoadProfile()
}

// @title Assignment API
// @version 1.0
// @description This is Simple Web Application.
// @schemes http
// @BasePath /
// @host localhost:8080
func main() {
	//return
	//fmt.Println(env.Get(env.DbDriveName,"not found"))
	cache.Init()
	echo.Init()
}
