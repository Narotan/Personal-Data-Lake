import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'
import path from 'path'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [react()],
  resolve: {
    alias: {
      "@": path.resolve(__dirname, "./src"),
    },
  },
  server: {
    port: 8000,
    host: '0.0.0.0',
    strictPort: true,
    hmr: {
      clientPort: 80,
      host: 'localhost',
    },
    watch: {
      usePolling: true,
    },
    proxy: {
      '/api': {
        target: 'http://app:8080',
        changeOrigin: true,
        secure: false,
      }
    }
  }
})
