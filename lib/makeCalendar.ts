import dayjs, { Dayjs } from "dayjs"
import { crlf } from "./lineend"

export type CalEvent = {
    date: Dayjs
    description: string
    title: string
    location: string
    duration: number
}

export const makeEvent = (event: CalEvent) => {
    return crlf`
        BEGIN:VEVENT
        DTSTAMP:${dayjs().format("YYYYMMDDTHHmmss")}
        DTSTART;TZID=Europe/Berlin:${event.date.format("YYYYMMDDTHHmmss")}
        DURATION:PT${event.duration}M
        SUMMARY:${event.title}
        DESCRIPTION:${event.description}
        LOCATION:${event.location}
        UID:cinestar_${event.date.unix()}
        END:VEVENT
    `
}
export const makeCalendar = (events: CalEvent[]) => {

    return crlf`BEGIN:VCALENDAR
        X-WR-CALDESC:cinestar
        X-WR-CALNAME:cinestar
        X-WR-TIMEZONE:Europe/Berlin
        CALSCALE:GREGORIAN
        PRODID:cinestar-parser
        VERSION:2.0
        ${events.map(makeEvent).join("")}
        BEGIN:VTIMEZONE
        TZID:Europe/Berlin
        BEGIN:STANDARD
        DTSTART:19701025T030000
        RRULE:FREQ=YEARLY;BYMONTH=10;BYDAY=-1SU
        TZNAME:CET
        TZOFFSETFROM:+0200
        TZOFFSETTO:+0100
        END:STANDARD
        BEGIN:DAYLIGHT
        DTSTART:19700329T020000
        RRULE:FREQ=YEARLY;BYMONTH=3;BYDAY=-1SU
        TZNAME:CEST
        TZOFFSETFROM:+0100
        TZOFFSETTO:+0200
        END:DAYLIGHT
        END:VTIMEZONE
        END:VCALENDAR
    `

}