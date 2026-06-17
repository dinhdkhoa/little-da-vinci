import { defineConfig, passthroughImageService, fontProviders } from "astro/config";
import tailwindcss from "@tailwindcss/vite";
import react from "@astrojs/react";
import keystatic from "@keystatic/astro";

export default defineConfig({
    site: 'https://dinhdkhoa.github.io',
    base: process.env.NODE_ENV === 'production' ? '/little-da-vinci' : undefined,
    prefetch: true,
    output: 'static',
    vite: {
        plugins: [tailwindcss()],
        envDir: "../../",
    },
    integrations: [
        react(),
        process.env.NODE_ENV === 'production' ? null : keystatic()
    ].filter(Boolean),
    image: {
        service: passthroughImageService()
    },
    fonts: [{
        name: "Phudu",
        provider: fontProviders.google(),
        cssVariable: "--font-phudu",
    }],
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
