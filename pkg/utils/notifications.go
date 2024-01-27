package utils

import (
	"ScheduleAPI/pkg/model"
	"log"
	"os"
	"strconv"

	"github.com/gofrs/uuid"

	"gopkg.in/gomail.v2"
	"gorm.io/gorm"
)

func CreateNotification(db *gorm.DB, notificationText, notificationType, userEmail string, userID uuid.UUID) error {
	notification := model.Notification{
		Type:      notificationType,
		UserID:    userID,
		UserEmail: userEmail,
		Text:      notificationText,
	}
	if err := db.Create(&notification).Error; err != nil {
		return err
	}
	//Sending email
	m := gomail.NewMessage()
	m.SetHeader("From", os.Getenv("EMAIL_USER"))
	m.SetHeader("To", userEmail)
	m.SetHeader("Subject", "New doctor appointment!")
	m.SetBody("text/html", notificationText)

	emailHost := os.Getenv("EMAIL_HOST")
	emailPort, err := strconv.Atoi(os.Getenv("EMAIL_PORT"))
	if err != nil {
		panic(err)
	}
	emailUser := os.Getenv("EMAIL_USER")
	emailPassword := os.Getenv("EMAIL_PASSWORD")

	d := gomail.NewDialer(emailHost, emailPort, emailUser, emailPassword)
	if err := d.DialAndSend(m); err != nil {
		log.Print(err)
	}
	return nil
}
