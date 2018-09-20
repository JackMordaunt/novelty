// Episode contains meta data about a watchable episode of a show. 
class Episode {
    constructor({
        name = "",
        releaseDate = new Date(),
        runtime = 0,
        synopsis = "",
        uri = "",
    }) {
        this.name = name
        this.releaseDate = releaseDate
        this.runtime = runtime
        this.synopsis = synopsis
        this.uri = uri
    }
}

export default Episode