import Vue from 'vue'
import ElementUI from 'element-ui';
import 'element-ui/lib/theme-chalk/index.css';
import App from './App.vue'
import axios from 'axios';
Vue.prototype.$axios = axios //全局注册，使用方法为:this.$axios
Vue.use(ElementUI);
Vue.config.productionTip = false

new Vue({
  el: '#app',
  render: h => h(App)
});
