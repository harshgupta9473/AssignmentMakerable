This is made in very short time as there are practicals and quizes and endsemd in last april.

Here there are three roles for authentication : admin,doctor,receptionist.
Not anyone can be admin . Admin email needs to be preseeded into admin table into db. Only those people with that email can create their account and need to verify their email via the verification link sent to the mail.
The doctors and receptionists can create their account. And need to verify their email.
They can't do anything until admin verified and approves them.
A receptionist can register a patient, read patients data, update patients and delete them(soft delete), can fetch all the doctors, doctors by uerid, fetch all patients.  
A doctor can fetch all the patients associated with him or her i.e. the one who is he or she dignosing. Can fetch a patient which is associiated with them by patient id.
Can comment or update on the patients diagnosing details.  
There are middlewares used for auth, role check, if email verified?, if approved by admin



## Documentation LINK
https://documenter.getpostman.com/view/34442065/2sB2ixjZDT


🛠️ Project Setup
## Clone the Repository

git clone https://github.com/harshgupta9473/AssignmentMakerable.git

cd  AssignmentMakerable



## Environment Setup

Set Up Environment Variables
Rename .env.example to .env:
     
     cp .env.example .env

Replace placeholder values in .env with your actual credentials.

You can generate a Gmail app password from Google Account Security if 2FA is enabled.
Add .env to your .gitignore


## Download Go Modules
go mod tidy

#RUN
go run main.go

