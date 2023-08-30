package main

import (
	"github.com/edos10/test_avito_service/internal/handlers"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
)

func main() {
	_ = godotenv.Load(".env")
	gin.SetMode(gin.ReleaseMode)
	log.Print("Сервер успешно запущен!")
	r := gin.Default()
	r.POST("/create_segment", handlers.CreateSegment)
	r.DELETE("/delete_segment", handlers.DeleteSegment)
	r.POST("/change_segments", handlers.ChangesUserSegments)
	r.GET("/get_user_segments", handlers.GetUserSegments)
	r.GET("/csv", handlers.GenerateCSV)
	log.Fatal(r.Run())
}
