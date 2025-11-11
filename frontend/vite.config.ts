import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

// https://vite.dev/config/
export default defineConfig({
    plugins: [react()],
    build: {
        rollupOptions: {
            input: {
                // Multiple indexing to help with
                // static hosting
                main: 'index.html',
                // login: 'login/index.html',
            },
        },
    },
})
