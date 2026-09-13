import { createApp } from 'vue'
import { createPinia } from 'pinia'
import './style.css'
import axios from 'axios'
import App from './App.vue'
import { initTelegramMock } from './telegram-mock.js'

if (import.meta.env.DEV) {
    console.warn(import.meta.env.DEV, "initing Telegram mock")
    initTelegramMock()
}



axios.defaults.baseURL = "/api/"

const pinia = createPinia()
const app = createApp(App)

app.use(pinia)

app.mount('#app')
