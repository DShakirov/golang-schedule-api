package controller

import (
	"ScheduleAPI/pkg/model"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

type AddMedicalRecordRequestBody struct {
	PatientID uuid.UUID `json:"patient_id"`
	Text      string    `json:"text"`
}

func GetMedicalRecorsList(db *gorm.DB) func(c *gin.Context) {
	//Request for fetching all MedicalRecord objects belongs to user
	return func(c *gin.Context) {
		userID := c.MustGet("uuid").(uuid.UUID)
		var medicalRecords []model.MedicalRecord
		db.Where("doctor_id = ? OR patient_id = ?", userID, userID).Find(&medicalRecords)
		c.JSON(http.StatusOK, medicalRecords)
	}
}

func GetMedicalRecord(db *gorm.DB) func(c *gin.Context) {
	//Request for fetching MedicalRecord object belongs to user
	return func(c *gin.Context) {
		id := c.Param("id")
		userID := c.MustGet("uuid").(uuid.UUID)
		var medicalRecord model.MedicalRecord
		db.Where("doctor_id = ? OR patient_id = ?", userID, userID).First(&medicalRecord, id)
		c.JSON(http.StatusOK, medicalRecord)
	}
}

func CreateMedicalRecord(db *gorm.DB) func(c *gin.Context) {
	//Request for creating MedicalRecord data
	//Only a doctor can create MedicalRecord
	//IMPORTANT: Structure of request
	//{"text": "deadly hemmoroids diagnosed",
	//"patient_id": "0ec638e3-c9aa-4fd3-9f6d-a738a42a9b5b"}
	//USE POST METHOD
	return func(c *gin.Context) {
		//Retrieving request body
		body := AddMedicalRecordRequestBody{}
		if err := c.BindJSON(&body); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
		//Checking user is doctor
		isDoctor, _ := c.Get("isDoctor")
		if isDoctor != true {
			c.AbortWithError(http.StatusBadRequest, errors.New("Only a doctor can create medical record"))
		}
		//Fetching user id
		uuidParam := c.MustGet("uuid").(uuid.UUID)
		//Creating MedicalRecord object
		var medicalRecord model.MedicalRecord
		medicalRecord.DoctorID = uuidParam
		medicalRecord.PatientID = body.PatientID
		medicalRecord.Text = body.Text
		if result := db.Create(&medicalRecord); result.Error != nil {
			c.AbortWithError(http.StatusNotFound, result.Error)
			return
		}
		c.JSON(http.StatusCreated, medicalRecord)
	}
}

func UpdateMedicalRecord(db *gorm.DB) func(c *gin.Context) {
	//Request for update MedicalRecord data
	//Only a owner can update MedicalRecord
	//IMPORTANT: Structure of request
	//NOTE: Go serializes duration in nanoseconds
	//{"text": "lightly hemmoroids",
	//"patient_id": "0ec638e3-c9aa-4fd3-9f6d-a738a42a9b5b"}
	//USE PUT METHOD
	return func(c *gin.Context) {
		//Fetch medical record
		var medicalRecord model.MedicalRecord
		id := c.Param("id")
		result := db.First(&medicalRecord, id)
		if result.Error != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Failed to fetch medical record"})
			return
		}
		//Check if medical record belongs to user
		uuidParam := c.MustGet("uuid").(uuid.UUID)
		if uuidParam != medicalRecord.DoctorID {
			c.JSON(http.StatusBadRequest, gin.H{"error": "This medical record does not belong to you"})
			return
		}
		//Retrieving request body
		body := AddMedicalRecordRequestBody{}
		if err := c.BindJSON(&body); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
		//Updating MedicalRecord object
		medicalRecord.DoctorID = uuidParam
		medicalRecord.PatientID = body.PatientID
		medicalRecord.Text = body.Text
		if result := db.Save(&medicalRecord); result.Error != nil {
			c.AbortWithError(http.StatusNotFound, result.Error)
			return
		}
		c.JSON(http.StatusOK, medicalRecord)
	}
}

func DeleteMedicalRecord(db *gorm.DB) func(c *gin.Context) {
	//Request for deleting MedicalRecord data
	//Only a owner can delete MedicalRecord
	//USE DELETE METHOD
	return func(c *gin.Context) {
		//Fetch medical record
		var medicalRecord model.MedicalRecord
		id := c.Param("id")
		result := db.First(&medicalRecord, id)
		if result.Error != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Failed to fetch medical record"})
			return
		}
		//Check if appointment belongs to user
		uuidParam := c.MustGet("uuid").(uuid.UUID)
		if uuidParam != medicalRecord.DoctorID {
			c.JSON(http.StatusBadRequest, gin.H{"error": "This medical record does not belong to you"})
			return
		}
		//Deleting object
		db.Delete(&medicalRecord)
		c.JSON(http.StatusNoContent, gin.H{"message": "The object has been succesfully deleted"})
	}
}
