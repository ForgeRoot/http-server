package main

import (
  "net/http"
  "encoding/json"
  "strings"
  "slices"
)

func handlerChirpsValidate(w http.ResponseWriter, r *http.Request) {
  type parameters struct {
    Body string `json:"body"`
  }

  type returnVals struct {
    CleanedBody string `json:"cleaned_body"`
  }

  decoder := json.NewDecoder(r.Body)
  params := parameters{}
  err := decoder.Decode(&params)
  if err != nil {
    respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
    return
  }

  const maxChirpLength = 140
  if len(params.Body) > maxChirpLength {
    respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
    return
  }
  
  splittedBody := strings.Split(params.Body, " ")

  forbidden_words := []string{"kerfuffle", "sharbert", "fornax"}
  for i := 0; i < len(splittedBody); i++ {
    if slices.Contains(forbidden_words, strings.ToLower(splittedBody[i])) {
      splittedBody[i] = "****"
    }
  }
  
  respondWithJSON(w, http.StatusOK, returnVals{
    CleanedBody: strings.Join(splittedBody, " "),
  })
}
