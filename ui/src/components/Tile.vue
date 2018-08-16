<template>
<div>
    <div class="cover">
        <img :src="img"/>
        <div class="filter"></div>
        <div class="overlay">
            <div @click.stop="liked()">
                <icon   name="heart"
                        scale="1.5"
                        class="like"
                        :liked="favourite"
                ></icon>
            </div>
            <div class="rating">
                <span class="rating-stars">
                    <div v-for="(icon, ii) in stars" :key="ii" class="star">
                        <icon v-if="icon.half">
                            <icon name="star-half-alt" :class="icon.class"></icon>
                        </icon>
                        <icon v-else name="star" :class="icon.class"></icon>
                    </div>
                </span>
                <span class="rating-value">
                    {{rating}}/10
                </span>
            </div>
        </div>
    </div>
    <div class="meta">
        <p class="name">{{name}}</p>
        <p class="year">{{year}}</p>
        <p   v-if="seasons > 0"
                class="seasons">
            {{seasons}} Seasons
        </p>
    </div>
</div>
</template>

<script>
export default {
    props: {
        name: String,
        year: String,
        img: String,
        tags: Array,
        favourite: Boolean,
        seasons: Array,
        rating: Number,
    },
    computed: {
        stars() {
            let stars = []
            let rating = this.rating/2
            for (let ii = 0; ii < Math.floor(rating); ii++) {
                stars.push({
                    half: false,
                    class: "colored",
                })
            }
            if (rating % 1 > 0) {
                stars.push({
                    half: true,
                    class: "colored",
                })
            }
            for (let ii = Math.ceil(rating); ii < 5; ii++) {
                stars.push({
                    half: false,
                    class: "",
                })
            }
            return stars
        }
    },
    methods: {
        liked() {
            this.$emit("liked")
        },
    }
}
</script>

<style scoped>

.cover {
    height: 200px;
    position: relative;
    color: white;
    background-color: #0E0E17;
    overflow: hidden;
}

.cover > img {
    width: 100%;
    height: auto;
}

.overlay {
    position: absolute;
    top: 0;
    visibility: hidden;
    min-width: 100%;
    min-height: 100%;
    z-index: 2;
}

.cover:hover {
    transition: 0.15s ease;
    border: 2px solid yellow;
    background-color: black;
}

.cover:hover .overlay {
    transition: 0.3s ease;
    opacity: 1;
    visibility: visible;
    z-index: 3;
}

.filter {
    visibility: hidden;
    background-color: black;
    opacity: 0;
}

.cover:hover .filter {
    position: absolute;
    width: 100%;
    height: 100%;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    transition: 0.3s ease;
    opacity: 0.3;
    visibility: visible;
    z-index: 2;
}

.meta {
    width: 100%;
    font-weight: bold;
    font-size: 0.8em;
}

.meta > .name {
    color: white;
    max-width: 100%;
    font-size: 1.3em;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
    margin: 2px 0 0;
}

.meta > .year {
    color: #979ca4;
    display: inline;
}

.meta > .seasons {
    color: #979ca4;
    float: right;
}

.like {
    z-index: 5;
    position: absolute;
    right: 8px;
    top: 8px;
    color: white;
}

.cover:hover .like:hover {
    transition: 1s ease;
    color: rgba(255, 0, 0, 0.3);
}

.cover:hover .like[liked]:hover {
    transition: 1s ease;
    color: red;
}

.like[liked] {
    color: red;
}

.rating {
    position: absolute;
    bottom: 0.5em;
    width: 100%;
    height: 1em;
    color: white;
    line-height: 1em;
    text-align: center;
    font-size: 1em;
    font-weight: bold;
}

.rating-stars {
    position: absolute;
    left: 5px;
    bottom: 1px;
}

.star {
    position: relative;
    display: inline;
    color: white;
}

.star .colored {
    color: gold;
}

.star svg {
    overflow: visible;
}

.rating-value {
    position: absolute;
    right: 5px;
}
</style>