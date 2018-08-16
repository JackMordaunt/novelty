// Show represents a watchable show. 
class Show {
    constructor({
        uuid,
        name = "",
        year = "",
        rating = 0,
        synopsis = "",
        img = "",
        runtime = "",
        status = "",
        favourite = false,
        genre = "",
        seasons = [],
    }) {
        this.uuid = uuid
        this.name = name
        this.year = year
        this.rating = rating
        this.synopsis = synopsis
        this.img = img
        this.runtime = runtime
        this.status = status
        this.favourite = favourite
        this.genre = genre
        this.seasons = seasons
    }
}

export default Show