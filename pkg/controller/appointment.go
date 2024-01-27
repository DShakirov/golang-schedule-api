package controller

import (
	"ScheduleAPI/pkg/model"
	"ScheduleAPI/pkg/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

type AddAppointmentRequestBody struct {
	DoctorID     uuid.UUID `json:"doctor_id"`
	DoctorEmail  string    `json:"doctor_email"`
	PatientID    uuid.UUID `json:"patient_id"`
	PatientEmail string    `json:"patient_email"`
	TimeStart    time.Time `json:"time_start"`
	TimeEnd      time.Time `json:"time_end"`
}

func GetAppointmentsList(db *gorm.DB) func(c *gin.Context) {
	//Fetching all Appointment objects belonging to user
	return func(c *gin.Context) {
		//Retrieving user ID from context
		userID := c.MustGet("uuid").(uuid.UUID)
		//Fetching all objects belongs to user
		var appointments []model.Appointment
		db.Where("doctor_id = ? OR patient_id = ?", userID).Find(&appointments)
		c.JSON(http.StatusOK, appointments)
	}
}

func GetAppointment(db *gorm.DB) func(c *gin.Context) {
	//Fetching Appointment object belonging to user
	return func(c *gin.Context) {
		//Retrieving user ID from context
		userID := c.MustGet("uuid").(uuid.UUID)
		//Retirieving object ID from context
		id := c.Param("id")
		var appointment model.Appointment
		db.Where("doctor_id = ? OR patient_id = ?", userID).First(&appointment, id)
		c.JSON(http.StatusOK, appointment)
	}
}

func CreateAppointment(db *gorm.DB) func(c *gin.Context) {
	//Request for creating Appointment data
	//IMPORTANT: Structure of request
	// {"time_start": "2023-12-01T12:00:00Z",
	//"time_end": "2023-12-01T16:00:00Z",
	//"doctor_id": "0ec638e3-c9aa-4fd3-9f6d-a738a42a9b5b",
	//"doctor_email":"doctor@test.com",
	//"patient_id": "0ec638e3-c9aa-4fd3-9f6d-a738a42a9b5b"
	//"patient_email": "patient@test.com"}
	//USE POST METHOD

	return func(c *gin.Context) {
		//Retrieving request body
		body := AddAppointmentRequestBody{}
		if err := c.BindJSON(&body); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
		//Checking for invalid values in request
		if body.TimeEnd.Before(body.TimeStart) || body.TimeEnd.Equal(body.TimeStart) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "TimeEnd must be after TimeStart"})
			return
		}
		doctorEmailValidate := utils.IsValidEmail(body.DoctorEmail)
		patientEmailValidate := utils.IsValidEmail(body.PatientEmail)
		if doctorEmailValidate != true || patientEmailValidate != true {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email"})
			return
		}
		//Checking if Shedule exists
		var schedule model.Schedule
		db.Where("time_start <= ? AND time_end >= ? AND doctor_id = ?", body.TimeStart, body.TimeEnd, body.DoctorID).First(&schedule)
		if schedule.ID == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "No free schedules for your request"})
			return
		}
		//Checking if time is not appointed
		var appointments []model.Appointment
		db.Where("time_start <= ? AND time_end >= ?", body.TimeEnd, body.TimeStart).Where("doctor_id = ?", body.DoctorID).Find(&appointments)
		for _, a := range appointments {
			if a.TimeStart.Equal(body.TimeStart) || a.TimeEnd.Equal(body.TimeEnd) || (a.TimeStart.Before(body.TimeStart) && a.TimeEnd.After(body.TimeStart)) || (a.TimeStart.Before(body.TimeEnd) && a.TimeEnd.After(body.TimeEnd)) {
				c.JSON(http.StatusBadRequest, gin.H{"error": "This time is already appointed"})
				return
			}
		}
		//Creating Appointment object
		var appointment model.Appointment
		appointment.DoctorID = body.DoctorID
		appointment.DoctorEmail = body.DoctorEmail
		appointment.PatientID = body.PatientID
		appointment.PatientEmail = body.PatientEmail
		appointment.TimeStart = body.TimeStart
		appointment.TimeEnd = body.TimeEnd
		if result := db.Create(&appointment); result.Error != nil {
			c.AbortWithError(http.StatusNotFound, result.Error)
			return
		}
		//Creating notifications for doctor and patient
		notificationType := "Create"
		notificationText := "Appointment data created"
		doctorEmail := body.DoctorEmail
		doctorID := body.DoctorID
		patientEmail := body.PatientEmail
		patientID := body.PatientID
		utils.CreateNotification(db, notificationText, notificationType, doctorEmail, doctorID)
		utils.CreateNotification(db, notificationText, notificationType, patientEmail, patientID)
		c.JSON(http.StatusCreated, appointment)
	}
}

