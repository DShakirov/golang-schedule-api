package main

import (
	"ScheduleAPI/pkg/config"
	"ScheduleAPI/pkg/controller"
	"ScheduleAPI/pkg/middleware"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	//Initialize Gin Router
	r := gin.Default()

	//Initialize DB connection
	db := config.SetupDatabaseConnection()
	defer config.CloseDatabaseConnection(db)

	//Adding middleware to router
	r.Use(middleware.AuthMiddleware(db))

	//Declaring API routes
	//Schedule objects rotes
	r.GET("api/schedules/", controller.GetShedulesList(db))
	r.GET("api/schedules/:id", controller.GetScheduleById(db))
	r.POST("api/schedules/", controller.CreateSchedule(db))
	r.PUT("api/schedules/:id", controller.UpdateSchedule(db))
	r.DELETE("api/schedules/:id", controller.DeleteSchedule(db))
	//Appointment objects routes
	r.GET("api/appointments/", controller.GetAppointmentsList(db))
	r.GET("api/appointments/:id", controller.GetAppointment(db))
	r.POST("api/appointments", controller.CreateAppointment(db))
	r.PUT("api/appointments/:id", controller.UpdateAppointment(db))
	r.DELETE("api/appointments/:id", controller.DeleteAppointment(db))
	//Notification objects routes
	r.GET("api/notifications", controller.GetNotificationsList(db))
	r.GET("api/notifications/:id", controller.GetNotification(db))
	//Prescription objects routes
	r.GET("api/prescriptions", controller.GetPrescriptionList(db))
	r.GET("api/prescriptions/:id", controller.GetPrescription(db))
	r.POST("api/prescriptions", controller.CreatePrescription(db))
	r.PUT("api/prescriptions/:id", controller.UpdatePrescription(db))
	r.DELETE("api/prescriptions/:id", controller.DeletePrescription(db))
	//MedicalRecord objects routes
	r.GET("api/medical_records/", controller.GetMedicalRecorsList(db))
	r.GET("api/medical_records/:id", controller.GetMedicalRecord(db))
	r.POST("api/medical_records/", controller.CreateMedicalRecord(db))
	r.PUT("api/medical_records/:id", controller.UpdateMedicalRecord(db))
	r.DELETE("api/medical_records/:id", controller.DeleteMedicalRecord(db))
	//start router
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to run server:", err)
	}
}
