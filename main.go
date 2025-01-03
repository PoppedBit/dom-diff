package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/poppedbit/dom-diff/handlers"
	"github.com/poppedbit/dom-diff/helpers"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	models "github.com/poppedbit/dom-diff/models"
)

func main() {

	// ENV
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	// DB
	dbConnection := os.Getenv("DB_CONNECTION")
	db, err := gorm.Open(sqlite.Open(dbConnection), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error connecting to database")
	}

	err = db.AutoMigrate(&models.Job{}, &models.Run{})
	if err != nil {
		log.Fatalf("Error migrating database")
	}

	router := mux.NewRouter()

	// Routes
	router.HandleFunc("/jobs", handlers.GetJobsHandler(db)).Methods("GET")
	router.HandleFunc("/job", handlers.GetCreateJobHandler()).Methods("GET")
	router.HandleFunc("/job", handlers.CreateJobHandler(db)).Methods("POST")
	router.HandleFunc("/job/{id}", handlers.GetJobHandler(db)).Methods("GET")
	router.HandleFunc("/job/{id}", handlers.DeleteJobHandler(db)).Methods("DELETE")
	router.HandleFunc("/job/{id}/run", handlers.RunJobHandler(db)).Methods("POST")
	router.HandleFunc("/job/{jobId}/run/{runId}", handlers.GetRunHandler(db)).Methods("GET")
	router.HandleFunc("/job/{jobId}/run/{runId}", handlers.DeleteRunHandler(db)).Methods("DELETE")

	router.HandleFunc("/", GetAppHandler).Methods("GET")
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	router.PathPrefix("/webfonts/").Handler(http.StripPrefix("/webfonts/", http.FileServer(http.Dir("static/webfonts"))))

	// Server
	port := os.Getenv("PORT")
	println("Server running at http://localhost:" + port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}

type AppData struct {
	helpers.BaseTemplateData
}

func GetAppHandler(w http.ResponseWriter, r *http.Request) {

	templates := []string{
		"templates/index.html",
	}

	tmpl, err := helpers.ParseFullPage(templates...)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	data := AppData{}
	data.BaseTemplateData.Init(r)

	err = tmpl.ExecuteTemplate(w, "base", data)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
