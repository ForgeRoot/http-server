package main

import (
  "net/http"
  "log"
  "io"
  "sync/atomic"
  "fmt"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        cfg.fileserverHits.Add(1)
        log.Printf("Hit! Count: %d", cfg.fileserverHits.Load())
        next.ServeHTTP(w, r)
    })
}

func main() {
  const filepathRoot = "."
  const port = "8080"
  apiCfg := apiConfig{}

  mux := http.NewServeMux()
  mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))))
  mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
	  w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	  w.WriteHeader(http.StatusOK)
	  io.WriteString(w, "OK")
  })
  mux.HandleFunc("/reset", func(w http.ResponseWriter, req *http.Request) {
	  w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	  w.WriteHeader(http.StatusOK)
	  apiCfg.fileserverHits.Store(0)
  })
  mux.HandleFunc("/metrics", func(w http.ResponseWriter, req *http.Request) {
	  w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	  w.WriteHeader(http.StatusOK)
	  io.WriteString(w, fmt.Sprintf("Hits: %v", apiCfg.fileserverHits.Load()))
  })

  srv := &http.Server{
    Addr:       ":" + port,
    Handler:    mux,
  }

  log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
  log.Fatal(srv.ListenAndServe())
}


