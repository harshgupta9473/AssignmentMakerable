package middlewares

import (
	"net/http"

	authHelper "github.com/harshgupta9473/assignment_makerable/helpers/auth"
	"github.com/harshgupta9473/assignment_makerable/utils"
)

func IsDoctor(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, err := GetUserContext(r)
		if err != nil {
			utils.WriteJson(w, http.StatusUnauthorized, utils.APIResponse{
				Status:  "error",
				Message: "Unauthorized: User context not found",
				Error:   err.Error(),
			})
			return
		}

		if user.Role != utils.DocotorRole {
			utils.WriteJson(w, http.StatusForbidden, utils.APIResponse{
				Status:  "error",
				Message: "Access denied: Only doctors are allowed to access this route.",
			})
			return
		}

		next.ServeHTTP(w, r)
	})
}

func IsDoctorApproved(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get user context
		user, err := GetUserContext(r)
		if err != nil {
			utils.WriteJson(w, http.StatusUnauthorized, utils.APIResponse{
				Status:  "error",
				Message: "Unauthorized: User context not found",
				Error:   err.Error(),
			})
			return
		}
		if user.IsApproved{
			next.ServeHTTP(w,r)
			return
		}


		isApproved, err := authHelper.IsUserApproved(user.UserID)
		if err != nil {
			utils.WriteJson(w, http.StatusInternalServerError, utils.APIResponse{
				Status:  "error",
				Message: "Could not check doctor's approval status",
				Error:   err.Error(),
			})
			return
		}

		if !isApproved {
			utils.WriteJson(w, http.StatusForbidden, utils.APIResponse{
				Status:  "error",
				Message: "Doctor is not approved. Please contact admin for approval.",
			})
			return
		}

		next.ServeHTTP(w, r)
	})
}
