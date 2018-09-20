import Vue from 'vue'
import BootstrapVue from 'bootstrap-vue'
import App from './App.vue'
import Settings from "./settings"
import WS from "./services/Websocket"

// import 'bootstrap/dist/css/bootstrap.css'
import './styles/bootstrap.min.css'
import 'bootstrap-vue/dist/bootstrap-vue.css'

import 'vue-awesome/icons'
import Icon from 'vue-awesome/components/Icon'

Vue.use(BootstrapVue)
Vue.component("icon", Icon)
Vue.config.productionTip = false

Vue.prototype.$ws = new WS(`ws://${Settings.server_address}/ws`)

new Vue({
  render: h => h(App),
}).$mount('#app')
