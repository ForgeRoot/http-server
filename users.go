package main

import (
	"encoding/json"
	"github.com/google/uuid"
	"net/http"
	"time"
)

func (cfg *apiConfig) handlerUsers(w http.ResponseWriter, r *http.Request) {

	type parameters struct {
		Email string `json:"email"`
	}

	type response struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Email     string    `json:"email"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't decode parameters", err)
		return
	}

	createdUser, err := cfg.dbQueries.CreateUser(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create an user", err)
		return
	}

	userResponse := response{
		ID:        createdUser.ID,
		CreatedAt: createdUser.CreatedAt,
		UpdatedAt: createdUser.UpdatedAt,
		Email:     createdUser.Email,
	}

	respondWithJSON(w, http.StatusCreated, userResponse)
}
