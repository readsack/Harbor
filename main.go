package main

import (
	"log"
	"net/http"
	"harbor/main/db"
	"harbor/main/routes"
	"github.com/joho/godotenv"
)

const PORT string = ":8080"

func main(){
	//u := db.User{Username: "tanisha", Email: "pagal@gmail.com", Password: "pagalhumain"}

	err := godotenv.Load()
	if(err != nil) {log.Fatal("Couldn't Load .env")}

	log.Println("Server starting At localhost"+PORT)
	db.InitDB()
	db.SetupDB()
	routes.SetupAuthRoutes()
	//log.Println(db.AppDB.QueryRow())
	log.Fatal(http.ListenAndServe(PORT, nil))
	
}


