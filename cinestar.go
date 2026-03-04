package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// getShows fetches showtimes from the Cinestar API and returns them as a slice
func getShows() ([]CinestarMovie, error) {
	resp, err := http.Get("https://www.cinestar.de/api/cinema/33/show/")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch shows: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %s", resp.Status)
	}

	var shows []CinestarMovie

	if err := json.NewDecoder(resp.Body).Decode(&shows); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	return shows, nil
}
