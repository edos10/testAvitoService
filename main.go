package main

import (
	"fmt"
	"github.com/edos10/test_avito_service/internal/handlers"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
)

func main() {
	_ = godotenv.Load(".env")
	gin.SetMode(gin.ReleaseMode)
	log.Print("Сервер успешно запущен!")
	r := gin.Default()
	r.POST("/create_segment", handlers.CreateSegment)
	r.DELETE("/delete_segment", handlers.DeleteSegment)
	r.PUT("/change_segments", handlers.ChangesUserSegments)
	r.GET("/get_user_segments", handlers.GetUserSegments)
	r.GET("/get_report_csv", handlers.GenerateCSV)
	serverPort := 5050
	if port := os.Getenv("GO_DOCKER_PORT"); port != "" {
		numPort, err := strconv.Atoi(port)
		if err == nil {
			serverPort = numPort
		}
	}
	log.Fatal(r.Run(fmt.Sprintf("0.0.0.0:%d", serverPort)))
}
