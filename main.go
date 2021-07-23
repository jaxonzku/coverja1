package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	gorm.Model
	Name     string `json:"name"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func signUp(w http.ResponseWriter, r *http.Request) {
	fmt.Println("signup hit")
	//fmt.Fprintf(w, "signup hit")
	db := dbconnect()

	var p User
	var password string

	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	password = p.Password
	hash, _ := HashPassword(password)

	sqlStatement := `INSERT INTO emp (name, phone,email,password) VALUES ($1,$2,$3,$4)`
	_, err = db.Exec(sqlStatement, p.Name, p.Phone, p.Email, hash)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println("error in here")
		panic(err)
	} else {
		fmt.Println("\nRow inserted successfully!")
	}

	w.WriteHeader(http.StatusOK)

}

func signIn(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "signin hit")
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Println("in homepage")
	fmt.Fprintf(w, "homepage hit")
}

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", homePage).Methods("GET")
	myRouter.HandleFunc("/api/signUp", signUp).Methods("POST")
	myRouter.HandleFunc("/api/signIn", signIn).Methods("GET")

	log.Fatal(http.ListenAndServe(":8081", myRouter))
}

func dbconnect() *sql.DB {
	connStr := "user=postgres dbname=blogapp password=root123 host=localhost sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Printf("\nSuccessfully connected to database!\n")
	//fmt.Println(db)
	return db

}

func main() {

	handleRequests()

}