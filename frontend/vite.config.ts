import { defineConfig, loadEnv } from "vite";
import { svelte } from "@sveltejs/vite-plugin-svelte";
import routify from "@roxi/routify/vite-plugin";
import tailwindcss from "@tailwindcss/vite";

const env = loadEnv(process.env.NODE_ENV as string, process.cwd(), "VITE_");

export default defineConfig({
    plugins: [
        routify({
            /* config */
        }),
        tailwindcss(),
        svelte(),
    ],
    server: {
        host: true,
        proxy: {
            "/api": {
                target: env.VITE_BACKEND_URL || "http://localhost:3000",
                changeOrigin: true,
            },
            "/socket.io": {
                target: env.VITE_BACKEND_URL || "http://localhost:3000",
                ws: true,
                changeOrigin: true,
            },
        },
    },
});
