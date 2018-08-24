/* eslint-disable */
import Vue from 'vue'
import BootstrapVue from 'bootstrap-vue'
import App from './App.vue'

// import 'bootstrap/dist/css/bootstrap.css'
import './styles/bootstrap.min.css'
import 'bootstrap-vue/dist/bootstrap-vue.css'

import 'vue-awesome/icons'
import Icon from 'vue-awesome/components/Icon'

Vue.use(BootstrapVue)
Vue.component("icon", Icon)
Vue.config.productionTip = false

new Vue({
  render: h => h(App),
}).$mount('#app')

import axios from "axios"

let http = axios.create({
  baseURL: "http://127.0.0.1",
  port: 9090,
})

axios.post("/open", {
  "torrent_path": "magnet:?xt=urn:btih:674D163D2184353CE21F3DE5196B0A6D7C2F9FC2&dn=bbb_sunflower_1080p_60fps_stereo_abl.mp4&tr=udp%3a%2f%2ftracker.openbittorrent.com%3a80%2fannounce&tr=udp%3a%2f%2ftracker.publicbt.com%3a80%2fannounce&ws=http%3a%2f%2fdistribution.bbb3d.renderfarming.net%2fvideo%2fmp4%2fbbb_sunflower_1080p_60fps_stereo_abl.mp4",
})

import Socket from "simple-websocket"

let ws = new Socket("ws://127.0.0.1:9090/1/stats")
ws.on("data", (data) => {
  console.log(data)
})

ws.on("error", (err) => {
  console.log(err)
})