package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/hashicorp/go-hclog"

	"github.com/aaronschweig/auto-sdb/extractor"
)

type ErrResponse struct {
	Message string `json:"message"`
}

func NewErrResponse(err error) *ErrResponse {
	return &ErrResponse{err.Error()}
}

func main() {
	log := hclog.Default()
	r := chi.NewRouter()

	// r.Use(middleware.SetHeader("content-type", "application/json"))
	r.Use(middleware.BasicAuth("SDB-Extractor", map[string]string{
		"admin": "admin",
	}))

	r.Get("/", func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Set("content-type", "text/html")
		http.ServeFile(rw, r, "frontend/index.html")
	})

	r.Post("/extract", func(rw http.ResponseWriter, r *http.Request) {
		r.ParseMultipartForm(10 << 20)
		rw.Header().Set("content-type", "application/json")

		file, _, _ := r.FormFile("file")
		defer file.Close()

		tempFile, _ := ioutil.TempFile(".", "sdb-*.pdf")
		defer tempFile.Close()
		defer os.Remove(tempFile.Name())

		content, _ := ioutil.ReadAll(file)

		tempFile.Write(content)

		cmd := exec.Command("gs", "-sDEVICE=txtwrite", "-dBATCH", "-dNOPAUSE", "-sOutputFile=-", tempFile.Name())

		out, err := cmd.Output()

		if err != nil {
			log.Error("could not process pdf with gs", "error", err)

			rw.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(rw).Encode(NewErrResponse(err))
			return
		}

		e := extractor.NewDefaultExtractor().WithContent(string(out))

		result := e.Extract()

		json.NewEncoder(rw).Encode(&result)
	})

	port := os.Getenv("PORT")

	if len(port) == 0 {
		port = "3000"
	}

	log.Info(fmt.Sprintf("Application is up and running on http://localhost:%s", port))
	http.ListenAndServe(fmt.Sprintf(":%s", port), r)
}
