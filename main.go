package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/kelseyhightower/envconfig"
	"github.com/renu-ramesh/robot-apocalypse-docker/handlers"
	"github.com/renu-ramesh/robot-apocalypse-docker/models"
)

func main() {
	fmt.Println("smbteam4 - Rest API v1.0")
	handleRequests()
}

func handleRequests() {

	//create Handlers object
	h := handlers.New()

	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/api/v1/add-survivors", h.CreateNewsurvivors).Methods("POST")
	myRouter.HandleFunc("/api/v1/update-location", h.UpdateSurvivorsLocation).Methods("POST")
	myRouter.HandleFunc("/api/v1/infected", h.UpdateInfection).Methods("POST")

	myRouter.HandleFunc("/api/v1/percentage/{spec}", h.PercentageSpecification)
	myRouter.HandleFunc("/api/v1/survivors/{spec}", h.ListSurvivors)
	myRouter.HandleFunc("/api/v1/robots/all", h.ListAllRobots)

	//Initailize Environment Variables
	var env models.EnvVariables
	err := envconfig.Process("robot_apocalypse", &env)
	if err != nil {
		log.Fatal("unable to initialize environment variables", err.Error())
	}
	log.Fatal(http.ListenAndServe(":"+env.Port, myRouter))
}
