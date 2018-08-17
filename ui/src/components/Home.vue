<template>
    <b-container fluid>
        <b-row>
            <b-col>
                <Nav/>
            </b-col>
        </b-row>
        <b-row id="show-browser"
               :visible="active == null">
            <b-col>
                <ul>
                    <li v-for="show in shows"
                        :key="show.uuid"
                        class="show"
                        @click="showDescription(show.uuid)">
                        <Tile :name="show.name"
                              :year="show.year"
                              :seasons="show.seasons"
                              :rating="show.rating"
                              :img="show.img"
                              :favourite="show.favourite"
                              @liked="toggleLiked(show.uuid)"
                        />
                    </li>
                </ul>
            </b-col>
        </b-row>
        <div class="detail-view"
             :visible="active != null">
            <div class="close"
                 @click="closeDescription">
                <icon name="times" scale="2"></icon>
            </div>
            <Detail v-if="active != null"
                    :name="active.name"
                    :year="active.year"
                    :synopsis="active.synopsis"
                    :runtime="active.runtime"
                    :status="active.status"
                    :genre="active.genre"
                    :rating="active.rating"
                    :cover="active.img"
                    :seasons="active.seasons"
            />
        </div>
    </b-container>
</template>

<script>
import Vue from "vue"
import Nav from "./Nav"
import Tile from "./Tile"
import Detail from "./detail/Detail"
import Show from "../models/Show"

export default {
    data() {
        return {
            shows: [],
            active: null
        }
    },
    created() {
        // TODO use api call instead of mocked fake data.
        this.shows = generateShows()
    },
    methods: {
        // TODO use api call instead of mocked fake data.
        toggleLiked(uuid) {
            let ii = uuid
            let show = this.shows[ii]
            show.favourite = !show.favourite
            Vue.set(this.shows, ii, show)
        },
        // TODO use api call instead of mocked fake data.
        showDescription(uuid) {
            let ii = uuid
            let show = this.shows[ii]
            this.active = show
        },
        closeDescription() {
            this.active = null
        }
    },
    components: {
        Nav,
        Tile,
        Detail
    }
}

const pad = number => {
    if (number < 10) {
        return "0" + number
    } else {
        return number
    }
}

import Season from "../models/Season"
import Episode from "../models/Episode"

const generateSeasons = () => {
  let seasons = []
  for (let ii = 0; ii < Math.floor(Math.random() * 8 + 2); ii++) {
    let episodes = []
    for (let kk = 0; kk < Math.floor(Math.random() * 12 + 3); kk++) {
        episodes.push(
            new Episode({
                name: `Episode ${kk}`,
                runtime: Math.floor(Math.random() * 80 + 20),
                synopsis:
                    "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum."
            })
        )
    }
    seasons.push(
        new Season({
            episodes,
        })
    )
  }
  return seasons
}

const generateShows = () => {
    let shows = []
    let genres = [
        "drama",
        "science fiction",
        "action",
        "horror",
        "fantasy",
        "documentary"
    ]
    for (let ii = 0; ii < 50; ii++) {
        shows.push(
            new Show({
                uuid: ii,
                name: "TV Show",
                year: "20" + pad(Math.floor(Math.random() * 20)),
                rating: Math.floor(Math.random() * 10),
                img: `https://source.unsplash.com/random/3${Math.floor(Math.random() * 99)}x6${Math.floor(Math.random() * 99)}`,
                synopsis: "Suits follows college drop-out Mike Ross, who accidentally lands a job with one of New York's best legal closers, Harvey Specter. They soon become a winning team with Mike's raw talent and photographic memory, and Mike soon reminds Harvey of why he went into the field of law in the first place.",
                favourite: Math.random() < 0.5,
                genre: genres[Math.floor(Math.random() * (genres.length - 1))],
                status: Math.random() < 0.5 ? "Returning Series" : "Cancelled",
                runtime: `${Math.floor(Math.random() * 100)} min`,
                seasons: generateSeasons(),
            })
        )
    }
    return shows
}
</script>

<style scoped>
#show-browser {
    display: none;
}

#show-browser[visible] {
    display: block;
}

.show {
    display: block;
    margin: 10px;
    float: left;
    width: 134px;
}

.detail-view {
    display: none;
    position: fixed;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    width: 100vw;
    height: 100vh;
    background-color: #101012;
}

.detail-view[visible] {
    display: block;
    z-index: 99;
}

.detail-view .close {
    position: absolute;
    right: 20px;
    top: 20px;
    color: red !important;
    z-index: 100;
    opacity: 0.5 !important;
}


</style>