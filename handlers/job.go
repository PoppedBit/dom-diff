package handlers

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/mux"
	"github.com/poppedbit/dom-diff/helpers"
	"github.com/poppedbit/dom-diff/models"
	"gorm.io/gorm"
)

type JobsData struct {
	helpers.BaseTemplateData
	Jobs []models.Job
}

func GetJobsHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var jobs []models.Job
		db.Find(&jobs)

		templates := []string{
			"templates/jobs.html",
		}

		tmpl, err := helpers.ParseFullPage(templates...)
		if err != nil {
			log.Println(err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		data := JobsData{
			Jobs: jobs,
		}
		data.BaseTemplateData.Init(r)

		err = tmpl.ExecuteTemplate(w, "base", data)
		if err != nil {
			log.Println(err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

	}
}

func GetCreateJobHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		tmpl, err := template.ParseFiles(
			"templates/_create_job_form.html",
		)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = tmpl.ExecuteTemplate(w, "createJobForm", nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func CreateJobHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()

		cemetery := models.Job{
			Url:          r.Form.Get("url"),
			ItemSelector: r.Form.Get("itemSelector"),
			TextSelector: r.Form.Get("textSelector"),
		}

		db.Create(&cemetery)

		id := cemetery.Id.String()

		// Create {outputDir}/{id}
		outputDir := filepath.Join(os.Getenv("OUTPUT_DIR"), id)
		err := os.MkdirAll(outputDir, os.ModePerm)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Redirect to cemeteries page
		w.Header().Set("HX-Location", "/job/"+id)
		w.WriteHeader(http.StatusCreated)
	}
}

func DeleteJobHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		id := params["id"]

		var job models.Job
		db.Where("id = ?", id).First(&job)
		db.Delete(&job)

		w.Header().Set("HX-Location", "/jobs")
		w.WriteHeader(http.StatusNoContent)
	}
}

type JobData struct {
	helpers.BaseTemplateData
	Job  models.Job
	Runs []models.Run
}

func GetJobHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		id := params["id"]

		var job models.Job
		db.Where("id = ?", id).First(&job)

		templates := []string{
			"templates/job.html",
		}

		tmpl, err := helpers.ParseFullPage(templates...)
		if err != nil {
			log.Println(err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		data := JobData{
			Job: job,
		}
		data.BaseTemplateData.Init(r)

		err = tmpl.ExecuteTemplate(w, "base", data)
		if err != nil {
			log.Println(err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}
}
