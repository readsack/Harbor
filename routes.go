package main

import (
    "net/http"
	"log"
)

func routeSignUp(w http.ResponseWriter, r *http.Request){
	r.ParseForm()
	var user string = r.FormValue("username")
	var passwd string = r.FormValue("password")
	var email string = r.FormValue("email")
	// 1. check validity here
	// 2. create account if valid details
	// 3. set jwt cookies
	log.Printf("%s  %s  %s", user, email, passwd)
	http.Redirect(w, r, "/", 301)
}

func routeLogIn(w http.ResponseWriter, r *http.Request) {
	r.ParseForm();
	var user string = r.FormValue("username")
	var passwd string = r.FormValue("password")
	// 1. check validity here
	// 2. set jwt cookies

	log.Printf("%s %s", user, passwd)
	http.Redirect(w, r, "/", 301)
}


func SetupPostRoutes(){
	http.HandleFunc("POST /signup", routeSignUp)
	http.HandleFunc("POST /login", routeLogIn)
}
