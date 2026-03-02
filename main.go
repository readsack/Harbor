package main

import (
	"net/http"
	"log"
)

const Port string = ":8080"

func main(){
	log.Println("Server running at localhost"+Port)
	SetupPostRoutes()
	http.Handle("/", http.FileServer(http.Dir("./static/")))
	log.Fatal(http.ListenAndServe(Port, nil))
}
