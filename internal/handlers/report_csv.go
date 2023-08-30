package handlers

import (
	"encoding/csv"
	"fmt"
	"github.com/edos10/test_avito_service/internal/databases"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func GenerateCSV(c *gin.Context) {
	db, err := databases.CreateDatabaseConnect()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "error to connect database"})
		return
	}
	var requestData struct {
		Year  int     `json:"year"`
		Month int     `json:"month"`
		Users []int64 `json:"users"`
	}

	if err := c.BindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}

	if requestData.Month < 1 || requestData.Month > 12 || requestData.Year > time.Now().Year() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid date"})
		return
	}
	dateStart := time.Date(requestData.Year, time.Month(requestData.Month), 1, 0, 0, 0, 0, time.Local)
	dateEnd := dateStart.AddDate(0, 1, 0).Add(-time.Second)
	query := `
		SELECT ush.user_id, ins.segment_name, ush.operation, ush.timestamp
		FROM user_segment_history ush
		JOIN id_name_segments ins ON ush.segment_id = ins.segment_id
		WHERE ush.timestamp >= $1 AND ush.timestamp <= $2
		AND ush.user_id IN $3
	`
	rows, err := db.Query(query, dateStart, dateEnd, requestData.Users)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch data from database"})
		return
	}
	defer rows.Close()
	var data [][]string
	data = append(data, []string{"идентификатор пользователя", "сегмент", "операция", "дата и время"})
	for rows.Next() {
		var userID int
		var segmentName string
		var operation string
		var timestamp time.Time
		if err := rows.Scan(&userID, &segmentName, &operation, &timestamp); err != nil {
			fmt.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to scan row"})
			return
		}
		data = append(data, []string{
			fmt.Sprintf("%d", userID),
			segmentName,
			operation,
			timestamp.String(),
		})
	}

	c.Header("Content-Disposition", "attachment; filename=report.csv")
	c.Header("Content-Type", "text/csv")

	writer := csv.NewWriter(c.Writer)
	defer writer.Flush()

	for _, record := range data {
		err := writer.Write(record)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to write data to CSV"})
			return
		}
	}
}
