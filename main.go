package main

import (
	"bytes"
	"embed"
	"encoding/json"
	"flag"
	"fmt"
	"io/fs"
	"net/http"
	"os/exec"

	"github.com/hashicorp/go-hclog"

	"github.com/aaronschweig/auto-sdb/extractor"
)

type ErrResponse struct {
	Message string `json:"message"`
}

var (
	//go:embed frontend/build/* frontend/build/_app/pages/* frontend/build/_app/assets/pages/*
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

		cmd := exec.Command("gs", "-sDEVICE=txtwrite", "-dBATCH", "-dNOPAUSE", "-sOutputFile=-", "-")

		cmd.Stdin = file

		var buffer bytes.Buffer
		cmd.Stdout = &buffer

		err = cmd.Run()

		if err != nil {
			log.Error("could not process pdf with gs", "error", err)

			writeError(rw, http.StatusInternalServerError, err)
			return
		}

		result := extractor.Extract(buffer.String(), log)

		rw.Header().Add("Content-Type", "application/json")
		json.NewEncoder(rw).Encode(&result)
	}
}

func post(f http.HandlerFunc) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.NotFound(rw, r)
			return
		}
		f.ServeHTTP(rw, r)
	}
}

func main() {
	port := flag.String("port", "3000", "the port for the application to run on")
	dev := flag.Bool("dev", false, "start server in dev mode and do not bundle frontend in binary")
	flag.Parse()

	log := hclog.Default()
	mux := http.NewServeMux()

	if *dev {
		mux.Handle("/", http.FileServer(http.Dir("./frontend/build")))
	} else {
		static, err := fs.Sub(frontend, "frontend/build")
		if err != nil {
			panic(err)
		}
		mux.Handle("/", http.FileServer(http.FS(static)))
	}

	mux.HandleFunc("/extract", post(extractSDB(log)))

	log.Info(fmt.Sprintf("Application is up and running on http://localhost:%s", *port))
	http.ListenAndServe(fmt.Sprintf(":%s", *port), mux)
}
