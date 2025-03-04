import flowbitePlugin from 'flowbite/plugin'
const { tailwindExtractor } = require("tailwindcss/lib/lib/purgeUnusedStyles");

export default {
    content: ['./src/**/*.{html,js,svelte,ts}', './node_modules/flowbite-svelte/**/*.{html,js,svelte,ts}'],
    darkMode: ["class", '[data-theme="dark"]'],
    theme: {
        extend: {
            fontFamily: {
                lexend: ['"Lexend Deca"', 'sans-serif'], // Add Lexend Deca here
            },
        },
    },
    purge: {
        content: [
          'src/app.html',
          'src/**/*.svelte',
        ],
        options: {
          defaultExtractor: (content) => [
            ...tailwindExtractor(content),
            ...[
              ...content.matchAll(/(?:class:)*([\w\d-/:%.]+)/gm)
            ].map(([_match, group, ..._rest]) => group),
          ],
          keyframes: true,
        },
    },
    plugins: [flowbitePlugin],
    corePlugins: {
        preflight: false,
    }
};
