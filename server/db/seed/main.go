package seed

import (
	"log"
	"os"

	"github.com/harshgupta9473/assignment_makerable/db"
)

func SeedAdminEmail() {
	if(os.Getenv("seed")=="true"){
		email:=os.Getenv("seedEmail")
	      query:=`INSERT INTO admins(email)
		  VALUES ($1)`
		  _,err:=db.GetDB().Exec(query,email)
		  if err!=nil{
			log.Println("Error seeding the admin mail")
			return
		  }
		  log.Println("email seeded")
	}
	return
}