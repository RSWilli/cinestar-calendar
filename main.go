package main

import (
	"flag"
	"log"
	"net/http"
	"slices"
	"strings"
	"time"
)

var ics []byte

var addr = flag.String("addr", ":8080", "HTTP network address")

func main() {
	flag.Parse()

	go func() {
		t := time.NewTicker(30 * time.Minute)

		updateCalendar()
		for range t.C {
			updateCalendar()
		}
	}()

	http.HandleFunc("/calendar", calendarHandler)
	log.Printf("Starting server on :%s...", *addr)
	if err := http.ListenAndServe(*addr, nil); err != http.ErrServerClosed {
		log.Fatalf("Server failed: %v", err)
	}
}

func updateCalendar() {
	shows, err := getShows()
	if err != nil {
		log.Printf("Error fetching shows: %v", err)
		return
	}

	var events []CalEvent
	for _, show := range shows {
		// filter non OV showtimes and events
		var OVs []Showtime
		for _, st := range show.Showtimes {
			if slices.Contains(st.Attributes, "OV") &&
				!slices.Contains(st.Attributes, "Alt ConEFF") &&
				!slices.Contains(st.Attributes, "Alt ConUFP") {
				OVs = append(OVs, st)
			}
		}

		if len(OVs) == 0 {
			continue
		}

		runtime, err := getMovieRuntime(show.Title)
		if err != nil {
			log.Printf("Error fetching runtime for %s: %v", show.Title, err)
			runtime = 120 // default to 2 hours if runtime is unavailable
		}

		for _, st := range show.Showtimes {
			summary := show.Title

			if slices.Contains(st.Attributes, "2D") {
				summary = "(2D)" + summary
			} else if slices.Contains(st.Attributes, "3D") {
				summary = "(3D)" + summary
			}

			// e.g. 2026-03-11 16:45 CET
			date, err := time.Parse("2006-01-02 15:04 MST", st.Datetime)

			if err != nil {
				log.Printf("Error parsing date for show %s: %v", show.Title, err)
				continue
			}

			events = append(events, CalEvent{
				Date:        date,
				Description: strings.Join(st.Attributes, " "),
				Title:       summary,
				Location:    "Cinestar Leipzig",
				Duration:    runtime,
			})
		}
	}

	ics = MakeCalendar(events)

	// debug output calendar to file
	// if err := os.WriteFile("calendar.ics", ics, 0644); err != nil {
	// 	log.Printf("Error writing calendar to file: %v", err)
	// }
}

func calendarHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/calendar; charset=utf-8")
	w.Write(ics)
}
