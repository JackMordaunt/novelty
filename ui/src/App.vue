<style>

/* Removing arbitrary text highlighting makes the app feel more native */ 
html {
  -webkit-touch-callout: none; /* iOS Safari */
    -webkit-user-select: none; /* Safari */
     -khtml-user-select: none; /* Konqueror HTML */
       -moz-user-select: none; /* Firefox */
        -ms-user-select: none; /* Internet Explorer/Edge */
            user-select: none; /* Non-prefixed version, currently
                                  supported by Chrome and Opera */
}

html, body {
    width: 100vw;
    height: 100vh;
    overflow: hidden;
    padding: 0;
    margin: 0;
}

body > div {
    height: 100vh;
    overflow: scroll;
    -webkit-overflow-scrolling: touch;
}

#app {
    font-family: "Avenir", Helvetica, Arial, sans-serif;
    -webkit-font-smoothing: antialiased;
    -moz-osx-font-smoothing: grayscale;
    color: #2c3e50;
    min-width: 100%;
    min-height: 100%;
    background-color: black;
}
</style>

<template>
    <div id="app">
        <Home/>
    </div>
</template>

<script>
import Home from "./components/Home.vue";

export default {
    name: "app",
    components: {
        Home
    }
};
</script>
