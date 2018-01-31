package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/rs/cors"

	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// Article data
type Article struct {
	gorm.Model

	Link        string
	Name        string
	Author      string
	Description string
}

var db *gorm.DB
var err error

func main() {
	router := mux.NewRouter()

	db, err = gorm.Open("postgres", "host=localhost user=postgres dbname=gogogo sslmode=disable password=pgadmin")

	if err != nil {
		panic("Connection to database failed.")
	}

	defer db.Close()

	db.AutoMigrate(&Article{})

	router.HandleFunc("/articles", GetArticles).Methods("GET")
	router.HandleFunc("/articles/{id}", GetArticle).Methods("GET")
	router.HandleFunc("/articles", CreateArticle).Methods("POST")

	handler := cors.Default().Handler(router)

	log.Fatal(http.ListenAndServe(":8080", handler))
}

// GetArticles from database
func GetArticles(w http.ResponseWriter, r *http.Request) {
	var articles []Article

	if err := db.Find(&articles).Error; err != nil {
		panic("Error while fetching data from database")
	}

	json.NewEncoder(w).Encode(&articles)
}

// GetArticle by request id
func GetArticle(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var article Article

	if err := db.First(&article, params["id"]).Error; err != nil {
		panic("Error while fetching data from database")
	}

	json.NewEncoder(w).Encode(article)
}

// CreateArticle from request
func CreateArticle(w http.ResponseWriter, r *http.Request) {
	var article Article
	json.NewDecoder(r.Body).Decode(&article)

	if err := db.Create(&article).Error; err != nil {
		panic("Error creating data")
	}

	json.NewEncoder(w).Encode(&article)
}
