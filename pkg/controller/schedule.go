package controller

import (
	"ScheduleAPI/pkg/model"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

type AddScheduleRequestBody struct {
	TimeStart time.Time `json:"time_start"`
	TimeEnd   time.Time `json:"time_end"`
}

func GetShedulesList(db *gorm.DB) func(c *gin.Context) {
	//Fetching all schedule objects
	return func(c *gin.Context) {
		var schedules []model.Schedule
		result := db.Find(&schedules)
		if result.Error != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Failed to fetch schedules"})
			return
		}
		c.JSON(http.StatusOK, schedules)
	}
}

func GetScheduleById(db *gorm.DB) func(c *gin.Context) {
	//Fetching schedule object by id
	return func(c *gin.Context) {
		var schedule model.Schedule
		id := c.Param("id")
		result := db.First(&schedule, id)
		if result.Error != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Failed to fetch schedule"})
			return
		}
		c.JSON(http.StatusOK, schedule)
	}
}

func CreateSchedule(db *gorm.DB) func(c *gin.Context) {
	//Creating schedule object
	//IMPORTANT! Structure of request:
	//  {"time_start": "2023-12-01T13:00:00Z",
	//	"time_end": "2023-12-01T15:00:00Z"}
	//USE POST METHOD
	return func(c *gin.Context) {
		// Checking if user is doctor
		isDoctor, _ := c.Get("isDoctor")
		if isDoctor != true {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Only a user with doctor role can create schedule"})
		}
		//Retrieving request body
		body := AddScheduleRequestBody{}
		if err := c.BindJSON(&body); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
		//Checking for invalid values in request
		if body.TimeEnd.Before(body.TimeStart) || body.TimeEnd.Equal(body.TimeStart) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "TimeEnd must be after TimeStart"})
			return
		}
		// Checking the time is not appointed
		doctorID := c.MustGet("uuid").(uuid.UUID)
		var schedules []model.Schedule
		db.Where("doctor_id = ?", doctorID).Find(&schedules)
		for _, s := range schedules {
			if s.TimeStart.Before(body.TimeStart) || s.TimeEnd.After(body.TimeEnd) {
				c.JSON(http.StatusBadRequest, gin.H{"error": "This time is already scheduled"})
				return
			}
		}
		//Creating schedule object
		var schedule model.Schedule
		schedule.DoctorID = doctorID
		schedule.TimeStart = body.TimeStart
		schedule.TimeEnd = body.TimeEnd
		if result := db.Create(&schedule); result.Error != nil {
			c.AbortWithError(http.StatusNotFound, result.Error)
			return
		}
		c.JSON(http.StatusCreated, &schedule)
	}
}

func UpdateSchedule(db *gorm.DB) func(c *gin.Context) {
	//Updating schedule object
	//IMPORTANT! Structure of request:
	//  {"time_start": "2023-12-01T13:00:00Z",
	//	"time_end": "2023-12-01T15:00:00Z"}
	//USE PUT METHOD
	return func(c *gin.Context) {
		//Fetch schedule
		var schedule model.Schedule
		id := c.Param("id")
		result := db.First(&schedule, id)
		if result.Error != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Failed to fetch schedule"})
			return
		}
		//Check if schedule belongs to user
		uuidParam := c.MustGet("uuid").(uuid.UUID)
		if uuidParam != schedule.DoctorID {
			c.JSON(http.StatusBadRequest, gin.H{"error": "This schedule does not belong to you"})
		}
		//Retrieving request body
		body := AddScheduleRequestBody{}
		if err := c.BindJSON(&body); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
		//Checking for invalid values in request
		if body.TimeEnd.Before(body.TimeStart) || body.TimeEnd.Equal(body.TimeStart) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "TimeEnd must be after TimeStart"})
			return
		}
		var schedules []model.Schedule
		db.Where("time_start <= ? AND time_end >= ?", body.TimeEnd, body.TimeStart).Where("doctor_id = ?", uuidParam).Find(&schedules)
		for _, s := range schedules {
			if s.ID == schedule.ID {
				continue //Do nothing with object instance
			}
			if s.TimeStart.Before(body.TimeStart) || s.TimeEnd.After(body.TimeEnd) {
				c.JSON(http.StatusBadRequest, gin.H{"error": "This time is already scheduled"})
				return
			}
		}
		//Updating Schedule object
		schedule.TimeStart = body.TimeStart
		schedule.TimeEnd = body.TimeEnd
		if result := db.Save(&schedule); result.Error != nil {
			c.AbortWithError(http.StatusNotFound, result.Error)
			return
		}
		c.JSON(http.StatusOK, &schedule)
	}
}

func DeleteSchedule(db *gorm.DB) func(c *gin.Context) {
	//Deleting schedule object
	//USE DELETE METHOD
	return func(c *gin.Context) {
		//Fetch schedule
		var schedule model.Schedule
		id := c.Param("id")
		result := db.First(&schedule, id)
		if result.Error != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Failed to fetch schedule"})
			return
		}
		//Check if schedule belongs to user
		uuidParam := c.MustGet("uuid").(uuid.UUID)
		if uuidParam != schedule.DoctorID {
			c.JSON(http.StatusBadRequest, gin.H{"error": "This schedule does not belong to you"})
		}
		//Deleting object
		db.Delete(&schedule)
		c.JSON(http.StatusNoContent, gin.H{"message": "The object has been succesfully deleted"})
	}
}
