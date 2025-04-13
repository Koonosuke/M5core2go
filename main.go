package main

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Message struct {
	Message   string `json:"message"`
	Timestamp int64  `json:"timestamp"`
}

var latest Message

func postHandler(c echo.Context) error {
	var m Message
	if err := c.Bind(&m); err != nil {
		return c.String(http.StatusBadRequest, "Invalid")
	}
	m.Timestamp = time.Now().Unix()
	latest = m
	return c.JSON(http.StatusOK, map[string]string{"status": "received"})
}

func getHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, latest)
}

func main() {
	e := echo.New()
	e.Use(middleware.CORS())
	e.POST("/post", postHandler)
	e.GET("/latest", getHandler)
	e.Logger.Fatal(e.Start(":8080"))
}
