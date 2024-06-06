package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"42ActivityAPI/internal/accessdb"
)

type UserRequestData struct {
	Uid string `json:"uid"`
	Login string `json:"login"`
	Wallet string `json:"wallet"`
}

func AddUser(c *gin.Context) {
	var requestData UserRequestData

	// JSONリクエストボディを解析してrequestDataに格納
	if err := c.BindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if requestData.Login == "" {
		// Loginは必須
		c.JSON(http.StatusBadRequest, gin.H{"error": "Login is required"})
		return
	}
	if accessdb.UserExists(requestData.Login) {
		c.JSON(http.StatusConflict, gin.H{"error": "User with this login already exists"})
		return
	}
	// データベースにUserを追加
	if err := accessdb.AddUserToDB(requestData.Uid, requestData.Login, requestData.Wallet); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	response := make(gin.H)

	response["login"] = requestData.Login
	if requestData.Uid != "" {
		response["uid"] = requestData.Uid
	}
	if requestData.Wallet != "" {
		response["wallet"] = requestData.Wallet
	}

	c.JSON(http.StatusOK, response)
	return
}

func EditUser(c *gin.Context) {
	var requestData UserRequestData

	// JSONリクエストボディを解析してrequestDataに格納
	if err := c.BindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if requestData.Login == "" {
		// Loginは必須
		c.JSON(http.StatusBadRequest, gin.H{"error": "Login is required"})
		return
	}
	// DB上のUserを編集
	if err := accessdb.EditUserInDB(requestData.Uid, requestData.Login, requestData.Wallet); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	response := make(gin.H)

	response["login"] = requestData.Login
	if requestData.Uid != "" {
		response["uid"] = requestData.Uid
	}
	if requestData.Wallet != "" {
		response["wallet"] = requestData.Wallet
	}

	c.JSON(http.StatusOK, response)
	return
}