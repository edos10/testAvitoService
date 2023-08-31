package handlers

import (
	"database/sql"
	"github.com/edos10/test_avito_service/internal/databases"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func DeleteSegment(c *gin.Context) {
	var requestData struct {
		SegmentName string `json:"segment_name"`
	}
	if err := c.BindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}
	if requestData.SegmentName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"Error": `invalid data, doesn't exit field "segment_name" or it's empty`})
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

	// проверка на существование сегмента
	if err := db.QueryRow(checkQuery, requestData.SegmentName).Scan(&segmentID); err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "segment not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error in query process", "err from db": err})
		}
		return
	}

	// проверочка на то, что сегмент не был добавлен автоматически некоторым пользователям
	// если добавлен - не удаляем его
	queryForCheckAutomaticSegment := `
		SELECT ush.user_id
		FROM user_segment_history ush
		JOIN id_name_segments ins ON ush.segment_id = ins.segment_id
		WHERE ush.timestamp = 'infinity'
	`

	var Id int64
	errCheck := db.QueryRow(queryForCheckAutomaticSegment).Scan(&Id)
	if errCheck == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": `segment already added automatically, 
		it isn't possible to delete it because there are users assigned to this segment permanently`})
		return
	} else if errCheck != sql.ErrNoRows {
		c.JSON(http.StatusInternalServerError, gin.H{"error": errCheck})
		return
	}

	transaction, errTransaction := db.Begin()
	if errTransaction != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": errTransaction})
		return
	}
	defer transaction.Rollback()
	_, errQuery := transaction.Exec("DELETE FROM users_segments WHERE segment_id = $1", segmentID)
	if errQuery != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete segment from users", "err from db": errQuery})
		return
	}
	_, errQuery = transaction.Exec("DELETE FROM id_name_segments WHERE segment_id = $1", segmentID)
	if errQuery != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete segment", "err from db": errQuery})
		return
	}
	_, errQuery = transaction.Exec("UPDATE user_segment_history SET timestamp = $1 WHERE timestamp > $1", time.Now().Format("2006-01-02 15:04:05"))
	if errQuery != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update history", "error from db": errQuery})
		return
	}
	if err := transaction.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "transaction commit error", "err from db": err})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Segment deleted successfully"})
}
