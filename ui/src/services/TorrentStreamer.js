import axios from "axios"

// TorrentStreamer consumes a torrent and returns a file. 
class TorrentStreamer {

    constructor({
        host = "http://localhost",
        port = "8080",
        baseURL = "player",
    }) {
        this.host = host
        this.port = `:${port}`
        this.baseURL = baseURL
        this.client = new axios()
    }

    watch(torrent) {
        if (typeof(torrent) != Torrent) {
            return new Promise.reject(`expected type Torrent, got ${typeof(torrent)}`)
        }
        this.post("/watch", torrent)
    }

    post(segment, torrent) {
        return this.client.post(`${this.host}${this.port}/${this.baseURL}/${segment}`, torrent)
    } 
}

// Torrent represents a torrent file.
// Can logically be a magnet link. 
class Torrent {
    constructor() {
        this.isMagnet = false
    }

}

export {
    TorrentStreamer,
    Torrent,
}