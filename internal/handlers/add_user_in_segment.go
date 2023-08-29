package handlers

import (
	"database/sql"
	"fmt"
	"github.com/edos10/test_avito_service/internal/databases"
	"github.com/gin-gonic/gin"
	"net/http"
)

func ChangesUserSegments(g *gin.Context) {
	var requestData struct {
		AddSegments    []string `json:"adding_segments"`
		RemoveSegments []string `json:"removing_segments"`
		UserID         int      `json:"user_id"`
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
	curTransaction, err := db.Begin()
	if err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": "error in transaction"})
		return
	}
	defer curTransaction.Rollback()
	fmt.Println(requestData.AddSegments)
	for _, segment := range requestData.AddSegments {
		queryCheck := "SELECT segment_id FROM id_name_segments WHERE segment_name = $1"
		var segmentID int
		if err := db.QueryRow(queryCheck, segment).Scan(&segmentID); err != nil {
			if err == sql.ErrNoRows {
				continue
			} else {
				g.JSON(http.StatusInternalServerError, gin.H{"error": "error in query process"})
				return
			}
		}
		queryAdd := "INSERT INTO users_segments (user_id, segment_id) VALUES ($1, (SELECT segment_id FROM id_name_segments WHERE segment_name = $2))"
		if _, err := curTransaction.Exec(queryAdd, requestData.UserID, segment); err != nil {
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
