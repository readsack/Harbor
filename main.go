package main

import (
	"log"
	"net/http"
	"harbor/main/db"
	"harbor/main/routes"
)

const PORT string = ":8080"

func main(){
	//u := db.User{Username: "tanisha", Email: "pagal@gmail.com", Password: "pagalhumain"}


	log.Println("Server starting At localhost"+PORT)
	db.InitDB()
	db.SetupDB()
	routes.SetupAuthRoutes()
	//log.Println(db.AppDB.QueryRow())
	log.Fatal(http.ListenAndServe(PORT, nil))
	
}


