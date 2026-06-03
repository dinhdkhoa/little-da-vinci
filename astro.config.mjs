import { defineConfig, passthroughImageService } from "astro/config";
import tailwindcss from "@tailwindcss/vite";

export default defineConfig({
    site: 'https://dinhdkhoa.github.io',
    base: '/little-da-vinci',
    prefetch: true,
    vite: {
        plugins: [tailwindcss()],
        envDir: "../../",
    },
    image: {
        service: passthroughImageService()
    },
    i18n: {
        locales: ["en", "vi"],
        defaultLocale: "vi",
        routing: {
            prefixDefaultLocale: false,
        },
        fallback: {
            en: "vi",
        }
    },
    build: {
        inlineStylesheets: "always",
    },
});
