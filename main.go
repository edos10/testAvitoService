package main

import (
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	log.Print("Сервер успешно запущен!")
	r := gin.Default()
	r.GET("/ping")
	log.Fatal(r.Run())
}
