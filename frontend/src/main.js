import { createApp } from 'vue'
import './style.scss'
import 'vue-color-kit/dist/vue-color-kit.css'
import App from './App.vue'
import Toast, { POSITION } from "vue-toastification";
import "vue-toastification/dist/index.css";

createApp(App).use(Toast, {
  position: POSITION.BOTTOM_RIGHT
}).mount('#app')
