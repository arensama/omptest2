package main

import (
	"log"
	"net/http"

	"github.com/arensama/testapi2/src/db"
	"github.com/arensama/testapi2/src/image"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

var DB = db.Init()
var imageService = image.ServiceInit(DB)
var imageController = image.Init(imageService)

func main() {

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	// Create a new router
	router := mux.NewRouter()

	router.PathPrefix("/image").Handler(imageController)

	// Start the server
	log.Fatal(http.ListenAndServe(":4444", router))
}
