package controller

import (
	"ScheduleAPI/pkg/model"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

type AddAppointmentRequestBody struct {
	DoctorID  uuid.UUID `json:"doctor_id"`
	PatientID uuid.UUID `json:"patient_id"`
	TimeStart time.Time `json:"time_start"`
	TimeEnd   time.Time `json:"time_end"`
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
	//"patient_id": "0ec638e3-c9aa-4fd3-9f6d-a738a42a9b5b"}
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
			if a.TimeStart.After(body.TimeStart) || a.TimeEnd.Before(body.TimeEnd) {
				c.JSON(http.StatusBadRequest, gin.H{"error": "This time is already appointed"})
				return
			}
		}
		//Creating Appointment object
		var appointment model.Appointment
		appointment.DoctorID = body.DoctorID
		appointment.PatientID = body.PatientID
		appointment.TimeStart = body.TimeStart
		appointment.TimeEnd = body.TimeEnd
		if result := db.Create(&appointment); result.Error != nil {
			c.AbortWithError(http.StatusNotFound, result.Error)
			return
		}
		//Creating notifications for both doctor and patient
		var doctor_notification model.Notification
		doctor_notification.Type = "Create"
		doctor_notification.UserID = appointment.DoctorID
		doctor_notification.Text = "Appointment data created"
		if result := db.Save(&doctor_notification); result.Error != nil {
			c.AbortWithError(http.StatusBadRequest, result.Error)
			return
		}
		var patient_notification model.Notification
		patient_notification.Type = "Create"
		patient_notification.UserID = appointment.PatientID
		patient_notification.Text = "Appointment data created"
		if result := db.Save(&patient_notification); result.Error != nil {
			c.AbortWithError(http.StatusBadRequest, result.Error)
			return
		}
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
			if a.TimeStart.After(body.TimeStart) || a.TimeEnd.Before(body.TimeEnd) {
				c.JSON(http.StatusBadRequest, gin.H{"error": "This time is already appointed"})
				return
			}
		}
		//Creating notifications for both doctor and patient
		appointment.DoctorID = body.DoctorID
		appointment.PatientID = body.PatientID
		appointment.TimeStart = body.TimeStart
		appointment.TimeEnd = body.TimeEnd
		if result := db.Save(&appointment); result.Error != nil {
			c.AbortWithError(http.StatusBadRequest, result.Error)
			return
		}
		var doctor_notification model.Notification
		doctor_notification.Type = "Change"
		doctor_notification.UserID = appointment.DoctorID
		doctor_notification.Text = "Appointment data changed"
		if result := db.Save(&doctor_notification); result.Error != nil {
			c.AbortWithError(http.StatusBadRequest, result.Error)
			return
		}
		var patient_notification model.Notification
		patient_notification.Type = "Change"
		patient_notification.UserID = appointment.PatientID
		patient_notification.Text = "Appointment data changed"
		if result := db.Save(&patient_notification); result.Error != nil {
			c.AbortWithError(http.StatusBadRequest, result.Error)
			return
		}

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
		//Deleting appointment
		db.Delete(&appointment)
		//Creating notifications for both doctor and patient
		var doctor_notification model.Notification
		doctor_notification.Type = "Delete"
		doctor_notification.UserID = appointment.DoctorID
		doctor_notification.Text = "Appointment has been cancelled"
		if result := db.Save(&doctor_notification); result.Error != nil {
			c.AbortWithError(http.StatusBadRequest, result.Error)
			return
		}
		var patient_notification model.Notification
		patient_notification.Type = "Delete"
		patient_notification.UserID = appointment.PatientID
		patient_notification.Text = "Appointment has ben cancelled"
		if result := db.Save(&patient_notification); result.Error != nil {
			c.AbortWithError(http.StatusBadRequest, result.Error)
			return
		}
		c.JSON(http.StatusNoContent, gin.H{"msg": "Appointment succesfully deleted"})
	}
}