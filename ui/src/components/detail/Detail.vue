<template>
    <div id="detail">
        <section id="show-description">
            <div class="background"
                 :style="backgroundImage">
            </div>
            <div class="cover-img">
                <img :src="cover"/>
            </div>
            <div class="meta">
                <h1>{{name}}</h1>
                <div class="infos">
                    <p>{{year}}</p>
                    <p>{{runtime}}</p>
                    <p>{{status}}</p>
                    <p>{{genre}}</p>
                    <p>{{rating}}/10</p>
                </div>
                <p class="synopsis">
                    {{synopsis}}
                </p>
            </div>
        </section>
        <section id="episode-browser">
            <div class="season-list">
                <h2>Seasons</h2>
                <ul class="scrollable list-group">
                    <li v-for="(season, ii) in seasons"
                        :key="ii"
                        :active="isActiveSeason(ii)"
                        @click="setActiveSeason(ii)"
                        class="list-group-item"
                    >
                        <div class="season">
                            Season {{ii+1}}
                        </div>
                    </li>
                </ul>
            </div>
            <div class="episode-list">
                <h2>Episodes</h2>
                <ul class="scrollable list-group">
                    <li v-for="(ep, kk) in season.episodes"
                        :key="ep.name"
                        :active="isActiveEpisode(kk)"
                        @click="setActiveEpisode(kk)"
                        class="list-group-item"
                    >
                        <div class="episode">
                            Episode {{kk+1}}
                        </div>
                    </li>
                </ul>
            </div>
            <div class="episode-description">
                <h2>{{episode.name}}</h2>
                <div class="number">
                    Season {{activeSeason+1}}, Episode {{activeEpisode + 1}}
                </div>
                <div class="date">
                    Aired Date: {{episode.releaseDate.toLocaleString()}}
                </div>
                <p class="episode-synopsis">
                    {{episode.synopsis}}
                </p>
            </div>
        </section>
    </div>
</template>

<script>
import Season from '../../models/Season'
import Episode from '../../models/Episode'

export default {
    props: {
        name: String,
        year: String,
        cover: String,
        synopsis: String,
        runtime: String,
        genre: String,
        rating: Number,
        status: String,
        seasons: Array,
    },
    data() {
        return {
            activeSeason: 0,
            activeEpisode: 0,
        }
    },
    computed: {
        backgroundImage() {
            return `background-image: url(${this.cover})`
        },
        season() {
            return this.seasons[this.activeSeason] || new Season()
        },
        episode() {
            return this.season.episodes[this.activeEpisode] || new Episode()
        }
    },
    methods: {
        setActiveSeason(ii) {
            this.activeSeason = ii
            this.activeEpisode = 0
        },
        isActiveSeason(ii) {
            return this.activeSeason === ii
        },
        setActiveEpisode(ii) {
            this.activeEpisode = ii
        },
        isActiveEpisode(ii) {
            return this.activeEpisode === ii
        }
    }
}
</script>

<style scoped>
h1 {
    font-size: 3em;
    margin: 0;
}

h2 {
    font-size: 2em;
    color: white;
    text-overflow: ellipsis;
    white-space: nowrap;
}

#detail {
    height: 100%;
    width: 100%;
}

#show-description {
    display: flex;
    flex-direction: row;
    overflow: hidden;
    height: 220px;
    padding: 10px;
}

#episode-browser {
    display: flex;
    flex-direction: row;
    height: 100%;
    max-height: calc(100vh - 220px);
    width: 100vw;
    background-color: black;
    padding: 10px;
}

#episode-browser > div {
    margin: 5px;
}

#episode-browser > div:nth-child(1) {
    flex-grow: 1;
}

#episode-browser > div:nth-child(2) {
    flex-grow: 2;
}

#episode-browser > div:nth-child(3) {
    flex-grow: 1;
}

.episode-description {
    max-width: 33vw;
    color: white;
}

.episode-synopsis {
    margin-top: 5px;
    padding-right: 5px;
    height: calc(100% - 90px);
    overflow-y: auto;
}

.scrollable {
    height: calc(100% - 40px);
    overflow-y: auto;
}

.cover-img {
    flex-grow: 1;
    height: 200px;
    min-width: 134px;
    max-width: 134px;
}

.cover-img > img {
    width: 100%;
    height: 100%;
    box-shadow: 0 0 10px black;
    border-radius: 3px;
}

.meta {
    height: 100%;
    flex-grow: 2;
    min-width: 200px;
    color: white;
    padding-left: 10px;
}

.synopsis {
    margin-top: 10px;
    overflow-x: hidden;
    overflow-y: auto;
    max-height: 90px;
}

.infos {
    overflow-x: auto;
    white-space: nowrap;
    text-overflow: ellipsis;
}

.infos > p {
    display: inline;
    padding: 0;
    margin: 0;
}

.infos > *:not(:first-child)::before {
    content: "•";
    margin-right: 20px;
    margin-left: 20px;
}

.background {
    position: absolute;
    background-size: cover;
    background-position: center;
    background-repeat: no-repeat;
    background-color: black;
    top: 0;
    left: 0;
    right: 0;
    width: 100vw;
    height: 220px;
    filter: blur(60px) brightness(0.7);
    z-index: -1;
}

.list-group-item {
    border-radius: 0px !important;
    background-color: black;
    color: white;
    margin-top: 1px;
    margin-bottom: 1px;
}

.list-group-item:nth-child(odd) {
    background-color: #131313;
}

.list-group-item:hover {
    transition: 0.1s ease;
    background-color: #232323;
    border-color: #232323;
}

.list-group-item[active] {
    box-shadow: 0;
    border-left: 1px solid  yellow;
}
</style>
