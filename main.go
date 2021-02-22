package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

// use godot package to load/read the .env file and
// return the value of the key
func goDotEnvVariable(key string) string {

	// load .env file
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}

// User data
type User struct {
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	City      string `json:"city"`
}

var user = []*User{
	{
		FirstName: "Hasbi",
		LastName:  "Qohar",
		City:      "JKT",
	},
	{
		FirstName: "Hadi",
		LastName:  "Mustofa",
		City:      "MDN",
	},
	{
		FirstName: "Haqi",
		LastName:  "Muttaqin",
		City:      "MDN",
	},
}

// ErrorUserNotFound is
func ErrorUserNotFound(name string) error {
	return fmt.Errorf("User %v not found", name)
}

func findUser(name string) (int, error) {
	for i, user := range user {
		dataFirstName := strings.ToLower(user.FirstName)
		dataLastName := strings.ToLower(user.LastName)
		input := strings.ToLower(name)
		if dataFirstName == input || dataLastName == input {
			return i, nil
		}
	}

	return -1, ErrorUserNotFound(name)
}

// HomeHandler for routing to homepage
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Bismillah Home")
}

// GetAllUserHandler (w http.ResponseWriter, r *http.Request)
func GetAllUserHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

// UserGetHandler for routing to homepage
func UserGetHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	i, err := findUser(vars["username"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user[i])
}

// ErrorUserByCityNotFound (cityName string) error {
func ErrorUserByCityNotFound(cityName string) error {
	return fmt.Errorf("User from %v city not found", cityName)
}

func findUserByCity(c string) ([]*User, error) {

	var userByCity []*User
	for _, data := range user {
		cityCode := strings.ToLower(data.City)
		input := strings.ToLower(c)
		if cityCode == input {
			userByCity = append(userByCity, data)
		} else {
			switch cityName := c; cityName {
			case "madiun":
				if cityCode == "mdn" {
					userByCity = append(userByCity, data)
				}
			case "jakarta":
				if cityCode == "jkt" {
					userByCity = append(userByCity, data)
				}
			}
		}
	}
	if len(userByCity) == 0 {
		return nil, ErrorUserByCityNotFound(c)
	}
	return userByCity, nil
}

// UserByCityGetHandler apasi
func UserByCityGetHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userByCity, err := findUserByCity(vars["city"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(userByCity)
}

func main() {
	// port := goDotEnvVariable("PORT")
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	var dir string

	flag.StringVar(&dir, "dir", ".", "the directory to serve files from. Defaults to the current dir")
	flag.Parse()

	r := mux.NewRouter()
	r.HandleFunc("/", HomeHandler)
	r.HandleFunc("/user/{username}", UserGetHandler).Methods("GET")
	r.HandleFunc("/city/{city}", UserByCityGetHandler).Methods("GET")
	r.HandleFunc("/user", GetAllUserHandler).Methods("GET")

	// This will serve files under http://localhost:8000/portfolio/<filename>
	r.PathPrefix("/portfolio/").Handler(http.StripPrefix("", http.FileServer(http.Dir(dir))))

	// This will serve files under http://localhost:8000/portfolio/<filename>
	r.PathPrefix("/static/").Handler(http.StripPrefix("", http.FileServer(http.Dir(dir))))

	s := &http.Server{
		Handler: r,
		Addr:    ":" + port,
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	fmt.Println("Starting server on PORT:" + port)
	go func() {
		err := s.ListenAndServe()
		if err != nil {
			log.Fatal(err)
		}
	}()

	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)

	sig := <-sigChan
	fmt.Println("Yaah, server nya lagi down", sig)
}
