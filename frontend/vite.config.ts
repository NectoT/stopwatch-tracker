import {defineConfig} from 'vite'
import {svelte} from '@sveltejs/vite-plugin-svelte'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [svelte()],
  resolve: {
    alias: {
      "@wails": "/wailsjs",
      "@assets": "/src/assets",
      "@components": "/src/components",
      "@App": "/src/App.svelte",
    }
  }
})
