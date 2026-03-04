package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
)

// TMDbMovie represents the relevant part of the TMDb API response
// Only the runtime field is needed for calendar generation

type TMDbMovie struct {
	Runtime int `json:"runtime"`
}

// getMovieRuntime fetches the runtime for a movie title from TMDb
func getMovieRuntime(title string) (int, error) {
	apiKey := os.Getenv("MOVIEDB_API_KEY")
	if apiKey == "" {
		return 0, fmt.Errorf("MOVIEDB_API_KEY not set")
	}

	// Search for the movie by title
	searchURL := fmt.Sprintf("https://api.themoviedb.org/3/search/movie?api_key=%s&query=%s", apiKey, url.QueryEscape(title))
	resp, err := http.Get(searchURL)
	if err != nil {
		return 0, fmt.Errorf("failed to search movie: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("unexpected status: %s", resp.Status)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("failed to read response: %w", err)
	}
	var searchResult struct {
		Results []struct {
			ID int `json:"id"`
		} `json:"results"`
	}
	if err := json.Unmarshal(body, &searchResult); err != nil {
		return 0, fmt.Errorf("failed to parse search JSON: %w", err)
	}
	if len(searchResult.Results) == 0 {
		return 0, fmt.Errorf("movie not found: %s", title)
	}
	movieID := searchResult.Results[0].ID
	// Fetch movie details for runtime
	detailURL := fmt.Sprintf("https://api.themoviedb.org/3/movie/%d?api_key=%s", movieID, apiKey)
	detailResp, err := http.Get(detailURL)
	if err != nil {
		return 0, fmt.Errorf("failed to fetch movie details: %w", err)
	}
	defer detailResp.Body.Close()
	if detailResp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("unexpected status: %s", detailResp.Status)
	}
	detailBody, err := io.ReadAll(detailResp.Body)
	if err != nil {
		return 0, fmt.Errorf("failed to read detail response: %w", err)
	}
	var movie TMDbMovie
	if err := json.Unmarshal(detailBody, &movie); err != nil {
		return 0, fmt.Errorf("failed to parse detail JSON: %w", err)
	}
	return movie.Runtime, nil
}
