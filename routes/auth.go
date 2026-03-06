package routes

import (
	"net/http"
	"fmt"
	"harbor/main/db"
	"log"
	"database/sql"
)

func signUp(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	var username string = r.FormValue("username")
	var email string = r.FormValue("email")
	var pass string = r.FormValue("pass")
	_, err := db.FindUserByEmail(email)
	if(err != sql.ErrNoRows){
		fmt.Fprintln(w, "Email Already Exists")
		return
	}
	usr := db.User{
		Username: username,
		Email: email,
		Password: pass,
	}
	err = db.CreateUser(usr)
	if(err != nil) {
		log.Fatal(err)
		return
	}
	fmt.Fprintln(w, "Success!")
	
}

func logIn(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	var email string = r.FormValue("email")
	var pass string = r.FormValue("pass")
	u, err := db.FindUserByEmail(email)
	if(err != nil){
		if(err == sql.ErrNoRows){
			fmt.Fprintln(w, "No Such Account Exists")
		} else{
			fmt.Fprintln(w, "Internal Server Error")
		}
		return
	}
	if(pass == u.Password){
		fmt.Fprintln(w, "Logged In Successfully!")
	} else{
		fmt.Fprintln(w, "Failed To Login")	
	}
	
}

func SetupAuthRoutes(){
	http.HandleFunc("POST /signup", signUp)
	http.HandleFunc("POST /login", logIn)

}