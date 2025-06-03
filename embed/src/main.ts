import './assets/main.css'

import { createApp } from 'vue'
// import { createPinia } from 'pinia'
// import vue3GoogleLogin from 'vue3-google-login'

import App from './App.vue'
import router from './router'

const app = createApp(App)

// app.use(createPinia())
app.use(router)

// app.use(vue3GoogleLogin, {
//   clientId: '114782695264-qbhlmm64mf883aetb4l07tf4m7jv4ek1.apps.googleusercontent.com'
// })

app.mount('#app')
