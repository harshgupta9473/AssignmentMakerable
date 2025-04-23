package routes

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/harshgupta9473/assignment_makerable/handlers/admin"
	"github.com/harshgupta9473/assignment_makerable/handlers/auth"
	"github.com/harshgupta9473/assignment_makerable/handlers/doctor"
	patient "github.com/harshgupta9473/assignment_makerable/handlers/patients.go"
	"github.com/harshgupta9473/assignment_makerable/middlewares"
)

func RegisterRoutes(r *mux.Router) {
	// signups
	r.HandleFunc("/auth/doctor/signup", auth.DoctorSignUp).Methods("POST")                     //#apichecked
	r.HandleFunc("/auth/receptionist/signup", auth.ReceptionistSignUp).Methods("POST")         //#apichecked
	r.HandleFunc("/auth/admin/signup", admin.SignupAdmin).Methods("POST")

	//login and email verification
	r.HandleFunc("/auth/login", auth.LoginHandler).Methods("POST")   //#apichecked
	r.HandleFunc("/auth/verify", auth.VerifyEmail).Methods("GET")    //#apichecked
	r.Handle("/auth/resendLink", middlewares.AuthMiddleware(http.HandlerFunc(auth.SendVerificationLink))).Methods("GET")   //#apichecked 

	//admin
	//http://localhost:3001/admin/doctors?approved=true
	r.Handle("/admin/doctors", middlewares.AuthMiddleware(middlewares.IsEmailVerified(middlewares.IsAdmin(http.HandlerFunc(admin.GetAllDoctorsHandler))))).Methods("GET")
	r.Handle("/admin/doctors/approve/{doctorID}", middlewares.AuthMiddleware(middlewares.IsEmailVerified(middlewares.IsAdmin(http.HandlerFunc(admin.ApproveDoctorHandler))))).Methods("GET")
	r.Handle("/admin/receptionists", middlewares.AuthMiddleware(middlewares.IsEmailVerified(middlewares.IsAdmin(http.HandlerFunc(admin.GetAllReceptionistsHandler))))).Methods("GET")
	r.Handle("/admin/receptionists/approve/{receptionistID}", middlewares.AuthMiddleware(middlewares.IsEmailVerified(middlewares.IsAdmin(http.HandlerFunc(admin.ApproveReceptionistHandler))))).Methods("GET")

	//doctor
	// all patient under that doctor
	r.Handle("/doctor/patients", middlewares.AuthMiddleware(middlewares.IsEmailVerified(middlewares.IsDoctor(middlewares.IsDoctorApproved(http.HandlerFunc(patient.GetAllPatientsByDoctors)))))).Methods("GET")
	// patient under that doctor by patientID
	r.Handle("/doctor/patients/{patientID}",middlewares.AuthMiddleware(middlewares.IsEmailVerified(middlewares.IsDoctor(middlewares.IsDoctorApproved(http.HandlerFunc(patient.GetAPatientPatientID)))))).Methods("GET")
	//update patient Diagnosis BY Doctor
	r.Handle("/doctor/patients/comments/{patientID}",middlewares.AuthMiddleware(middlewares.IsEmailVerified(middlewares.IsDoctor(middlewares.IsDoctorApproved(http.HandlerFunc(patient.UpdateDoctorDiagnosisHandler)))))).Methods("PUT")

	// receptionist
            	//patient
	r.Handle("/receptionist/patients/register", middlewares.AuthMiddleware(middlewares.IsEmailVerified(middlewares.IsReceptionist(middlewares.IsReceptionistApproved(http.HandlerFunc(patient.RegisterPatient)))))).Methods("POST")
	r.Handle("/receptionist/patients/", middlewares.AuthMiddleware(middlewares.IsEmailVerified(middlewares.IsReceptionist(middlewares.IsReceptionistApproved(http.HandlerFunc(patient.GetAllPatients)))))).Methods("GET")
	r.Handle("/receptionist/patients/{patientID}", middlewares.AuthMiddleware(middlewares.IsEmailVerified(middlewares.IsReceptionist(middlewares.IsReceptionistApproved(http.HandlerFunc(patient.GetPatientByID)))))).Methods("GET")
	// r.Handle("/receptionist/patients/{doctoID}/{patientID}", middlewares.AuthMiddleware(middlewares.IsEmailVerified(middlewares.IsReceptionist(middlewares.IsReceptionistApproved(http.HandlerFunc(patient.GetAllPatientsByDoctorID)))))).Methods("GET")
	r.Handle("/receptionist/patients/{patientID}", middlewares.AuthMiddleware(middlewares.IsEmailVerified(middlewares.IsReceptionist(middlewares.IsReceptionistApproved(http.HandlerFunc(patient.UpdatePatientHandler)))))).Methods("PUT")
	r.Handle("/receptionist/patients/{patientID}", middlewares.AuthMiddleware(middlewares.IsEmailVerified(middlewares.IsReceptionist(middlewares.IsReceptionistApproved(http.HandlerFunc(patient.DeletePatientById)))))).Methods("DELETE")
	//doctor
	r.Handle("/receptionist/doctors", middlewares.AuthMiddleware(middlewares.IsEmailVerified(middlewares.IsReceptionist(middlewares.IsReceptionistApproved(http.HandlerFunc(doctor.GetAllApprovedDoctorsDetails)))))).Methods("GET")
	r.Handle("/receptionist/doctors/{doctorID}", middlewares.AuthMiddleware(middlewares.IsEmailVerified(middlewares.IsReceptionist(middlewares.IsReceptionistApproved(http.HandlerFunc(doctor.GetDoctorById)))))).Methods("GET")
	
	
	
}
