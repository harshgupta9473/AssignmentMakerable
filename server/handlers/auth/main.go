package auth

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	authHelper "github.com/harshgupta9473/assignment_makerable/helpers/auth"
	"github.com/harshgupta9473/assignment_makerable/middlewares"
	"github.com/harshgupta9473/assignment_makerable/models"
	"github.com/harshgupta9473/assignment_makerable/utils"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func ReceptionistSignUp(w http.ResponseWriter,r *http.Request){
	var req models.SignupRequestReceptionist
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteJson(w, http.StatusBadRequest, utils.APIResponse{
			Status:  "fail",
			Message: "Invalid request body",
			Error:   err.Error(),
		})
		return
	}

	exists, err := authHelper.ISUserExistsByEmail(req.Email)
	if err != nil {
		utils.WriteJson(w, http.StatusInternalServerError, utils.APIResponse{
			Status:  "error",
			Message: "Failed to check user existence",
			Error:   err.Error(),
		})
		return
	}
	if exists {
		utils.WriteJson(w, http.StatusConflict, utils.APIResponse{
			Status:  "fail",
			Message: "User with this email already exists",
		})
		return
	}
	hashedPassword, err := HashPassword(req.Password)
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
	expiryTime := time.Now().Add(24 * time.Hour)

	req.Password=hashedPassword
	userId, err := authHelper.InsertIntoReceptionist(req,token,expiryTime)
	if err != nil {
		utils.WriteJson(w, http.StatusInternalServerError, utils.APIResponse{
			Status:  "error",
			Message: "Failed to register Receptionist",
			Error:   err.Error(),
		})
		return
	}

	// Attempt to send verification email (even if it fails, return partial success)
	err = utils.SendVerificationEmail(req.Email, token, userId)
	if err != nil {
		utils.WriteJson(w, http.StatusOK, utils.APIResponse{
			Status:  "success",
			Message: "Receptionist registered successfully, but failed to send verification email. Please log in and request verification again.",
		})
		return
	}

	utils.WriteJson(w, http.StatusCreated, utils.APIResponse{
		Status:  "success",
		Message: "Receptionist registered successfully. Please check your email to verify your account.",
	})

}

func DoctorSignUp(w http.ResponseWriter, r *http.Request) {
	var req models.SignupRequestDoctor
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteJson(w, http.StatusBadRequest, utils.APIResponse{
			Status:  "fail",
			Message: "Invalid request body",
			Error:   err.Error(),
		})
		return
	}

	// Check if user already exists
	exists, err := authHelper.ISUserExistsByEmail(req.Email)
	if err != nil {
		utils.WriteJson(w, http.StatusInternalServerError, utils.APIResponse{
			Status:  "error",
			Message: "Failed to check user existence",
			Error:   err.Error(),
		})
		return
	}
	if exists {
		utils.WriteJson(w, http.StatusConflict, utils.APIResponse{
			Status:  "fail",
			Message: "User with this email already exists",
		})
		return
	}

	// Hash the password
	hashedPassword, err := HashPassword(req.Password)
	if err != nil {
		utils.WriteJson(w, http.StatusInternalServerError, utils.APIResponse{
			Status:  "error",
			Message: "Failed to hash password",
			Error:   err.Error(),
		})
		return
	}

	// Generate verification token
	token, err := utils.GenerateOTPToken()
	if err != nil {
		utils.WriteJson(w, http.StatusInternalServerError, utils.APIResponse{
			Status:  "error",
			Message: "Failed to generate verification token",
			Error:   err.Error(),
		})
		return
	}
	expiryTime := time.Now().Add(24 * time.Hour)

	// Insert into users and gov ID table with transaction
	req.Password=hashedPassword
	userId, err := authHelper.InsertDoctorWithGovID(req, token, expiryTime)
	if err != nil {
		utils.WriteJson(w, http.StatusInternalServerError, utils.APIResponse{
			Status:  "error",
			Message: "Failed to register doctor",
			Error:   err.Error(),
		})
		return
	}

	// Attempt to send verification email (even if it fails, return partial success)
	err = utils.SendVerificationEmail(req.Email, token, userId)
	if err != nil {
		utils.WriteJson(w, http.StatusOK, utils.APIResponse{
			Status:  "success",
			Message: "Doctor registered successfully, but failed to send verification email. Please log in and request verification again.",
		})
		return
	}

	utils.WriteJson(w, http.StatusCreated, utils.APIResponse{
		Status:  "success",
		Message: "Doctor registered successfully. Please check your email to verify your account.",
	})
}

