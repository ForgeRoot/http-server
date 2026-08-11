package main

import (
  "encoding/json"
  "github.com/google/uuid"
  "net/http"
  "time"
  "strings"
  "github.com/ForgeRoot/http-server/internal/database"
  "fmt"
)

type response struct {
  ID        uuid.UUID `json:"id"`
  CreatedAt time.Time `json:"created_at"`
  UpdatedAt time.Time `json:"updated_at"`
  Body      string    `json:"body"`
  User_id   uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {
  chirps, err := cfg.dbQueries.GetChirps(r.Context())
  if err != nil {
    respondWithError(w, http.StatusInternalServerError, "Couldn't get all chirps", err)
    return
  }

  chirpsResponse := []response{}

  for i := 0; i < len(chirps); i++ {
    chirpsResponse = append(chirpsResponse,
    response{
      ID:        chirps[i].ID,
      CreatedAt: chirps[i].CreatedAt,
      UpdatedAt: chirps[i].UpdatedAt,
      Body:      chirps[i].Body,
      User_id:   chirps[i].UserID,
    })
  }
  respondWithJSON(w, http.StatusOK, chirpsResponse)
}

func (cfg *apiConfig) handlerGetChirpsId(w http.ResponseWriter, r *http.Request) {
  idString := r.PathValue("id")
  id, err := uuid.Parse(idString)
  if err != nil {
    respondWithError(w, http.StatusBadRequest, "Couldn't parse id string into UUID", err)
    return
  }

  chirp, err := cfg.dbQueries.GetChirpsId(r.Context(), id)
  if err != nil {
    respondWithError(w, http.StatusNotFound, fmt.Sprintf("Couldn't find chirp %s", idString), err)
    return
  }

  chirpResponse := response{
    ID:        chirp.ID,
    CreatedAt: chirp.CreatedAt,
    UpdatedAt: chirp.UpdatedAt,
    Body:      chirp.Body,
    User_id:   chirp.UserID,
  }

  respondWithJSON(w, http.StatusOK, chirpResponse)
}

func (cfg *apiConfig) handlerPostChirps(w http.ResponseWriter, r *http.Request) {
  type parameters struct {
    Body string `json:"body"`
    User_id uuid.UUID `json:"user_id"`
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
