package main

import (
	"BMTimetable/web/app/controllers"
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func main() {

	router := mux.NewRouter()
	router.HandleFunc("/api/helloworld", controllers.HelloWorld)

	/*
		/api/arrets => liste des arrêts
		/api/arret/xxx/lignes => liste des lignes à cet arret
		/api/arret/xxx/lignes/A/sens => sens possible pour une ligne à un arrêt
		/api/arret/154/ligne/A/sens/A/passages?next=3

		/api/passages?
	*/

	port := os.Getenv("PORT") //Get port from .env file, we did not specify any port so this should return an empty string when tested locally
	if port == "" {
		port = "8000" //localhost
	}

	fmt.Println(port)

	err := http.ListenAndServe(":"+port, router) //Launch the app, visit localhost:8000/api
	if err != nil {
		fmt.Print(err)
	}
}
