package main

import _ "github.com/lib/pq"
import (
	"database/sql"
	"github.com/ForgeRoot/http-server/internal/database"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	dbQueries      *database.Queries
	platform       string
}

func main() {
	godotenv.Load()

	dbURL := os.Getenv("DB_URL")
	platform := os.Getenv("PLATFORM")

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
		dbQueries:      dbQueries,
		platform:       platform,
	}

	mux := http.NewServeMux()

	fsHandler := apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))))
	mux.Handle("/app/", fsHandler)

	mux.HandleFunc("GET /api/healthz", handlerReadiness)

	mux.HandleFunc("POST /api/users", apiCfg.handlerUsers)
	mux.HandleFunc("POST /api/chirps", apiCfg.handlerChirps)

	mux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)
	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)


	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving files on port: %s\n", port)
	log.Fatal(srv.ListenAndServe())
}
