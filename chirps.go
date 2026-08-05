package main

import (
	"encoding/json"
	"github.com/google/uuid"
	"net/http"
	"time"
	"strings"
	"github.com/ForgeRoot/http-server/internal/database"
)

func (cfg *apiConfig) handlerChirps(w http.ResponseWriter, r *http.Request) {

	type parameters struct {
		Body string `json:"body"`
		User_id uuid.UUID `json:"user_id"`
	}

	type response struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Body      string    `json:"body"`
		User_id   uuid.UUID `json:"user_id"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't decode parameters", err)
		return
	}

	const maxChirpLength = 140
	if len(params.Body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}

	badWords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}

	cleaned := getCleanedBody(params.Body, badWords)
	
	chirpParams := database.CreateChirpParams{
		Body:  cleaned,
		UserID: params.User_id,
	}

	createdChirp, err := cfg.dbQueries.CreateChirp(r.Context(), chirpParams)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create an chirp", err)
		return
	}

	chirpResponse := response{
		ID:        createdChirp.ID,
		CreatedAt: createdChirp.CreatedAt,
		UpdatedAt: createdChirp.UpdatedAt,
		Body:      createdChirp.Body,
		User_id:   createdChirp.UserID,
	}

	respondWithJSON(w, http.StatusCreated, chirpResponse)
}

func getCleanedBody(body string, badWords map[string]struct{}) string {
	words := strings.Split(body, " ")
	for i, word := range words {
		loweredWord := strings.ToLower(word)
		if _, ok := badWords[loweredWord]; ok {
			words[i] = "****"
		}
	}
	cleaned := strings.Join(words, " ")
	return cleaned
}
