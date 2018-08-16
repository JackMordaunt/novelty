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
