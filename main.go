package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// 新しい構造体に対応
type Message struct {
	DeviceID   string  `json:"device_id"`
	Detected   bool    `json:"detected"`
	Confidence float64 `json:"confidence"`
	Timestamp  int64   `json:"timestamp"`
}

var latest Message

func postHandler(c echo.Context) error {
	var m Message
	if err := c.Bind(&m); err != nil {
		return c.String(http.StatusBadRequest, "Invalid JSON")
	}

	// 保存
	latest = m

	// JST時刻へ変換
	jst := time.FixedZone("Asia/Tokyo", 9*60*60)
	t := time.Unix(m.Timestamp, 0).In(jst)
	tStr := t.Format("2006-01-02 15:04:05")

	// 1. ターミナルにログ出力
	fmt.Printf("[受信] device: %s, detected: %v, confidence: %.2f, at %s\n",
		m.DeviceID, m.Detected, m.Confidence, tStr)

	// 2. ファイルにログ保存
	f, err := os.OpenFile("log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err == nil {
		defer f.Close()
		logLine := fmt.Sprintf("[受信] device: %s, detected: %v, confidence: %.2f, at %s\n",
			m.DeviceID, m.Detected, m.Confidence, tStr)
		f.WriteString(logLine)
	}

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
	log.Println("Server started on :8080")
	e.Logger.Fatal(e.Start(":8080"))
}

// package main

// import (
// 	"net/http"
// 	"time"

// 	"github.com/labstack/echo/v4"
// 	"github.com/labstack/echo/v4/middleware"
// )

// type Message struct {
// 	Message   string `json:"message"`
// 	Timestamp int64  `json:"timestamp"`
// }

// var latest Message

// func postHandler(c echo.Context) error {
// 	var m Message
// 	if err := c.Bind(&m); err != nil {
// 		return c.String(http.StatusBadRequest, "Invalid")
// 	}
// 	m.Timestamp = time.Now().Unix()
// 	latest = m
// 	return c.JSON(http.StatusOK, map[string]string{"status": "received"})
// }

// func getHandler(c echo.Context) error {
// 	return c.JSON(http.StatusOK, latest)
// }

// func main() {
// 	e := echo.New()
// 	e.Use(middleware.CORS())
// 	e.POST("/post", postHandler)
// 	e.GET("/latest", getHandler)
// 	e.Logger.Fatal(e.Start(":8080"))
// }
