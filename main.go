package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
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

type Luser struct {
	gorm.Model

	Email    string `json:"email"`
	Password string `json:"password"`
}

func CreateToken(email string) (string, error) {
	var err error
	//Creating Access Token
	os.Setenv("ACCESS_SECRET", "jdnfksdmfksd") //this should be in an env file
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["user_id"] = userid
	atClaims["exp"] = time.Now().Add(time.Minute * 15).Unix()
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, err := at.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	if err != nil {
		return "", err
	}
	return token, nil
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func createblog(w http.ResponseWriter, r *http.Request) {

}

func signUp(w http.ResponseWriter, r *http.Request) {
	fmt.Println("signup hit")
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

	sqlStatement := `INSERT INTO users (name, phone,email,password) VALUES ($1,$2,$3,$4)`
	_, err = db.Exec(sqlStatement, p.Name, p.Phone, p.Email, hash)
	if err != nil {
		w.WriteHeader(http.StatusOK)
		fmt.Println("error in here")
		fmt.Fprintf(w, "email is  already there")
		//panic(err)
	} else {
		fmt.Println("\nRow inserted successfully!")
		fmt.Fprintf(w, "succes")

	}

}

func signIn(w http.ResponseWriter, r *http.Request) {
	fmt.Println("signin hit")
	db := dbconnect()
	var l Luser
	err := json.NewDecoder(r.Body).Decode(&l)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	result := db.QueryRow("select password from users where email=$1", l.Email)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	storedCreds := l
	err = result.Scan(&storedCreds.Password)
	if err != nil {

		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	match := CheckPasswordHash(l.Password, storedCreds.Password)
	fmt.Println("Match:   ", match)

	if err = bcrypt.CompareHashAndPassword([]byte(storedCreds.Password), []byte(l.Password)); err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Println("wrong pass")
	} else {
		fmt.Fprintf(w, "password matches")
		w.WriteHeader(http.StatusAccepted)

	}

}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Println("in homepage")
	fmt.Fprintf(w, "homepage hit")
	type Welcome struct {
		Name string
		Time string
	}

	welcome := Welcome{"Anonymous", time.Now().Format(time.Stamp)}
	templates := template.Must(template.ParseFiles("templates/home.html"))

	if err := templates.ExecuteTemplate(w, "home.html", welcome); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", homePage).Methods("GET")
	myRouter.HandleFunc("/api/signUp", signUp).Methods("POST")
	myRouter.HandleFunc("/api/signIn", signIn).Methods("POST")
	myRouter.HandleFunc("/api/createblog", createblog).Methods("POST")

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

	return db

}

func main() {
	fmt.Printf("running. . . . . . ")

	handleRequests()

}
