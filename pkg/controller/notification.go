package controller

import (
	"ScheduleAPI/pkg/model"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

func GetNotificationsList(db *gorm.DB) func(c *gin.Context) {
	//Fetching all Notifications objects belongs to user
	return func(c *gin.Context) {
		var notifications []model.Notification
		userID := c.MustGet("uuid").(uuid.UUID)
		db.Where("user_id = ?", userID).Find(&notifications)
		c.JSON(http.StatusOK, notifications)
	}
}

func GetNotification(db *gorm.DB) func(c *gin.Context) {
	//Fetching Notification object belongs to user
	return func(c *gin.Context) {
		userID := c.MustGet("uuid").(uuid.UUID)
		id := c.Param("id")
		fmt.Println()
		var notification model.Notification
		db.Where("user_id = ?", userID).First(&notification, id)
		c.JSON(http.StatusOK, notification)
	}
}
