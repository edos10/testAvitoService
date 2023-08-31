package handlers

import (
	"database/sql"
	"fmt"
	"github.com/edos10/test_avito_service/internal/databases"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type InfoAddSegments struct {
	SegmentName string    `json:"segment_name"`
	DeleteTime  time.Time `json:"delete_time"`
}

func ChangesUserSegments(g *gin.Context) {

	var requestData struct {
		AddSegments    []InfoAddSegments `json:"adding_segments"`
		RemoveSegments []string          `json:"removing_segments"`
		UserID         int64             `json:"user_id"`
	}
	if err := g.BindJSON(&requestData); err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}

	db, errDb := databases.CreateDatabaseConnect()
	if errDb != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": "couldn't open the database"})
		return
	}
	defer db.Close()

	// хотим до транзакции проверить наличие id пользователя и на всякий случай вставить в таблицу users
	flagInsertUser := false
	queryCheckUserId := "SELECT user_id FROM users WHERE user_id = $1"
	if err := db.QueryRow(queryCheckUserId, requestData.UserID).Scan(nil); err != nil {
		if err == sql.ErrNoRows {
			flagInsertUser = true
		} else {
			g.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}
	}
	curTransaction, err := db.Begin()
	if err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	defer curTransaction.Rollback()
	fmt.Println(flagInsertUser)
	if flagInsertUser {
		queryInsertUser := "INSERT INTO users (user_id) VALUES ($1)"
		if _, err := curTransaction.Exec(queryInsertUser, requestData.UserID); err != nil {
			g.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}
	}

	for _, segmentInfo := range requestData.AddSegments {
		segment := segmentInfo.SegmentName
		timeToDelete := segmentInfo.DeleteTime

		queryCheckSegment := "SELECT segment_id FROM id_name_segments WHERE segment_name = $1"
		var segmentID int
		if err := db.QueryRow(queryCheckSegment, segment).Scan(&segmentID); err != nil {
			if err == sql.ErrNoRows {
				continue
			} else {
				g.JSON(http.StatusInternalServerError, gin.H{"error": "error in query process"})
				return
			}
		}

		queryCheckUserSegment := "SELECT COUNT(*) FROM users_segments WHERE user_id = $1 AND segment_id = $2"
		var count int
		if err := db.QueryRow(queryCheckUserSegment, requestData.UserID, segmentID).Scan(&count); err != nil {
			g.JSON(http.StatusInternalServerError, gin.H{"error": "error in query process"})
			return
		}
		if count > 0 || timeToDelete.Sub(time.Now()) <= 0 {
			continue
		}

		queryAddUsersSegments := `INSERT INTO users_segments 
    				 (user_id, segment_id, endtime) 
					 VALUES ($1, (SELECT segment_id FROM id_name_segments
					 WHERE segment_name = $2), $3)`
		if _, err := curTransaction.Exec(queryAddUsersSegments, requestData.UserID, segment, timeToDelete); err != nil {
			g.JSON(http.StatusInternalServerError, gin.H{"error": "failed to add segment"})
			return
		}

		queryAddHistory := `INSERT INTO user_segment_history
					(user_id, segment_name, operation, timestamp) VALUES 
					($1, $2, $3, $4),
					($5, $6, $7, $8)
		`
		if _, err := curTransaction.Exec(queryAddHistory,
			requestData.UserID, segmentInfo.SegmentName, "add", time.Now().Format("2006-01-02 15:04:05"),
			requestData.UserID, segmentInfo.SegmentName, "remove", timeToDelete); err != nil {

			g.JSON(http.StatusInternalServerError, gin.H{"error": "failed to add segment"})
			return
		}
	}

	for _, segment := range requestData.RemoveSegments {
		queryDel := "DELETE FROM users_segments WHERE user_id = $1 AND segment_id = (SELECT segment_id FROM id_name_segments WHERE segment_name = $2)"
		if _, err := curTransaction.Exec(queryDel, requestData.UserID, segment); err != nil {
			g.JSON(http.StatusInternalServerError, gin.H{"error": "failed to remove segment"})
			return
		}
	}

	if err := curTransaction.Commit(); err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": "transaction commit error"})
		return
	}

	g.JSON(http.StatusOK, gin.H{"message": "segments updated successfully"})
}
