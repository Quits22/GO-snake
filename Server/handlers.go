package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type Score struct {
	Name   string `json:"name"`
	Points string `json:"points"`
	Date   string `json:"date"`
}

// this gets the score from the db
func getScoreHandler(w http.ResponseWriter, r *http.Request) {
	scores, err := store.GetScores()
	if err != nil {
		log.Println("Failed to get scores:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	//marshal the json to make it readable by go
	scoreListBytes, err := json.Marshal(scores)
	if err != nil {
		log.Println("Failed to marshal scores:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	//sending the json
	w.Write(scoreListBytes)
}

// this is the post used by the client to recive the score and store it in the db
func createScoreHandler(w http.ResponseWriter, r *http.Request) {
	score := Score{}

	err := json.NewDecoder(r.Body).Decode(&score)
	if err != nil {
		log.Println("Failed to parse score:", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// Process the score data and store it in the database
	err = store.CreateScore(&score)
	if err != nil {
		log.Println("Failed to create score:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Respond with a success message if everything goes well
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Score created successfully"))
}
