import { defineConfig } from "vite";
import { svelte } from "@sveltejs/vite-plugin-svelte";
import routify from "@roxi/routify/vite-plugin";
import tailwindcss from "@tailwindcss/vite";

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
                target: "http://10.10.39.145:3000",
                changeOrigin: true,
            },
            "/socket.io": {
                target: "http://10.10.39.145:3000",
                ws: true,
                changeOrigin: true,
            },
        },
    },
});
