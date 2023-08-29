package handlers

import (
	"database/sql"
	"github.com/edos10/test_avito_service/internal/databases"
	"github.com/gin-gonic/gin"
	"net/http"
)

func CreateSegment(c *gin.Context) {
	var requestData struct {
		SegmentName string `json:"segment_name"`
	}
	db, errDb := databases.CreateDatabaseConnect()
	if errDb != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "couldn't open the database"})
		return
	}
	defer db.Close()
	if err := c.BindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "invalid data, try please again"})

		return
	}
	var existID int64
	queryCheck := "SELECT segment_id FROM id_name_segments WHERE segment_name = $1"
	queryInsert := "INSERT INTO id_name_segments (segment_name) VALUES ($1) RETURNING segment_id"
	errCheck := db.QueryRow(queryCheck, requestData.SegmentName).Scan(&existID)
	if errCheck == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "segment already exists"})
		return
	} else if errCheck != sql.ErrNoRows {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error in query process"})
		return
	}
	var newSegmentID int
	errIns := db.QueryRow(queryInsert, requestData.SegmentName).Scan(&newSegmentID)
	if errIns != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to insert segment"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "segment is successfully created"})
}
