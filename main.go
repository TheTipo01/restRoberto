package main

import (
	"net/http"
	"os/exec"
	"strings"

	"github.com/bwmarrin/lit"
	"github.com/gorilla/mux"
	"github.com/kkyr/fig"
)

var (
	token      map[string]bool
	balconPath string
	address    string
)

func init() {
	var cfg config

	lit.LogLevel = lit.LogError

	err := fig.Load(&cfg, fig.File("config.yml"))
	if err != nil {
		lit.Error(err.Error())
		return
	}

	balconPath = cfg.BalconPath
	address = cfg.Address

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
	_ = http.ListenAndServe(address, nil)
}

func audio(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	if !token[query.Get("token")] {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	tts := exec.Command(balconPath, "-i", "-o", "-enc", "utf8", "-n", "Roberto", "-fr", "16", "-ch", "1", "-bt", "16")
	tts.Stdin = strings.NewReader(query.Get("text"))
	tts.Stdout = w
	_ = tts.Start()

	w.WriteHeader(http.StatusAccepted)

	_ = tts.Wait()
}
