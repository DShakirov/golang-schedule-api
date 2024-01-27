package controller

import (
	"ScheduleAPI/pkg/model"
	"ScheduleAPI/pkg/utils"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

type AddPrescriptionRequestBody struct {
	PatientID    uuid.UUID     `json:"patient_id"`
	PatientEmail string        `json:"patient_email"`
	Dosage       string        `json:"dosage"`
	Duration     time.Duration `json:"duration"`
	DrugName     string        `json:"drug_name"`
}

func GetPrescriptionList(db *gorm.DB) func(c *gin.Context) {
	//Request for fetching all Prescription objects belongs to user
	return func(c *gin.Context) {
		userID := c.MustGet("uuid").(uuid.UUID)
		var prescriptions []model.Prescription
		db.Where("doctor_id = ? OR patient_id = ?", userID, userID).Find(&prescriptions)
		c.JSON(http.StatusOK, prescriptions)
	}
}

func GetPrescription(db *gorm.DB) func(c *gin.Context) {
	//Request for fetching Prescription object belongs to user
	return func(c *gin.Context) {
		id := c.Param("id")
		userID := c.MustGet("uuid").(uuid.UUID)
		var prescription model.Prescription
		db.Where("doctor_id = ? OR patient_id = ?", userID, userID).First(&prescription, id)
		c.JSON(http.StatusOK, prescription)
	}
}

func CreatePrescription(db *gorm.DB) func(c *gin.Context) {
	//Request for creating Prescription data
	//Only a doctor can create Prescription
	//IMPORTANT: Structure of request
	//NOTE: Go serializes duration in nanoseconds
	//{"dosage": "2 capsules",
	//"duration": 1,
	//"patient_id": "0ec638e3-c9aa-4fd3-9f6d-a738a42a9b5b",
	//"patient_email": "patient@test.com"}
	//USE POST METHOD
	return func(c *gin.Context) {
		//Retrieving request body
		body := AddPrescriptionRequestBody{}
		if err := c.BindJSON(&body); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
		patientEmail := body.PatientEmail
		validatePatientEmail := utils.IsValidEmail(patientEmail)
		if validatePatientEmail != true {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email"})
			return
		}
		//Checking user is doctor
		isDoctor, _ := c.Get("isDoctor")
		if isDoctor != true {
			c.AbortWithError(http.StatusBadRequest, errors.New("Only a doctor can create prescription"))
		}
		//Fetching user id and email
		uuidParam := c.MustGet("uuid").(uuid.UUID)
		doctorEmail := c.MustGet("email").(string)
		//Creating Prescription object
		var prescription model.Prescription
		prescription.DoctorID = uuidParam
		prescription.DoctorEmail = doctorEmail
		prescription.PatientID = body.PatientID
		prescription.PatientEmail = body.PatientEmail
		prescription.DrugName = body.DrugName
		prescription.Duration = body.Duration
		prescription.Dosage = body.Dosage
		if result := db.Create(&prescription); result.Error != nil {
			c.AbortWithError(http.StatusNotFound, result.Error)
			return
		}
		c.JSON(http.StatusCreated, prescription)
	}
}

func UpdatePrescription(db *gorm.DB) func(c *gin.Context) {
	//Request for update Prescription data
	//Only a owner can update Prescription
	//IMPORTANT: Structure of request
	//NOTE: Go serializes duration in nanoseconds
	//{"dosage": "2 capsules",
	//"duration": 1,
	//"patient_id": "0ec638e3-c9aa-4fd3-9f6d-a738a42a9b5b",
	//"patient_email": "patient@test.com"}
	//USE PUT METHOD
	return func(c *gin.Context) {
		//Fetch prescription
		var prescription model.Prescription
		id := c.Param("id")
		result := db.First(&prescription, id)
		if result.Error != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Failed to fetch prescription"})
			return
		}
		//Check if prescription belongs to user
		uuidParam := c.MustGet("uuid").(uuid.UUID)
		if uuidParam != prescription.DoctorID {
			c.JSON(http.StatusBadRequest, gin.H{"error": "This prescription does not belong to you"})
			return
		}
		//Fetching user email
		doctorEmail := c.MustGet("email").(string)
		//Retrieving request body
		body := AddPrescriptionRequestBody{}
		if err := c.BindJSON(&body); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
		patientEmail := body.PatientEmail
		validatePatientEmail := utils.IsValidEmail(patientEmail)
		if validatePatientEmail != true {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email"})
			return
		}
		//Updating Prescription object
		prescription.DoctorID = uuidParam
		prescription.DoctorEmail = doctorEmail
		prescription.PatientID = body.PatientID
		prescription.PatientEmail = body.PatientEmail
		prescription.DrugName = body.DrugName
		prescription.Duration = body.Duration
		prescription.Dosage = body.Dosage
		if result := db.Save(&prescription); result.Error != nil {
			c.AbortWithError(http.StatusNotFound, result.Error)
			return
		}
		c.JSON(http.StatusOK, prescription)
	}
}

func DeletePrescription(db *gorm.DB) func(c *gin.Context) {
	//Request for deleting Prescription data
	//Only a owner can delete Prescription
	//USE DELETE METHOD
	return func(c *gin.Context) {
		//Fetch prescription
		var prescription model.Prescription
		id := c.Param("id")
		result := db.First(&prescription, id)
		if result.Error != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Failed to fetch prescription"})
			return
		}
		//Check if appointment belongs to user
		uuidParam := c.MustGet("uuid").(uuid.UUID)
		if uuidParam != prescription.DoctorID {
			c.JSON(http.StatusBadRequest, gin.H{"error": "This prescription does not belong to you"})
			return
		}
		//Deleting object
		db.Delete(&prescription)
		c.JSON(http.StatusNoContent, gin.H{"message": "The object has been succesfully deleted"})
	}
}
