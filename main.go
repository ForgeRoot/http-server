package main

import _ "github.com/lib/pq"
import (
  "log"
  "net/http"
  "sync/atomic"
  "godotenv"
  "os"
  "database/sql"
)

type apiConfig struct {
  fileserverHits atomic.Int32
}

func main() {
  godotenv.Load()

  dbURL := os.Getenv("DB_URL")

  db, err := sql.Open("postgres", dbURL)
  if err != nil {
    log.Printf("Can't open db url: %v", err)
    return
  }
  dbQueries := database.New(db)

  const filepathRoot = "."
  const port = "8080"

  apiCfg := apiConfig{
    fileserverHits: atomic.Int32{},
    dbQueries: *database.Queries
  }

  mux := http.NewServeMux()

  fsHandler := apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))))
  mux.Handle("/app/", fsHandler)

  mux.HandleFunc("GET /api/healthz", handlerReadiness)
  mux.HandleFunc("POST /api/validate_chirp", handlerChirpsValidate)
  
  mux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)
  mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)

  srv := &http.Server{
    Addr:       ":" + port,
    Handler:    mux,
  }

  log.Printf("Serving files on port: %s\n", port)
  log.Fatal(srv.ListenAndServe())
}



