package main

import (
	"embed"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/fs"
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

var (
	//go:embed frontend/*
	frontend embed.FS
)

func writeError(rw http.ResponseWriter, statusCode int, err error) {
	rw.WriteHeader(statusCode)

	rw.Header().Add("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(&ErrResponse{err.Error()})
}

func extractSDB(log hclog.Logger) func(http.ResponseWriter, *http.Request) {
	return func(rw http.ResponseWriter, r *http.Request) {
		file, _, err := r.FormFile("file")
		if err != nil {
			log.Error("error reading file", "error", err)

			writeError(rw, http.StatusBadRequest, err)
			return
		}
		defer file.Close()

		tempFile, err := os.CreateTemp(".", "sdb-*.pdf")
		if err != nil {
			log.Error("Error creating temp file", "error", err)

			writeError(rw, http.StatusInternalServerError, err)
			return
		}
		defer tempFile.Close()
		defer os.Remove(tempFile.Name())

		_, err = io.Copy(tempFile, file)
		if err != nil {
			log.Error("error saving file", "error", err)

			writeError(rw, http.StatusInternalServerError, err)
			return
		}

		cmd := exec.Command("gs", "-sDEVICE=txtwrite", "-dBATCH", "-dNOPAUSE", "-sOutputFile=-", tempFile.Name())

		out, err := cmd.Output()

		if err != nil {
			log.Error("could not process pdf with gs", "error", err)

			writeError(rw, http.StatusInternalServerError, err)
			return
		}

		e := extractor.NewDefaultExtractor(extractor.WithContent(string(out)), extractor.WithLogger(log))

		result := e.Extract()

		rw.Header().Add("Content-Type", "application/json")
		json.NewEncoder(rw).Encode(&result)
	}
}

func main() {
	port := flag.String("port", "3000", "the port for the application to run on")
	flag.Parse()

	log := hclog.Default()
	r := chi.NewRouter()

	r.Use(middleware.BasicAuth("SDB-Extractor", map[string]string{"admin": "admin"}))

	static, err := fs.Sub(frontend, "frontend")
	if err != nil {
		panic(err)
	}

	r.Handle("/", http.FileServer(http.FS(static)))

	r.Post("/extract", extractSDB(log))

	log.Info(fmt.Sprintf("Application is up and running on http://localhost:%s", *port))
	http.ListenAndServe(fmt.Sprintf(":%s", *port), r)
}
