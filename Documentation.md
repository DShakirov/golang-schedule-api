Doctors Schedule API documentation

MODELS:
    Appointment 
	    DoctorID  UUID
	    PatientID UUID
	    CreatedAt time 
	    TimeStart time
	    TimeEnd   time
    MedicalRecord
	    DoctorID  UUID
	    PatientID UUID
	    CreatedAt time
	    Text      string    
    Notification 
        Type      string
        UserID    UUID
        CreatedAt time
        Text      string
    Prescription
        DrugName  string
        Dosage    string
        Duration  time.Duration
        DoctorID  UUID
        PatientID UUID
        CreatedAt time
    Schedule
        DoctorID  UUID
        TimeStart time
        TimeEnd   time
        CreatedAt time

APIs:
	GET "api/schedules/"
        Fetching all schedule objects
	GET "api/schedules/:id"
        Fetching schedule object by id
	POST "api/schedules/"
            Creating schedule object
        IMPORTANT! Structure of request:
          {"time_start": "2023-12-01T13:00:00Z",
        	"time_end": "2023-12-01T15:00:00Z"}
	PUT "api/schedules/:id"
        Updating schedule object
	    IMPORTANT! Structure of request:
	      {"time_start": "2023-12-01T13:00:00Z",
	    	"time_end": "2023-12-01T15:00:00Z"}
	DELETE "api/schedules/:id"
        Deleting schedule object
	GET "api/appointments/"
        Fetching all Appointment objects belonging to user
	GET "api/appointments/:id"
        Fetching Appointment object belonging to user
	POST "api/appointments"
        Request for creating Appointment data
	    IMPORTANT: Structure of request
	     {"time_start": "2023-12-01T12:00:00Z",
	    "time_end": "2023-12-01T16:00:00Z",
	    "doctor_id": "0ec638e3-c9aa-4fd3-9f6d-a738a42a9b5b",
	    "patient_id": "0ec638e3-c9aa-4fd3-9f6d-a738a42a9b5b"}
	PUT "api/appointments/:id"
        Request for creating Appointment data
	    IMPORTANT: Structure of request
	     {"time_start": "2023-12-01T12:00:00Z",
	    "time_end": "2023-12-01T16:00:00Z",
	    "doctor_id": "0ec638e3-c9aa-4fd3-9f6d-a738a42a9b5b",
	    "patient_id": "0ec638e3-c9aa-4fd3-9f6d-a738a42a9b5b"}
	DELETE "api/appointments/:id"
    	Request for deleting Appointment data
	    Only a user with doctor role can do this
	GET "api/notifications"
        Fetching all Notifications objects belongs to user
	GET "api/notifications/:id"
        Fetching Notification object belongs to user
	GET"api/prescriptions"
        Request for fetching all Prescription objects belongs to user
	GET "api/prescriptions/:id"
        Request for fetching Prescription object belongs to user
	POST "api/prescriptions"
        Request for creating Prescription data
        Only a doctor can create Prescription
        IMPORTANT: Structure of request
        NOTE: Go serializes duration in nanoseconds
        {"dosage": "2 capsules",
        "duration": 1,
        "patient_id": "0ec638e3-c9aa-4fd3-9f6d-a738a42a9b5b"}
	PUT "api/prescriptions/:id"
        Request for update Prescription data
        Only a owner can update Prescription
        IMPORTANT: Structure of request
        NOTE: Go serializes duration in nanoseconds
        {"dosage": "2 capsules",
        "duration": 1,
        "patient_id": "0ec638e3-c9aa-4fd3-9f6d-a738a42a9b5b"}
	DELETE "api/prescriptions/:id"
    	Request for deleting Prescription data
	    Only a owner can delete Prescription
	GET "api/medical_records/"
        Request for fetching all MedicalRecord objects belongs to user
	GET "api/medical_records/:id"
        Request for fetching MedicalRecord object belongs to user
	POST "api/medical_records/"
    	Request for creating MedicalRecord data
        Only a doctor can create MedicalRecord
        IMPORTANT: Structure of request
        {"text": "deadly hemmoroids diagnosed",
        "patient_id": "0ec638e3-c9aa-4fd3-9f6d-a738a42a9b5b"}
	PUT "api/medical_records/:id"
    	Request for update MedicalRecord data
        Only a owner can update MedicalRecord
        IMPORTANT: Structure of request
        NOTE: Go serializes duration in nanoseconds
        {"text": "lightly hemmoroids",
        "patient_id": "0ec638e3-c9aa-4fd3-9f6d-a738a42a9b5b"}
	DELETE "api/medical_records/:id"
    	Request for deleting MedicalRecord data
        Only a owner can delete MedicalRecord