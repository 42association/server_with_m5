package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"github.com/jinzhu/now"
	"errors"
	"42ActivityAPI/internal/accessdb"
)

type ActivityRequestData struct {
	Mac string `json:"mac"`
	Uid string `json:"uid"`
}

func GetActivityCleanData(c *gin.Context) {
	//start_timeとend_timeを取得
	start_time, end_time, err := getQueryAboutTime(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid query"})
		return
	}

	//roleがcleaningのactivityを取得
	Activities, err := accessdb.GetActivitiesFromDB(start_time, end_time, "cleaning")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get activities"})
		return
	}
	c.JSON(http.StatusOK, Activities)
}

func AddActivity(c *gin.Context) {
	var requestData ActivityRequestData

	if err := c.BindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if requestData.Mac == "" || requestData.Uid == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "All parameters are required"})
		return
	}

	status, uid, mac, err := accessdb.AddActivityToDB(requestData.Uid, requestData.Mac)
	if err != nil {
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(status, gin.H{"uid": uid, "mac": mac})
	return
}

func getQueryAboutTime(c *gin.Context) (int64, int64, error) {
	var start_time int64
	var end_time int64
	var err error

	start := c.Query("start")
	if start == "" {
		//startパラメータがない場合は当日の0時0分0秒を取得
		start_time = now.BeginningOfDay().Unix()
	} else {
		start_time, err = strconv.ParseInt(start, 10, 64)
		if err != nil {
			return 0, 0, err
		}
	}
	end := c.Query("end")
	if end == "" {
		//endパラメータがない場合はstart_timeの24時間後を取得
		end_time = start_time + 24*60*60
	} else {
		end_time, err = strconv.ParseInt(end, 10, 64)
		if err != nil {
			return 0, 0, err
		}
	}
	if start_time >= end_time {
		return 0, 0, errors.New("Invalid time range")
	}
	return start_time, end_time, nil
}