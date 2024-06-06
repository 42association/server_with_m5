package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
	"regexp"
	"errors"
	"42ActivityAPI/internal/accessdb"
)

func GetShiftData(c *gin.Context) {
	date, err := getQueryAboutDate(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid query"})
		return
	}

	//roleがcleaningのactivityを取得
	shifts, err := accessdb.GetShiftFromDB(date)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get shift"})
		return
	}
	c.JSON(http.StatusOK, shifts)
}

func AddShiftData(c *gin.Context) {
	var schedule []accessdb.Schedule

	if err := c.BindJSON(&schedule); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if len(schedule) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Shift is required"})
		return
	}
	if date, err := accessdb.AddShiftToDB(schedule); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{"date": date})
	}
}

func getQueryAboutDate(c *gin.Context) (string, error) {
	date := c.Query("date")
	if date == "" {
		date = time.Now().Format("2006-01-02")
	}
	if !regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`).MatchString(date) {
		return "", errors.New("Invalid date format. It should be in YYYY/MM/DD format")
	}
	return date, nil
}