package profile

import "github.com/andrewgerez/multithreaded-prefetching/internal/model"

var profiles = []model.Profile{
    {"1", "User 1", "28"},    // Ação
    {"2", "User 2", "35"},    // Comédia
    {"3", "User 3", "18"},    // Drama
}

func GetAll() []model.Profile {
    return profiles
}

func GetByID(id string) *model.Profile {
    for _, p := range profiles {
        if p.ID == id {
            return &p
        }
    }
    return nil
}
