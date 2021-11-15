import { CalEvent, makeCalendar } from "./lib/makeCalendar"
import express from "express"
import { writeFile } from "fs/promises"
import nocache from "nocache"
import { getShows } from "./lib/getShows"
import dayjs from "dayjs"
import customParseFormat from 'dayjs/plugin/customParseFormat'
import { config } from "dotenv"
import { MovieDb } from "moviedb-promise"
import { inspect } from "util"
dayjs.extend(customParseFormat)

config()

const movieDB = new MovieDb(process.env.MOVIEDB_API_KEY!)


// get movie runtime from movieDB
async function getRuntime(title: string) {
    // filter "the" from string
    const titleNoThe = title.replace(/^(the\s+)?/i, "")

    const movie = await movieDB.searchMovie({ query: titleNoThe })

    if (!movie.results || movie.results.length === 0 || movie.results[0].id === undefined) {
        throw new Error(`No movie found for ${title}`)
    }

    const movieByID = await movieDB.movieInfo({ id: movie.results[0].id! })

    if (!movieByID.runtime) {
        throw new Error(`No runtime found for ${title}`)
    }

    return movieByID.runtime
}

getRuntime("The Eternals")

let ics: string

const app = express()

app.use(nocache())
app.set("etag", false)
app.disable('view cache')

async function refetchCalendar() {
    const shows = await getShows()

    const events: CalEvent[] = []

    for (const show of shows) {
        const { title } = show

        const showtimes = show.showtimes.filter(showtime =>
            showtime.attributes.includes("OV")
            && !showtime.attributes.includes("Alt ConEFF")
        )
        if (showtimes.length === 0) {
            continue
        }

        const runtime = await getRuntime(title).catch(() => 120)

        events.push(...showtimes.map<CalEvent>(showtime => {

            let summary = title

            if (showtime.attributes.includes("2D")) {
                summary = "(2D) " + summary
            } else if (showtime.attributes.includes("3D")) {
                summary = "(3D) " + summary
            }

            const date = dayjs(showtime.datetime, "YYYY-MM-DD HH:mm")

            return {
                date,
                description: showtime.attributes.join(" "),
                location: "",
                title: summary,
                duration: runtime
            }
        }))
    }



    ics = makeCalendar(events)

    writeFile("calendar.ics", ics)
}

refetchCalendar()

setInterval(refetchCalendar, 1000 * 60 * 30)

app.get("/calendar", (req, res) => {
    res.contentType("text/calendar")
    res.send(ics)
})

app.listen(3001, "0.0.0.0")