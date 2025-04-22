package patient

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	patientsHelper "github.com/harshgupta9473/assignment_makerable/helpers/patients"
	"github.com/harshgupta9473/assignment_makerable/middlewares"
	"github.com/harshgupta9473/assignment_makerable/models"
	"github.com/harshgupta9473/assignment_makerable/utils"
)

func RegisterPatient(w http.ResponseWriter, r *http.Request) {

	var patientreq models.PatientRequest
	err := json.NewDecoder(r.Body).Decode(&patientreq)
	if err != nil {
		utils.WriteJson(w, http.StatusBadRequest, utils.APIResponse{
			Status:  "error",
			Message: "Invalid request body",
			Error:   err.Error(),
		})
		return
	}

	if patientreq.DoctorDiagnosis == "" {
		patientreq.DoctorDiagnosis = "NULL"
	}

	patientID, err := patientsHelper.CreatePatient(patientreq)
	if err != nil {
		utils.WriteJson(w, http.StatusInternalServerError, utils.APIResponse{
			Status:  "error",
			Message: "Failed to register patient",
			Error:   err.Error(),
		})
		return
	}

	patient := models.Patient{
		ID:              patientID,
		FirstName:       patientreq.FirstName,
		LastName:        patientreq.LastName,
		Email:           patientreq.Email,
		DOB:             patientreq.DOB,
		Phone:           patientreq.Phone,
		Gender:          patientreq.Gender,
		ComplaintType:   patientreq.ComplaintType,
		DoctorDiagnosis: patientreq.DoctorDiagnosis,
		DoctorID:        patientreq.DoctorId,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	utils.WriteJson(w, http.StatusCreated, utils.APIResponse{
		Status:  "success",
		Message: "Patient registered successfully",
		Data:    patient,
	})
}

// for doctors
func GetAllPatientsByDoctors(w http.ResponseWriter, r *http.Request) {
	user, err := middlewares.GetUserContext(r)
	if err != nil {
		utils.WriteJson(w, http.StatusUnauthorized, utils.APIResponse{
			Status:  "error",
			Message: "Unable to retrieve user from context",
		})
		return
	}

	patients, err := patientsHelper.GetPatientsByDoctorID(user.UserID)
	if err != nil {
		utils.WriteJson(w, http.StatusInternalServerError, utils.APIResponse{
			Status:  "error",
			Message: fmt.Sprintf("Error fetching patients: %v", err),
		})
		return
	}

	utils.WriteJson(w, http.StatusOK, utils.APIResponse{
		Status:  "success",
		Message: "Patients fetched successfully",
		Data:    patients,
	})
}

// for doctors  /{patientID}
func GetAPatientPatientID(w http.ResponseWriter, r *http.Request) {
	user, err := middlewares.GetUserContext(r)
	if err != nil {
		utils.WriteJson(w, http.StatusUnauthorized, utils.APIResponse{
			Status:  "error",
			Message: "Unable to retrieve user from context",
		})
		return
	}

	params := mux.Vars(r)
	patientIDStr := params["patientID"]
	patientID, err := strconv.ParseInt(patientIDStr, 10, 64)
	if err != nil || patientID <= 0 {
		utils.WriteJson(w, http.StatusBadRequest, utils.APIResponse{
			Status:  "error",
			Message: "Invalid patient ID",
		})
		return
	}

	patient, err := patientsHelper.GetPatientByDoctorIDAndPatientID(user.UserID, patientID)
	if err != nil {
		utils.WriteJson(w, http.StatusInternalServerError, utils.APIResponse{
			Status:  "error",
			Message: fmt.Sprintf("Error fetching patient: %v", err),
		})
		return
	}

	utils.WriteJson(w, http.StatusOK, utils.APIResponse{
		Status:  "success",
		Message: "Patient fetched successfully",
		Data:    patient,
	})
}




// receptionist
func GetAllPatients(w http.ResponseWriter, r *http.Request) {
	patients, err := patientsHelper.GetAllPatients()
	if err != nil {
		utils.WriteJson(w, http.StatusInternalServerError, utils.APIResponse{
			Status:  "error",
			Message: fmt.Sprintf("Error fetching patients: %v", err),
		})
		return
	}

	utils.WriteJson(w, http.StatusOK, utils.APIResponse{
		Status:  "success",
		Message: "All patients fetched successfully",
		Data:    patients,
	})
}

func GetPatientByID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	patientIDStr := params["patientID"]
	patientID, err := strconv.ParseInt(patientIDStr, 10, 64)
	if err != nil || patientID <= 0 {
		utils.WriteJson(w, http.StatusBadRequest, utils.APIResponse{
			Status:  "error",
			Message: "Invalid patient ID",
		})
		return
	}
	if err != nil || patientID <= 0 {
		utils.WriteJson(w, http.StatusBadRequest, utils.APIResponse{
			Status:  "error",
			Message: "Invalid JSON body or patient ID",
		})
		return
	}

	patient, err := patientsHelper.GetPatient(patientID)
	if err != nil {
		utils.WriteJson(w, http.StatusInternalServerError, utils.APIResponse{
			Status:  "error",
			Message: fmt.Sprintf("Error fetching patient: %v", err),
		})
		return
	}

	utils.WriteJson(w, http.StatusOK, utils.APIResponse{
		Status:  "success",
		Message: "Patient fetched successfully",
		Data:    patient,
	})
}

