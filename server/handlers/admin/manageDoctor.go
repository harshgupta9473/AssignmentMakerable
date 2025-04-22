package admin

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	adminHelper "github.com/harshgupta9473/assignment_makerable/helpers/admin"
	"github.com/harshgupta9473/assignment_makerable/utils"
)

func GetAllDoctorsHandler(w http.ResponseWriter, r *http.Request) {
	approved := r.URL.Query().Get("approved")

	doctors, err := adminHelper.GetUsersByApproval(approved, utils.DocotorRole)
	if err != nil {
		utils.WriteJson(w, http.StatusInternalServerError, utils.APIResponse{
			Status:  "error",
			Message: "Failed to fetch doctors",
			Error:   err.Error(),
		})
		return
	}

	utils.WriteJson(w, http.StatusOK, utils.APIResponse{
		Status:  "success",
		Message: "Doctors fetched successfully",
		Data:    doctors,
	})
}

func ApproveDoctorHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	doctorID, err := strconv.ParseInt(vars["doctorID"], 10, 64)
	if err != nil {
		utils.WriteJson(w, http.StatusBadRequest, utils.APIResponse{
			Status:  "error",
			Message: "Invalid doctor ID",
			Error:   err.Error(),
		})
		return
	}

	err = adminHelper.ApproveUserByRole(doctorID, utils.DocotorRole)
	if err != nil {
		utils.WriteJson(w, http.StatusInternalServerError, utils.APIResponse{
			Status:  "error",
			Message: "Failed to approve doctor",
			Error:   err.Error(),
		})
		return
	}

	utils.WriteJson(w, http.StatusOK, utils.APIResponse{
		Status:  "success",
		Message: "Doctor approved successfully",
	})
}

func GetAllReceptionistsHandler(w http.ResponseWriter, r *http.Request) {
	approved := r.URL.Query().Get("approved")

	receptionists, err := adminHelper.GetUsersByApproval(approved, utils.ReceptionistRole)
	if err != nil {
		utils.WriteJson(w, http.StatusInternalServerError, utils.APIResponse{
			Status:  "error",
			Message: "Failed to fetch receptionists",
			Error:   err.Error(),
		})
		return
	}

	utils.WriteJson(w, http.StatusOK, utils.APIResponse{
		Status:  "success",
		Message: "Receptionists fetched successfully",
		Data:    receptionists,
	})
}

func ApproveReceptionistHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	receptionistID, err := strconv.ParseInt(vars["receptionistID"], 10, 64)
	if err != nil {
		utils.WriteJson(w, http.StatusBadRequest, utils.APIResponse{
			Status:  "error",
			Message: "Invalid receptionist ID",
			Error:   err.Error(),
		})
		return
	}

	err = adminHelper.ApproveUserByRole(receptionistID, utils.ReceptionistRole)
	if err != nil {
		utils.WriteJson(w, http.StatusInternalServerError, utils.APIResponse{
			Status:  "error",
			Message: "Failed to approve receptionist",
			Error:   err.Error(),
		})
		return
	}

	utils.WriteJson(w, http.StatusOK, utils.APIResponse{
		Status:  "success",
		Message: "Receptionist approved successfully",
	})
}
