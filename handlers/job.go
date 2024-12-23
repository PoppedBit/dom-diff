package handlers

import (
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/PuerkitoBio/goquery"
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

		var runs []models.Run
		db.Where("job_id = ?", id).Find(&runs)

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
			Job:  job,
			Runs: runs,
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

func RunJobHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		id := params["id"]

		var job models.Job
		db.Where("id = ?", id).First(&job)

		run := models.Run{
			JobId: id,
		}

		db.Create(&run)

		runId := run.Id.String()

		// Get html from url
		response, err := http.Get(job.Url)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		defer response.Body.Close()
		body, err := io.ReadAll(response.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Save html to {outputDir}/response.html
		outputDir := filepath.Join(os.Getenv("OUTPUT_DIR"), runId)
		err = os.MkdirAll(outputDir, os.ModePerm)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = os.WriteFile(filepath.Join(outputDir, "response.html"), body, os.ModePerm)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		html := string(body)

		// get all elements that match itemSelector
		doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
		if err != nil {
			log.Fatalf("Failed to parse HTML: %v", err)
		}

		// Select elements with the class "example"
		doc.Find(job.ItemSelector).Each(func(index int, container *goquery.Selection) {
			// inside element, find job.TextSelector, and get the inner text
			elements := container.Find(job.TextSelector)

			elements.Each(func(index int, element *goquery.Selection) {
				// text := element.Text()
			})
		})

		// Redirect to run page
		w.Header().Set("HX-Location", "/job/"+id+"/run/"+runId)
		w.WriteHeader(http.StatusCreated)
	}
}
