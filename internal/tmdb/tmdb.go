package tmdb

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type TMDBClient struct {
    APIKey string
    BaseURL string
    Client  *http.Client
}

func NewClient(apiKey string) *TMDBClient {
    return &TMDBClient{
        APIKey:  apiKey,
        BaseURL: "https://api.themoviedb.org/3",
        Client:  http.DefaultClient,
    }
}

func (c *TMDBClient) DiscoverByGenre(genre string) (interface{}, error) {
    url := fmt.Sprintf("%s/discover/movie?api_key=%s&with_genres=%s", c.BaseURL, c.APIKey, genre)
    resp, err := c.Client.Get(url)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, err
    }
    var result map[string]interface{}
    if err := json.Unmarshal(body, &result); err != nil {
        return nil, err
    }
    return result, nil
}