func UpdateAppointment(db *gorm.DB) func(c *gin.Context) {
	//Request for creating Appointment data
	//IMPORTANT: Structure of request
	// {"time_start": "2023-12-01T12:00:00Z",
	//"time_end": "2023-12-01T16:00:00Z",
	//"doctor_id": "0ec638e3-c9aa-4fd3-9f6d-a738a42a9b5b",
	//"patient_id": "0ec638e3-c9aa-4fd3-9f6d-a738a42a9b5b"}
	//USE PUT METHOD
	return func(c *gin.Context) {
		//Retrieving Appointment object
		var appointment model.Appointment
		if err := db.First(&appointment, c.Param("id")).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Appointment not found"})
			return
		}
		//Retrieving request body
		body := AddAppointmentRequestBody{}
		if err := c.BindJSON(&body); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
		//Checking for invalid values in request
		if body.TimeEnd.Before(body.TimeStart) || body.TimeEnd.Equal(body.TimeStart) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "TimeEnd must be after TimeStart"})
			return
		}
		doctorEmailValidate := utils.IsValidEmail(body.DoctorEmail)
		patientEmailValidate := utils.IsValidEmail(body.PatientEmail)
		if doctorEmailValidate != true || patientEmailValidate != true {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email"})
			return
		}
		//Checking if Shedule exists
		var schedule model.Schedule
		db.Where("time_start <= ? AND time_end >= ? AND doctor_id = ?", body.TimeStart, body.TimeEnd, body.DoctorID).First(&schedule)
		if schedule.ID == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "No free schedules for your request"})
			return
		}
		//Checking if time is not appointed
		var appointments []model.Appointment
		db.Where("time_start <= ? AND time_end >= ?", body.TimeEnd, body.TimeStart).Where("doctor_id = ?", body.DoctorID).Find(&appointments)
		for _, a := range appointments {
			if a.ID == appointment.ID {
				continue //Do nothing with object instance
			}
			if a.TimeStart.Equal(body.TimeStart) || a.TimeEnd.Equal(body.TimeEnd) || (a.TimeStart.Before(body.TimeStart) && a.TimeEnd.After(body.TimeStart)) || (a.TimeStart.Before(body.TimeEnd) && a.TimeEnd.After(body.TimeEnd)) {
				c.JSON(http.StatusBadRequest, gin.H{"error": "This time is already appointed"})
				return
			}
		}
		//Creating appointment
		appointment.DoctorID = body.DoctorID
		appointment.DoctorEmail = body.DoctorEmail
		appointment.PatientID = body.PatientID
		appointment.PatientEmail = body.PatientEmail
		appointment.TimeStart = body.TimeStart
		appointment.TimeEnd = body.TimeEnd
		if result := db.Save(&appointment); result.Error != nil {
			c.AbortWithError(http.StatusBadRequest, result.Error)
			return
		}
		//Creating notifications for doctor and patient
		notificationType := "Change"
		notificationText := "Appointment data has changed"
		doctorEmail := appointment.DoctorEmail
		doctorID := appointment.DoctorID
		patientEmail := appointment.PatientEmail
		patientID := appointment.PatientID
		utils.CreateNotification(db, notificationText, notificationType, doctorEmail, doctorID)
		utils.CreateNotification(db, notificationText, notificationType, patientEmail, patientID)
		c.JSON(http.StatusCreated, appointment)

		c.JSON(http.StatusOK, appointment)
	}
}

func DeleteAppointment(db *gorm.DB) func(c *gin.Context) {
	//Request for deleting Appointment data
	//Only a user with doctor role can do this
	return func(c *gin.Context) {
		//Fetch appointment
		var appointment model.Appointment
		id := c.Param("id")
		result := db.First(&appointment, id)
		if result.Error != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Cannot fetch appointment"})
		}
		//Checking user is doctor
		isDoctor, _ := c.Get("isDoctor")
		if isDoctor != true {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Only a doctor can delete appointment"})
		}
		//Creating notifications for both doctor and patient
		notificationType := "Cancel"
		notificationText := "Appointment data has cancelled"
		doctorEmail := appointment.DoctorEmail
		doctorID := appointment.DoctorID
		patientEmail := appointment.PatientEmail
		patientID := appointment.PatientID
		utils.CreateNotification(db, notificationText, notificationType, doctorEmail, doctorID)
		utils.CreateNotification(db, notificationText, notificationType, patientEmail, patientID)
		//Deleting appointment
		db.Delete(&appointment)

		c.JSON(http.StatusNoContent, gin.H{"msg": "Appointment succesfully deleted"})
	}
}
