package main

import (
	"crypto/sha1"
	"encoding/base32"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/gorilla/mux"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

const (
	audioExtension = ".mp3"
)

var (
	token = make(map[string]bool)
)

func init() {
	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found
			fmt.Println("Config file not found")
			return
		}
	}

	// Loads tokens
	for _, t := range strings.Split(viper.GetString("token"), ",") {
		token[t] = true
	}

	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		log.Println("Reloading config file")
		token = make(map[string]bool)

		for _, t := range strings.Split(viper.GetString("token"), ",") {
			token[t] = true
		}
	})

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

	uuid := genAudio(query.Get("text"))

	w.WriteHeader(http.StatusAccepted)
	_, _ = fmt.Fprintf(w, "/temp/"+uuid+audioExtension)
}

// Generates audio from a string. Checks if it already exist before generating it
func gen(text string, uuid string) {
	_, err := os.Stat("./temp/" + uuid + audioExtension)

	if err != nil {
		tts := exec.Command("balcon", "-i", "-o", "-enc", "utf8", "-n", "Roberto")
		tts.Stdin = strings.NewReader(text)
		ttsOut, _ := tts.StdoutPipe()
		_ = tts.Start()

		ffmpeg := exec.Command("ffmpeg", "-i", "pipe:0", "-f", "s16le", "-ar", "48000", "-ac", "2", "-f", "mp3", "./temp/"+uuid+audioExtension)
		ffmpeg.Stdin = ttsOut
		_ = ffmpeg.Run()

		_ = tts.Wait()
	}
}

// genAudio generates a mp3 file from a string, returning it's UUID (aka SHA1 hash of the text)
func genAudio(text string) string {
	h := sha1.New()
	h.Write([]byte(text))
	uuid := strings.ToUpper(base32.HexEncoding.EncodeToString(h.Sum(nil)))

	gen(text, uuid)

	return uuid
}
