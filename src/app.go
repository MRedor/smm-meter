package main

import (
	"app"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func StartServer() {
	e := echo.New()
	e.Use(middleware.CORSWithConfig(middleware.DefaultCORSConfig))


	e.GET("/api/channel/:id", app.GetChannel)
	e.GET("api/video/:id", app.GetVideo)
	e.GET("/api/timeline/:id", app.GetTimeline)
	e.POST("/api/timeline/start", app.StartWatching)

	e.Logger.Fatal(e.Start(":1323"))
}

func main() {
	StartServer()
}
