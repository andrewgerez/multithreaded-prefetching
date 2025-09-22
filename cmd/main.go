package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/andrewgerez/multithreaded-prefetching/internal/handler"
	"github.com/andrewgerez/multithreaded-prefetching/internal/tmdb"

	"github.com/joho/godotenv"
)

func main() {
    // Carregar variáveis do .env
    if err := godotenv.Load(); err != nil {
        log.Println(".env file not found, proceeding with system env vars")
    }

    apiKey := os.Getenv("TMDB_API_KEY")
    if apiKey == "" {
        log.Fatal("TMDB_API_KEY not set")
    }

    tmdbClient := tmdb.NewClient(apiKey)
    h := &handler.Handler{TmdbClient: tmdbClient}

    http.HandleFunc("/profiles", h.Profiles)
    http.HandleFunc("/profiles/", func(w http.ResponseWriter, r *http.Request) {
        if r.URL.Path == "/profiles/homes" {
            h.AllProfilesHomes(w, r)
        } else if strings.HasSuffix(r.URL.Path, "/home") {
            h.ProfileHome(w, r)
        } else {
            http.NotFound(w, r)
        }
    })

    log.Println("Servidor ouvindo na porta :8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}
