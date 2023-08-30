package handlers

import (
	"database/sql"
	"fmt"
	"github.com/edos10/test_avito_service/internal/databases"
	"github.com/gin-gonic/gin"
	"math"
	"math/rand"
	"net/http"
)

func GetAllUserID(db *sql.DB) ([]int64, error) {
	query := "SELECT DISTINCT user_id FROM users_segments"
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var userIDs []int64
	for rows.Next() {
		var userID int64
		if err := rows.Scan(&userID); err != nil {
			return nil, err
		}
		userIDs = append(userIDs, userID)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return userIDs, nil
}

func GetRandSlice(in []int64) []int64 {
	ans := make([]int64, len(in), len(in))
	shuffle := rand.Perm(len(in))
	for i, v := range shuffle {
		ans[i] = in[v]
	}
	return ans
}

func CreateSegment(c *gin.Context) {
	var requestData struct {
		SegmentName string `json:"segment_name"`
		Percents    int    `json:"percents"`
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
	if !(0 <= requestData.Percents && requestData.Percents <= 100) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON field - percents, please, make this field the number from 0 to 100 including"})
		return
	}
	var existID int64
	queryCheck := "SELECT segment_id FROM id_name_segments WHERE segment_name = $1"
	errCheck := db.QueryRow(queryCheck, requestData.SegmentName).Scan(&existID)
	if errCheck == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "segment already exists"})
		return
	} else if errCheck != sql.ErrNoRows {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error in query process"})
		return
	}

	transaction, errTx := db.Begin()
	if errTx != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error starting transaction"})
		return
	}
	defer transaction.Rollback()

	var newSegmentID int
	queryInsert := "INSERT INTO id_name_segments (segment_name) VALUES ($1) RETURNING segment_id"
	errIns := transaction.QueryRow(queryInsert, requestData.SegmentName).Scan(&newSegmentID)
	if errIns != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to insert segment"})
		return
	}

	if err := transaction.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "transaction commit error"})
		return
	}

	transaction, errTx = db.Begin()
	if errTx != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error starting transaction"})
		return
	}
	defer transaction.Rollback()

	if requestData.Percents < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid percents for users, please try again, percents >= 0"})
		return
	}

	allUsers, errGet := GetAllUserID(db)
	if errGet != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "users don't receive this segment because all user IDs aren't retrieved"})
		return
	}
	newRandSlice := GetRandSlice(allUsers)
	countElems := int(math.Ceil(float64(requestData.Percents) * float64(len(allUsers)) / 100.0))
	fmt.Println(newRandSlice)
	for i := 0; i < countElems; i++ {
		curUser := newRandSlice[i]
		queryInsUserSeg := "INSERT INTO users_segments (user_id, segment_id, endtime) VALUES ($1, $2, 'infinity')"
		_, err := transaction.Exec(queryInsUserSeg, curUser, newSegmentID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to insert user segment"})
			return
		}
		queryInsHist := `INSERT INTO user_segment_history (user_id, segment_id, operation, timestamp)
						 VALUES ($1, $2, $3, 'infinity')`
		_, err = transaction.Exec(queryInsHist, curUser, newSegmentID, "add")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to insert add operation"})
			return
		}
	}
	if err := transaction.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "transaction commit error"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "segment is successfully created"})
}
