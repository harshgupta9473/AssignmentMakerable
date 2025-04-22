package middlewares

import (
	"net/http"

	"github.com/harshgupta9473/assignment_makerable/utils"
)

func IsAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, err := GetUserContext(r)
		if err != nil {
			utils.WriteJson(w, http.StatusUnauthorized, utils.APIResponse{
				Status:  "fail",
				Message: "Unauthorized: User not found in context",
				Error:   err.Error(),
			})
			return
		}

		if user.Role != utils.AdminRole {
			utils.WriteJson(w, http.StatusForbidden, utils.APIResponse{
				Status:  "fail",
				Message: "Forbidden: Admin access required",
				Error:   "INSUFFICIENT_PERMISSIONS",
			})
			return
		}

		next.ServeHTTP(w, r)
	})
}
