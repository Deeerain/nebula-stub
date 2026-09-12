import { createApp } from 'vue'
import { createPinia } from 'pinia'
import './style.css'
import axios from 'axios'
import App from './App.vue'

axios.defaults.baseURL = "/api/"

const pinia = createPinia()
const app = createApp(App)

app.use(pinia)

app.mount('#app')
