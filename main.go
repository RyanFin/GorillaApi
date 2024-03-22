package main

import (
	"RyanFin/GorillaApi/pkg/model"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
)

const (
	addr = "127.0.0.1"
	port = ":8080"
)

func main() {

	r := mux.NewRouter()
	r.HandleFunc("/cars/", getAllCarsHandler).Methods(http.MethodGet)
	r.HandleFunc("/cars/{id}", getCarByIDHandler).Methods(http.MethodGet)
	http.Handle("/", r)
	srv := &http.Server{
		Handler: r,
		Addr:    fmt.Sprintf("%s%s", addr, port),
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	fmt.Printf("listening on port %s ...", port)
	log.Fatal(srv.ListenAndServe())

}

func getAllCarsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin ", "*")

	cars, err := loadDataFromJSONFile()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// send write code to response with encoder.Encode(obj)
	encoder := json.NewEncoder(w)
	err = encoder.Encode(cars)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func getCarByIDHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	// search for the car with the id passed in
	cars, err := loadDataFromJSONFile()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for _, car := range cars {
		if car.UUID == id {
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("Access-Control-Allow-Origin ", "*")
			json.NewEncoder(w).Encode(car)
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
	fmt.Fprintf(w, "Car with ID %s not found", id)
}

func loadDataFromJSONFile() ([]model.Car, error) {
	//get data from file
	f, err := os.Open("cars.json")
	if err != nil {
		return nil, err
	}

	defer f.Close()
	// Unmarshal JSON data into car slice
	var cars []model.Car
	decoder := json.NewDecoder(f)
	decoder.Decode(&cars)
	if err != nil {
		return nil, err
	}

	return cars, nil
}

func simpleMw(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do stuff here
		log.Println(r.RequestURI)
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}
