package main

import (
	"bytes"
	"fmt"
	"io"
	"time"
)

var locationBerlin *time.Location

func init() {
	var err error
	locationBerlin, err = time.LoadLocation("Europe/Berlin")
	if err != nil {
		panic("Failed to load timezone: " + err.Error())
	}
}

type CalEvent struct {
	Date        time.Time
	Description string
	Title       string
	Location    string
	Duration    int // in minutes
}

func makeEvent(w io.Writer, event CalEvent) {
	fmt.Fprint(w, "BEGIN:VEVENT\r\n")
	fmt.Fprintf(w, "DTSTAMP:%s\r\n", time.Now().Format("20060102T150405"))
	fmt.Fprintf(w, "DTSTART;TZID=Europe/Berlin:%s\r\n", event.Date.In(locationBerlin).Format("20060102T150405"))
	fmt.Fprintf(w, "DURATION:PT%dM\r\n", event.Duration)
	fmt.Fprintf(w, "SUMMARY:%s\r\n", event.Title)
	fmt.Fprintf(w, "DESCRIPTION:%s\r\n", event.Description)
	fmt.Fprintf(w, "LOCATION:%s\r\n", event.Location)
	fmt.Fprintf(w, "UID:cinestar_%d\r\n", event.Date.Unix())
	fmt.Fprint(w, "END:VEVENT\r\n")
}

func MakeCalendar(events []CalEvent) []byte {
	var sb bytes.Buffer

	sb.WriteString("BEGIN:VCALENDAR\r\n")
	sb.WriteString("X-WR-CALDESC:cinestar\r\n")
	sb.WriteString("X-WR-CALNAME:cinestar\r\n")
	sb.WriteString("X-WR-TIMEZONE:Europe/Berlin\r\n")
	sb.WriteString("CALSCALE:GREGORIAN\r\n")
	sb.WriteString("PRODID:cinestar-parser\r\n")
	sb.WriteString("VERSION:2.0\r\n")

	for _, event := range events {
		makeEvent(&sb, event)
	}

	sb.WriteString("BEGIN:VTIMEZONE\r\n")
	sb.WriteString("TZID:Europe/Berlin\r\n")
	sb.WriteString("BEGIN:STANDARD\r\n")
	sb.WriteString("DTSTART:19701025T030000\r\n")
	sb.WriteString("RRULE:FREQ=YEARLY;BYMONTH=10;BYDAY=-1SU\r\n")
	sb.WriteString("TZNAME:CET\r\n")
	sb.WriteString("TZOFFSETFROM:+0200\r\n")
	sb.WriteString("TZOFFSETTO:+0100\r\n")
	sb.WriteString("END:STANDARD\r\n")
	sb.WriteString("BEGIN:DAYLIGHT\r\n")
	sb.WriteString("DTSTART:19700329T020000\r\n")
	sb.WriteString("RRULE:FREQ=YEARLY;BYMONTH=3;BYDAY=-1SU\r\n")
	sb.WriteString("TZNAME:CEST\r\n")
	sb.WriteString("TZOFFSETFROM:+0100\r\n")
	sb.WriteString("TZOFFSETTO:+0200\r\n")
	sb.WriteString("END:DAYLIGHT\r\n")
	sb.WriteString("END:VTIMEZONE\r\n")
	sb.WriteString("END:VCALENDAR\r\n")

	return sb.Bytes()
}