// only by receptionist
func DeletePatientById(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	patientIDStr := params["patientID"]
	patientID, err := strconv.ParseInt(patientIDStr, 10, 64)
	if err != nil || patientID <= 0 {
		utils.WriteJson(w, http.StatusBadRequest, utils.APIResponse{
			Status:  "error",
			Message: "Invalid patient ID",
		})
		return
	}

	err = patientsHelper.SoftDeletePatientByID(patientID)
	if err != nil {
		utils.WriteJson(w, http.StatusInternalServerError, utils.APIResponse{
			Status:  "error",
			Message: fmt.Sprintf("Error deleting patient: %v", err),
		})
		return
	}

	utils.WriteJson(w, http.StatusOK, utils.APIResponse{
		Status:  "success",
		Message: "Patient deleted successfully",
	})
}

func UpdatePatientHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	patientIDStr := params["patientID"]
	patientID, err := strconv.ParseInt(patientIDStr, 10, 64)
	if err != nil || patientID <= 0 {
		utils.WriteJson(w, http.StatusBadRequest, utils.APIResponse{
			Status:  "error",
			Message: "Invalid patient ID",
		})
		return
	}

	var patient models.Patient
	if err := json.NewDecoder(r.Body).Decode(&patient); err != nil {
		utils.WriteJson(w, http.StatusBadRequest, utils.APIResponse{
			Status:  "error",
			Message: "Invalid JSON body",
			Error:   err.Error(),
		})
		return
	}

	patient.ID = patientID

	if err := patientsHelper.UpdatePatient(patient); err != nil {
		utils.WriteJson(w, http.StatusInternalServerError, utils.APIResponse{
			Status:  "error",
			Message: "Failed to update patient",
			Error:   err.Error(),
		})
		return
	}

	utils.WriteJson(w, http.StatusOK, utils.APIResponse{
		Status:  "success",
		Message: "Patient updated successfully",
	})
}

func UpdateDoctorDiagnosisHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	patientIDStr := params["patientID"]
	patientID, err := strconv.ParseInt(patientIDStr, 10, 64)
	if err != nil || patientID <= 0 {
		utils.WriteJson(w, http.StatusBadRequest, utils.APIResponse{
			Status:  "error",
			Message: "Invalid patient ID",
		})
		return
	}

	var requestBody struct {
		DoctorDiagnosis string `json:"doctor_diagnosis"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		utils.WriteJson(w, http.StatusBadRequest, utils.APIResponse{
			Status:  "error",
			Message: "Invalid JSON body",
			Error:   err.Error(),
		})
		return
	}

	if requestBody.DoctorDiagnosis == "" {
		utils.WriteJson(w, http.StatusBadRequest, utils.APIResponse{
			Status:  "error",
			Message: "Doctor diagnosis is required",
		})
		return
	}

	err = patientsHelper.UpdateDoctorDiagnosis(int(patientID), requestBody.DoctorDiagnosis)
	if err != nil {
		utils.WriteJson(w, http.StatusInternalServerError, utils.APIResponse{
			Status:  "error",
			Message: "Failed to update doctor diagnosis",
			Error:   err.Error(),
		})
		return
	}

	utils.WriteJson(w, http.StatusOK, utils.APIResponse{
		Status:  "success",
		Message: "Doctor diagnosis updated successfully",
	})
}
