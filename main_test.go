package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"bytes"
	"encoding/json"

	"ScheduleAPI/pkg/config"
	"ScheduleAPI/pkg/controller"
	"ScheduleAPI/pkg/middleware"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestCreateScheduleHandler(t *testing.T) {
	// Create test data
	payload := map[string]interface{}{
		"time_start": "2023-12-01T20:00:00Z",
		"time_end":   "2023-12-01T21:00:00Z",
	}
	jsonPayload, _ := json.Marshal(payload)
	// Declaring router
	r := gin.Default()
	// Declaring DB
	db := config.SetupDatabaseConnection()
	defer config.CloseDatabaseConnection(db)
	//Adding middleware
	r.Use(middleware.AuthMiddleware(db))
	//Adding route
	r.POST("api/schedules/create/", controller.CreateSchedule(db))

	// Create POST request
	req, err := http.NewRequest("POST", "/api/schedules/create", bytes.NewBuffer(jsonPayload))
	if err != nil {
		t.Fatal(err)
	}

	// Adding authorization header
	req.Header.Set("Authorization", "Bearer your_token")

	// Creating fake ResponseWriter and request handler
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	//responseData, _ := ioutil.ReadAll(w.Body)
	assert.Equal(t, http.StatusCreated, w.Code)
}
