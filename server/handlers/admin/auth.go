package admin

import (
	"encoding/json"
	"net/http"
	"time"

	adminHelper "github.com/harshgupta9473/assignment_makerable/helpers/admin"
	authHelper "github.com/harshgupta9473/assignment_makerable/helpers/auth"
	"github.com/harshgupta9473/assignment_makerable/models"
	"github.com/harshgupta9473/assignment_makerable/utils"
	"golang.org/x/crypto/bcrypt"
)

func SignupAdmin(w http.ResponseWriter, r *http.Request) {
	var req models.SignUpRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteJson(w, http.StatusBadRequest, utils.APIResponse{
			Status:  "error",
			Message: "Invalid request body",
			Error:   err.Error(),
		})
		return
	}

	userExists, err := authHelper.ISUserExistsByEmail(req.Email)
	if err != nil {
		utils.WriteJson(w, http.StatusInternalServerError, utils.APIResponse{
			Status:  "error",
			Message: "Failed to check user existence",
			Error:   err.Error(),
		})
		return
	}
	if userExists {
		utils.WriteJson(w, http.StatusConflict, utils.APIResponse{
			Status:  "error",
			Message: "User with this email already exists",
			Error:   "EMAIL_ALREADY_REGISTERED",
		})
		return
	}

	adminApproved, err := adminHelper.EmailExistsInAdmins(req.Email)
	if err != nil {
		utils.WriteJson(w, http.StatusInternalServerError, utils.APIResponse{
			Status:  "error",
			Message: "Error checking admin approval",
			Error:   err.Error(),
		})
		return
	}
	if !adminApproved {
		utils.WriteJson(w, http.StatusForbidden, utils.APIResponse{
			Status:  "error",
			Message: "User is not allowed to create an admin account",
			Error:   "EMAIL_NOT_APPROVED",
		})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		utils.WriteJson(w, http.StatusInternalServerError, utils.APIResponse{
			Status:  "error",
			Message: "Failed to hash password",
			Error:   err.Error(),
		})
		return
	}

	token, err := utils.GenerateOTPToken()
	if err != nil {
		utils.WriteJson(w, http.StatusInternalServerError, utils.APIResponse{
			Status:  "error",
			Message: "Failed to generate verification token",
			Error:   err.Error(),
		})
		return
	}

	expiry := time.Now().Add(24 * time.Hour)
	req.Password = string(hashedPassword)

	userId, err := adminHelper.CreateAdminWithVerification(req, token, expiry)
	if err != nil {
		utils.WriteJson(w, http.StatusInternalServerError, utils.APIResponse{
			Status:  "error",
			Message: "Failed to create admin account",
			Error:   err.Error(),
		})
		return
	}

	err = utils.SendVerificationEmail(req.Email, token, userId)
	if err != nil {
		utils.WriteJson(w, http.StatusOK, utils.APIResponse{
			Status:  "success",
			Message: "Admin account created, but failed to send verification email. Please log in and request verification again.",
		})
		return
	}

	utils.WriteJson(w, http.StatusCreated, utils.APIResponse{
		Status:  "success",
		Message: "Admin registered successfully. Please check your email to verify your account before logging in.",
	})
}
