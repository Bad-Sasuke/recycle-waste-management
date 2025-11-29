import '../src/assets/main.css'
import type { Preview } from '@storybook/vue3-vite'
import '../src/assets/main.css'

import { setup } from '@storybook/vue3'
import { createPinia } from 'pinia'
import { createRouter, createWebHistory } from 'vue-router'
import i18n from '../src/i18n'

const pinia = createPinia()

const router = createRouter({
  history: createWebHistory(),
  routes: [
    { path: '/', component: { template: '<div>Home</div>' } },
    { path: '/about', component: { template: '<div>About</div>' } },
    { path: '/login', component: { template: '<div>Login</div>' } },
    // Add other dummy routes as needed to prevent warnings
  ],
})

setup((app) => {
  app.use(pinia)
  app.use(i18n)
  app.use(router)
})

const preview: Preview = {
  parameters: {
    controls: {
      matchers: {
        color: /(background|color)$/i,
        date: /Date$/i,
      },
    },
  },
}

export default preview