func VerifyEmail(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	token := r.URL.Query().Get("token")
	userIDStr := r.URL.Query().Get("userID")

	if token == "" || userIDStr == "" {
		utils.WriteJson(w, http.StatusBadRequest, utils.APIResponse{
			Status:  "fail",
			Message: "Token and userID are required",
		})
		return
	}

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		utils.WriteJson(w, http.StatusBadRequest, utils.APIResponse{
			Status:  "fail",
			Message: "Invalid userID",
			Error:   err.Error(),
		})
		return
	}

	storedToken, expiry, isVerified, err := authHelper.GetTokenExpiryIsVerified(userID)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.WriteJson(w, http.StatusNotFound, utils.APIResponse{
				Status:  "fail",
				Message: "User not found",
			})
			return
		}
		utils.WriteJson(w, http.StatusInternalServerError, utils.APIResponse{
			Status:  "error",
			Message: "Database error",
			Error:   err.Error(),
		})
		return
	}

	if isVerified {
		utils.WriteJson(w, http.StatusOK, utils.APIResponse{
			Status:  "success",
			Message: "Email is already verified",
		})
		return
	}

	if token != storedToken {
		utils.WriteJson(w, http.StatusUnauthorized, utils.APIResponse{
			Status:  "fail",
			Message: "Invalid verification token",
		})
		return
	}

	if time.Now().After(expiry) {
		utils.WriteJson(w, http.StatusUnauthorized, utils.APIResponse{
			Status:  "fail",
			Message: "Verification token has expired",
		})
		return
	}

	err = authHelper.VerifyUserwithUserID(userID)
	if err != nil {
		utils.WriteJson(w, http.StatusInternalServerError, utils.APIResponse{
			Status:  "error",
			Message: "Failed to update verification status",
			Error:   err.Error(),
		})
		return
	}

	utils.WriteJson(w, http.StatusOK, utils.APIResponse{
		Status:  "success",
		Message: "Email verified successfully",
	})
}

func SendVerificationLink(w http.ResponseWriter, r *http.Request) {

	claims, err:=middlewares.GetUserContext(r)
	if err!=nil {
		utils.WriteJson(w, http.StatusUnauthorized, utils.APIResponse{
			Status:  "error",
			Message: "Unauthorized",
			Error:   "Invalid token claims",
		})
		return
	}

	if claims.IsVerified {
		utils.WriteJson(w, http.StatusBadRequest, utils.APIResponse{
			Status:  "error",
			Message: "Already verified",
			Error:   "User already verified",
		})
		return
	}

	userID := claims.UserID

	user, err := authHelper.GetUserByUserId(userID)
	if err != nil {
		utils.WriteJson(w, http.StatusInternalServerError, utils.APIResponse{
			Status:  "error",
			Message: "User fetch failed",
			Error:   err.Error(),
		})
		return
	}
	if user.IsVerified{
		utils.WriteJson(w, http.StatusOK, utils.APIResponse{
			Status:  "error",
			Message: "User Already Verified",
		})
		return
	}

	token, err := utils.GenerateOTPToken()
	if err != nil {
		utils.WriteJson(w, http.StatusInternalServerError, utils.APIResponse{
			Status:  "error",
			Message: "Token generation failed",
			Error:   err.Error(),
		})
		return
	}

	expiry := time.Now().Add(30 * time.Minute)

	err = authHelper.UpdateUserVerificationToken(userID, token, expiry)
	if err != nil {
		utils.WriteJson(w, http.StatusInternalServerError, utils.APIResponse{
			Status:  "error",
			Message: "Failed to update verification token",
			Error:   err.Error(),
		})
		return
	}

	err = utils.SendVerificationEmail(user.Email, token, userID)
	if err != nil {
		utils.WriteJson(w, http.StatusInternalServerError, utils.APIResponse{
			Status:  "partial",
			Message: "Token saved, but email sending failed",
			Error:   err.Error(),
		})
		return
	}

	utils.WriteJson(w, http.StatusOK, utils.APIResponse{
		Status:  "success",
		Message: "Verification email sent",
	})
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteJson(w, http.StatusBadRequest, utils.APIResponse{
			Status:  "error",
			Message: "Invalid request",
			Error:   err.Error(),
		})
		return
	}

	
	user, err := authHelper.GetUserByUserByEmailAndRole(req.Email, req.Role)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.WriteJson(w, http.StatusBadRequest, utils.APIResponse{
				Status:  "error",
				Message: "Invalid Email or Password",
				Error:   "EMAIL_NOT_FOUND",
			})
		} else {
			utils.WriteJson(w, http.StatusInternalServerError, utils.APIResponse{
				Status:  "error",
				Message: "Internal Server Error",
				Error:   err.Error(),
			})
		}
		return
	}

	// Check password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		utils.WriteJson(w, http.StatusBadRequest, utils.APIResponse{
			Status:  "error",
			Message: "Invalid Email or Password",
			Error:   "INVALID_CREDENTIALS",
		})
		return
	}

	

	// JWT generation
	token, err := middlewares.GenerateJWT(user.ID, user.Email, user.Role, user.IsVerified, false)
	if err != nil {
		utils.WriteJson(w, http.StatusInternalServerError, utils.APIResponse{
			Status:  "error",
			Message: "JWT generation failed",
			Error:   err.Error(),
		})
		return
	}

	
	utils.WriteJson(w, http.StatusOK, utils.APIResponse{
		Status:  "success",
		Message: "Login successful",
		Data: map[string]any{
			"token": token,
			"user": map[string]any{
				"id":         user.ID,
				"email":      user.Email,
				"first_name": user.FirstName,
				"last_name":  user.LastName,
				"role":       user.Role,
				"is_approved":user.IsApproved,
			},
		},
	})
}
