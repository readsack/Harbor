package main

import (
	"harbor/main/db"
	"harbor/main/routes"
	"log"
	"net/http"

	"github.com/joho/godotenv"
)

const PORT string = ":8080"

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Couldn't Load .env")
	}

	log.Println("Server starting At localhost" + PORT)
	db.InitDB()
	db.SetupDB()
	routes.SetupAuthRoutes()
	routes.SetupOrgRoutes()
	routes.SetupTeamRoutes()
	log.Fatal(http.ListenAndServe(PORT, nil))
}
