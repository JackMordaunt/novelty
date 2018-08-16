// Episode contains meta data about a watchable episode of a show. 
class Episode {
    constructor({
        name = "",
        releaseDate = new Date(),
        runtime = 0,
        synopsis = "",
    }) {
        this.name = name
        this.releaseDate = releaseDate
        this.runtime = runtime
        this.synopsis = synopsis
    }
}

export default Episode