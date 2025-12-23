package main

import (
	"net/http"
	"strings"

	libroberto "github.com/TheTipo01/libRoberto"
	"github.com/bwmarrin/lit"
	"github.com/gorilla/mux"
	"github.com/kkyr/fig"
)

var (
	token map[string]bool
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
	token = make(map[string]bool, len(cfg.Token))
	for _, t := range cfg.Token {
		token[t] = true
	}
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/audio", audio).Methods("GET")

	http.Handle("/", r)
	_ = http.ListenAndServe(":8087", nil)
}

func audio(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	if !token[query.Get("token")] {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	cmds := libroberto.GenRawAudio(query.Get("text"))
	cmds[0].Stdout = w
	libroberto.CmdsStart(cmds)

	w.WriteHeader(http.StatusAccepted)

	libroberto.CmdsWait(cmds)
}
