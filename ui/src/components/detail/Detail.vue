<template>
    <b-container fluid>
        <b-row id="show-description">
            <b-col>
                <b-container fluid>
                    <b-row>
                        <div class="background"
                             :style="backgroundImage">
                        </div>
                        <b-col sm="4" md="3" lg="2" >
                            <div class="cover-img">
                                <img :src="cover"/>
                            </div>
                        </b-col>
                        <b-col sm="7" md="8" lg="9" >
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
                        </b-col>
                    </b-row>
                </b-container>
            </b-col>
        </b-row>
        <b-row id="episode-browser">
            <b-col sm="3" md="3" lg="2">
                <div class="season-list">
                    <h2>Seasons</h2>
                    <b-list-group flush>
                        <b-list-group-item v-for="(season, ii) in seasons"
                                           :key="ii"
                                           :active="isActiveSeason(ii)"
                                           @click="setActiveSeason(ii)">
                            <div class="season">
                                Season {{ii+1}}
                            </div>
                        </b-list-group-item>
                    </b-list-group>
                </div>
            </b-col>
            <b-col sm="4" md="5" lg="6">
                <div class="episode-list">
                    <h2>Episodes</h2>
                    <b-list-group flush>
                        <b-list-group-item v-for="(ep, kk) in season.episodes"
                                           :key="ep.name"
                                           :active="isActiveEpisode(kk)"
                                           @click="setActiveEpisode(kk)">
                            <div class="episode">
                                Episode {{kk+1}}
                            </div>
                        </b-list-group-item>
                    </b-list-group>
                </div>
            </b-col>
            <b-col sm="5" md="4">
                <div class="episode-description">
                    <h2>{{episode.name}}</h2>
                    <div class="number">
                        Season {{activeSeason+1}}, Episode {{activeEpisode + 1}}
                    </div>
                    <div class="date">
                        Aired Date: {{episode.releaseDate.toLocaleString()}}
                    </div>
                    <p>{{episode.synopsis}}</p>
                    <div class="controls">
                        <div id="health"></div>
                        <div id="watch"></div>
                    </div>
                </div>
            </b-col>
        </b-row>
    </b-container>
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

#show-description {
    padding: 0;
    overflow: hidden;
    height: 220px;
}

#episode-browser {
    height: calc(100vh - 220px);
    width: 100vw;
    overflow: auto;
}

.season-list {
    height: 100%;
    padding-top: 20px;
}

.episode-list {
    height: 100%;
    padding-top: 20px;
}

.episode-description {
    height: 100%;
    padding-top: 20px;
}

.episode-description > p {
    margin-top: 10px;
    color: white;
}

.episode-description .controls {
    position: relative;
    height: 100%;
    width: 100%;
}

#watch {
    position: absolute;
    background: green;
    bottom: 10px;
    left: 10px;
    width: 100px;
    height: 100px;
}

.list-group {
    overflow: auto;
}

h2 {
    font-size: 2em;
    color: white;
}

.list-group-item {
    border: none;
    color: white;
    background-color: black;
    margin-top: 1px;
    margin-bottom: 1px;
}

.list-group-item:nth-child(odd) {
    background-color: #131313;
}

.list-group-item:hover {
    border: none;
    transition: 0.1s ease;
    background-color: #232323;
}

.list-group-item.active {
    border-left: 2px solid yellow;
    background-color: #232323;
}

.cover-img {
    margin-top: 10px;
    margin-bottom: 10px;
    height: 200px;
    width: 130px;
    margin-left: auto;
    margin-right: auto;
}

.cover-img > img {
    width: 100%;
    height: 100%;
    box-shadow: 0 0 10px black;
    border-radius: 3px;
}

.meta {
    color: white;
    height: 100%;
}

.infos {
    overflow: auto;
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
    bottom: 0;
    width: 100vw;
    height: 100%;
    min-height: 200px;
    filter: blur(60px) brightness(0.7);
}

.synopsis {
    padding-top: 10px;
    overflow: auto;
}
</style>
