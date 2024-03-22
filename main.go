package main

import (
	"RyanFin/GorillaApi/pkg/model"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
)

const (
	addr = "127.0.0.1"
	port = ":8080"
)

var jwtKey = []byte("my_secret_key")

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func main() {

	r := mux.NewRouter()

	r.HandleFunc("/token", generateTokenHandler).Methods(http.MethodGet)
	r.HandleFunc("/cars", authenticate(getAllCarsHandler)).Methods(http.MethodGet)
	r.HandleFunc("/cars/{id}", authenticate(getCarByIDHandler)).Methods(http.MethodGet)
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

func generateTokenHandler(w http.ResponseWriter, r *http.Request) {
	expirationTime := time.Now().Add(5 * time.Minute)
	claims := &Claims{
		Username: "user",
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error generating token: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
}

// Authentication middleware
func authenticate(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/token" {
			next.ServeHTTP(w, r)
			return
		}

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprintf(w, "Unauthorized")
			return
		}

		tokenString := authHeader[len("Bearer "):]

		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})
		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				w.WriteHeader(http.StatusUnauthorized)
				fmt.Fprintf(w, "Unauthorized")
				return
			}
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Bad Request")
			return
		}

		if !token.Valid {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprintf(w, "Unauthorized")
			return
		}

		next.ServeHTTP(w, r)
	})
}
