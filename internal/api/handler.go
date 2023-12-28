package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"prompt-generator/internal/db"
)

// GetPromptsHandler handles the GET request to fetch all prompts
func GetPromptsHandler(w http.ResponseWriter, r *http.Request) {
	prompts, err := db.GetPrompts()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(prompts)
}

// PostPromptHandler handles the POST request to create a new prompt
func PostPromptHandler(w http.ResponseWriter, r *http.Request) {
	var p db.Prompt
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	promptID, err := db.CreatePrompt(p.Prompt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]int{"id": promptID})
}

// GetPromptHandler handles the GET request to fetch a single prompt by ID
func GetPromptHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	prompt, err := db.GetPrompt(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(prompt)
}

// DeletePromptHandler handles the DELETE request to delete a prompt
func DeletePromptHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = db.DeletePrompt(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
