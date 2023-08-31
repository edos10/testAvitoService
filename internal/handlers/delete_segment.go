package handlers

import (
	"database/sql"
	"github.com/edos10/test_avito_service/internal/databases"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func DeleteSegment(c *gin.Context) {
	var requestData struct {
		SegmentName string `json:"segment_name"`
		UserID      int    `json:"user_id"`
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
	checkQuery := "SELECT segment_id FROM id_name_segments WHERE segment_name = $1"
	var segmentID int
	if err := db.QueryRow(checkQuery, requestData.SegmentName).Scan(&segmentID); err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "segment not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error in query process"})
		}
		return
	}
	// проверочка на то, что сегмент не был добавлен автоматически некоторым пользователям
	queryForCheckAutomaticSegment := `
		SELECT ush.user_id
		FROM user_segment_history ush
		JOIN id_name_segments ins ON ush.segment_id = ins.segment_id
		WHERE ush.timestamp = 'infinity'
	`
	var Id int64
	errCheck := db.QueryRow(queryForCheckAutomaticSegment).Scan(&Id)
	if errCheck == nil {
		c.JSON(http.StatusConflict, gin.H{"error": `segment already added automatically, 
		it isn't possible to delete it because there are users assigned to this segment permanently`})
		return
	} else if errCheck != sql.ErrNoRows {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error in query process"})
		return
	}

	transaction, errTransaction := db.Begin()
	if errTransaction != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error in transaction process"})
		return
	}
	defer transaction.Rollback()
	_, errQuery := transaction.Exec("DELETE FROM users_segments WHERE segment_id = $1", segmentID)
	if errQuery != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete segment from users"})
		return
	}
	_, errQuery = transaction.Exec("DELETE FROM id_name_segments WHERE segment_id = $1", segmentID)
	if errQuery != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete segment"})
		return
	}
	if err := transaction.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "transaction commit error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Segment deleted successfully"})
}
