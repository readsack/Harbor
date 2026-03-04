package main

import (
	"log"
	"net/http"
	"harbor/main/db"
)

const PORT string = ":8080"

func main(){
	log.Println("Server starting At localhost"+PORT)
	db.InitDB()
	db.SetupDB()
	log.Fatal(http.ListenAndServe(PORT, nil))
	
}


