package middlewares

import (
	"net/http"

	authHelper "github.com/harshgupta9473/assignment_makerable/helpers/auth"
	"github.com/harshgupta9473/assignment_makerable/utils"
)

func IsEmailVerified(next http.Handler) http.Handler {
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

		if !user.IsVerified {
			
			updatedUser, err := authHelper.GetUserByUserId(user.UserID)
			if err != nil {
				utils.WriteJson(w, http.StatusInternalServerError, utils.APIResponse{
					Status:  "error",
					Message: "Could not verify email status from database",
					Error:   err.Error(),
				})
				return
			}

			if updatedUser.IsVerified {
				newToken, err := GenerateJWT(updatedUser.ID, updatedUser.Email,updatedUser.Role, true,updatedUser.IsApproved)
				if err != nil {
					utils.WriteJson(w, http.StatusInternalServerError, utils.APIResponse{
						Status:  "error",
						Message: "Failed to generate updated token",
						Error:   err.Error(),
					})
					return
				}

				w.Header().Set("newtoken", newToken)

				next.ServeHTTP(w,r)
				return
			}

			
			utils.WriteJson(w, http.StatusForbidden, utils.APIResponse{
				Status:  "fail",
				Message: "Please verify your email to continue.",
			})
			return
		}

		next.ServeHTTP(w, r)
	})
}
