package doctor

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	doctorHelper "github.com/harshgupta9473/assignment_makerable/helpers/doctor"
	"github.com/harshgupta9473/assignment_makerable/utils"
)

func GetAllApprovedDoctorsDetails(w http.ResponseWriter, r *http.Request) {
	doctors, err := doctorHelper.GetAllApprovedAndVerifiedDoctorsDetailByRole()
	if err != nil {
		utils.WriteJson(w, http.StatusInternalServerError, utils.APIResponse{
			Status:  "error",
			Message: "Internal Server Error",
			Error:   "DB_ERR",
		})
		return
	}

	if len(doctors) == 0 {
		utils.WriteJson(w, http.StatusNoContent, utils.APIResponse{
			Status:  "error",
			Message: "No Doctors Found",
			Error:   "NO_CONTENT",
		})
		return
	}

	utils.WriteJson(w, http.StatusOK, utils.APIResponse{
		Status:  "success",
		Message: "Doctors Fetched Successfully",
		Data:    doctors,
	})
}


func GetDoctorById(w http.ResponseWriter, r *http.Request) {
	
	doctorIdStr := mux.Vars(r)["doctorID"]
	doctorID, err := strconv.ParseInt(doctorIdStr, 10, 64)
	if err != nil {
		utils.WriteJson(w, http.StatusBadRequest, utils.APIResponse{
			Status:  "error",
			Message: "Invalid doctor ID",
			Error:   "INVALID_ID",
		})
		return
	}

	
	doc, err := doctorHelper.GetADoctorById(doctorID)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.WriteJson(w, http.StatusNotFound, utils.APIResponse{
				Status:  "error",
				Message: "Doctor not found or not approved/verified",
				Error:   "NOT_FOUND",
			})
			return
		}
		utils.WriteJson(w, http.StatusInternalServerError, utils.APIResponse{
			Status:  "error",
			Message: "Failed to fetch doctor",
			Error:   err.Error(),
		})
		return
	}

	utils.WriteJson(w, http.StatusOK, utils.APIResponse{
		Status:  "success",
		Message: "Doctor fetched successfully",
		Data:    doc,
	})
}
