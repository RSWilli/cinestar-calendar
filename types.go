package main

type CinestarMovie struct {
	Type             Type              `json:"_type"`
	ID               int64             `json:"id"`
	Cinema           int64             `json:"cinema"`
	Blockbuster      bool              `json:"blockbuster"`
	Title            string            `json:"title"`
	Subtitle         *string           `json:"subtitle"`
	HasTrailer       bool              `json:"hasTrailer"`
	Showtimes        []Showtime        `json:"showtimes"`
	Date             *string           `json:"date,omitempty"`
	Movie            *int64            `json:"movie,omitempty"`
	RelatedShows     []int64           `json:"relatedShows,omitempty"`
	PosterPreload    string            `json:"poster_preload"`
	Poster           string            `json:"poster"`
	DetailLink       string            `json:"detailLink"`
	StartDate        *string           `json:"startDate,omitempty"`
	ScreeningWeek    *int64            `json:"screeningWeek,omitempty"`
	ShowtimeSchedule *ShowtimeSchedule `json:"showtimeSchedule,omitempty"`
	Trailer          *int64            `json:"trailer,omitempty"`
	Event            *int64            `json:"event,omitempty"`
	EndDate          *string           `json:"endDate"`

	// Exists, but is hard to parse due to different formats depending on type
	// Attributes       []string          `json:"attributes"`
}

type ShowtimeSchedule struct {
	ID       int64   `json:"id"`
	Datetime string  `json:"datetime"`
	Text     string  `json:"text"`
	Type     *string `json:"type"`
}

type Showtime struct {
	ID         int64    `json:"id"`
	Name       string   `json:"name"`
	Cinema     int64    `json:"cinema"`
	Datetime   string   `json:"datetime"`
	Emv        *int64   `json:"emv"`
	Fsk        int64    `json:"fsk"`
	SystemID   string   `json:"systemId"`
	System     System   `json:"system"`
	Show       int64    `json:"show"`
	Attributes []string `json:"attributes"`
	Screen     *int64   `json:"screen"`
}

type System string

const (
	Vista System = "vista"
)

type Type string

const (
	Event Type = "event"
	Movie Type = "movie"
)
