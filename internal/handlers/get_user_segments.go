package handlers

import (
	"github.com/edos10/test_avito_service/internal/databases"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetUserSegments(c *gin.Context) {
	var requestData struct {
		UserID int `json:"user_id"`
	}
	if err := c.BindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}
	db, errDb := databases.CreateDatabaseConnect()
	if errDb != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "couldn't open database"})
		return
	}
	defer db.Close()

	query := `
		SELECT ids.segment_name
		FROM users_segments us
		INNER JOIN id_name_segments ids ON us.segment_id = ids.segment_id
		WHERE us.user_id = 1
	`
	data, errGet := db.Query(query, requestData.UserID)
	if errGet != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to parse user segments"})
		return
	}
	defer data.Close()

	var segments []string
	for data.Next() {
		var segmentName string
		if err := data.Scan(&segmentName); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}
		segments = append(segments, segmentName)
	}

	if err := data.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error in query process"})
		return
	}
	c.JSON(http.StatusOK, segments)
}
