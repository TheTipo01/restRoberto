package main

import (
	"fmt"
	libroberto "github.com/TheTipo01/libRoberto"
	"github.com/bwmarrin/lit"
	"github.com/gorilla/mux"
	"github.com/kkyr/fig"
	"net/http"
	"strings"
	"time"
)

const (
	audioExtension = ".mp3"
)

var (
	token = make(map[string]bool)
)

func init() {
	var cfg config

	lit.LogLevel = lit.LogError

	err := fig.Load(&cfg, fig.File("config.yml"))
	if err != nil {
		lit.Error(err.Error())
		return
	}

	// Set lit.LogLevel to the given value
	switch strings.ToLower(cfg.LogLevel) {
	case "logwarning", "warning":
		lit.LogLevel = lit.LogWarning

	case "loginformational", "informational":
		lit.LogLevel = lit.LogInformational

	case "logdebug", "debug":
		lit.LogLevel = lit.LogDebug
	}

	// Loads tokens
	for _, t := range cfg.Token {
		token[t] = true
	}
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/audio", audio).Methods("GET")
	r.PathPrefix("/temp/").Handler(http.StripPrefix("/temp/", http.FileServer(http.Dir("./temp"))))

	http.Handle("/", r)
	http.ListenAndServe(":8087", nil)
}

func audio(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	if !token[query.Get("token")] {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	uuid := libroberto.GenAudio(query.Get("text"), audioExtension, 30*time.Second)

	w.WriteHeader(http.StatusAccepted)
	_, _ = fmt.Fprintf(w, "/temp/"+uuid+audioExtension)
}
