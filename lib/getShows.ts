import fetch from "node-fetch"
import { Movie } from "./types"

export const getShows = async () => {
    const res = await fetch("https://www.cinestar.de/api/cinema/33/show/")

    if (res.status !== 200) {
        throw new Error("Failed to fetch shows")
    }

    const json = await res.json() as Movie[]

    return json

}
