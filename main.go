package main

import (
	"encoding/json"
	"fmt"
	"gospamfilter/core"
	"log"
	"net/http"
	"strings"

	"github.com/joho/godotenv"
)

type RequestBody struct {
	Text string `json:"text"`
}

type Response struct {
	Text   string `json:"text"`
	IsSpam bool   `json:"is_spam"`
}

func detectSpam(text string) bool {
	system := "You are a Spam Detector."
	prompt := fmt.Sprintf(`
		Return "true" if the given text is spam.
		Return "false" if the given text is not spam.
		Text: "%s"
	`, text)
	result := core.GetCompletion(system, prompt)

	fmt.Println(result)

	return strings.Trim(result, " ") == "true"
}

func handler(w http.ResponseWriter, r *http.Request) {
	var body RequestBody
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response := Response{
		Text:   body.Text,
		IsSpam: detectSpam(body.Text),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	http.HandleFunc("/detect-spam", handler)
	http.ListenAndServe(":8080", nil)
}
