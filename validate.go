package main

import (
  "net/http"
  "encoding/json"
  "log"
)

func validateChirp(w http.ResponseWriter, r *http.Request) {
  type parameter struct {
    Body string `json:"body"`
  }

  type returnError struct {
    Error string `json:"error"`
  }
  
  type returnVal struct {
    Valid bool `json:"valid"`
  }

  decoder := json.NewDecoder(r.Body)
  param := parameter{}
  err := decoder.Decode(&param)
  if err != nil {
    log.Printf("Error decoding parameters: %s", err)
    w.WriteHeader(500)
    return
  }

  if len(param.Body) > 140 {
    respErr := returnError{ Error: "Chirp is too long" }
    dat, err := json.Marshal(respErr)
    if err != nil {
      log.Printf("Error marshalling JSON: %s", err)
      w.WriteHeader(500)
      return
    }
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(400)
    w.Write(dat)
    return
  }

  respVal := returnVal{ Valid: true }
  dat, err := json.Marshal(respVal)
  if err != nil {
    log.Printf("Error marshalling JSON: %s", err)
    w.WriteHeader(500)
    return
  }
  w.Header().Set("Content-Type", "application/json")
  w.WriteHeader(200)
  w.Write(dat)
  return

}
