package handler

import (
	"encoding/json"
	"net/http"
	"strings"
	"sync"

	"github.com/andrewgerez/multithreaded-prefetching/internal/model"
	"github.com/andrewgerez/multithreaded-prefetching/internal/profile"
	"github.com/andrewgerez/multithreaded-prefetching/internal/tmdb"
)

type Handler struct {
    TmdbClient *tmdb.TMDBClient
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    json.NewEncoder(w).Encode(data)
}

func (h *Handler) Profiles(w http.ResponseWriter, r *http.Request) {
    profiles := profile.GetAll()
    writeJSON(w, http.StatusOK, profiles)
}

func (h *Handler) ProfileHome(w http.ResponseWriter, r *http.Request) {
    parts := strings.Split(r.URL.Path, "/")
    if len(parts) < 4 {
        http.Error(w, "Invalid path", http.StatusBadRequest)
        return
    }
    profileID := parts[2]
    prof := profile.GetByID(profileID)
    if prof == nil {
        http.Error(w, "Profile not found", http.StatusNotFound)
        return
    }
    data, err := h.TmdbClient.DiscoverByGenre(prof.Genre)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    writeJSON(w, http.StatusOK, data)
}

func (h *Handler) AllProfilesHomes(w http.ResponseWriter, r *http.Request) {
    profilesList := profile.GetAll()
    type result struct {
        ProfileID string
        Home      interface{}
        Err       error
    }
    results := make([]result, len(profilesList))

    var wg sync.WaitGroup

    for i, prof := range profilesList {
        wg.Add(1)
        go func(i int, prof model.Profile) {
            defer wg.Done()
            home, err := h.TmdbClient.DiscoverByGenre(prof.Genre)
            results[i] = result{
                ProfileID: prof.ID,
                Home:      home,
                Err:       err,
            }
        }(i, prof)
    }
    wg.Wait()

    // Mapeia resposta by profileID
    resp := make(map[string]interface{})
    for _, res := range results {
        if res.Err != nil {
            resp[res.ProfileID] = map[string]string{"error": res.Err.Error()}
        } else {
            resp[res.ProfileID] = res.Home
        }
    }
    writeJSON(w, http.StatusOK, resp)
}
