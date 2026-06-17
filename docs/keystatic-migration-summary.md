# Keystatic & Astro Fonts Migration Summary

## Overview
Replaced the previous Sanity CMS architecture with Keystatic to manage localized content (English and Vietnamese) locally in the Git repository. Also implemented local font preloading using Astro's built-in Fonts API.

## What was implemented

### 1. Dependencies and Configuration
*   **Upgraded Astro:** Updated Astro from `v5.0.0` to the latest `v6.x` to fully support the stable Astro Fonts API (`fontProviders` exported from `astro/config`).
*   **Keystatic Installed:** Added `@keystatic/core`, `@keystatic/astro`, `@astrojs/react`, `react`, and `react-dom`.
*   **`astro.config.mjs` Updated:**
    *   Added the React and Keystatic integrations.
    *   Configured the Astro Fonts API (under `fonts: []` in the root config object) to preload the local `Phudu` font variants.

### 2. Font Asset Migration
*   Copied the missing `.woff2` font files from the old `@temp/apps/web/src/assets/fonts/` directory to `src/assets/fonts/` to resolve the build errors.

### 3. Keystatic Schema Definition (`keystatic.config.ts`)
*   **GitHub Storage:** Configured Keystatic to dynamically switch between `local` storage during development (`process.env.NODE_ENV === 'development'`) and `github` storage (`dinhdkhoa/little-da-vinci`) for production. This enables the admin panel to function flawlessly when hosted statically on GitHub Pages by directly communicating with the GitHub API.
*   Created a configuration that closely mirrors the old Sanity schema, particularly the field-level localization strategy (e.g., `title.en` and `title.vi` grouped using `fields.object()`).
*   **Singletons:** Defined schemas for **Site Settings** and **Home Page** matching the structure of the existing `little-da-vinci` UI components.
*   **Collections:** Defined schemas for **Authors**, **Categories**, and **Blog Posts**.
*   **Cloudinary Support:** Implemented a Cloudinary fallback structure (`secure_url` and `public_id`) so media URLs can be pasted from Cloudinary and used seamlessly in Astro components, bypassing the need for a complex custom widget.

### 4. Content Integration & Refactoring
*   **Default Content:** Created the initial Keystatic JSON files at `src/content/home.json` and `src/content/settings.json` based on the previously hardcoded text in the `.astro` components. This ensured the site compiled perfectly on the first run.
*   **Astro Components Updated:** Refactored `index.astro` (and `en/index.astro`) to load the JSON data and pass it down as props. The UI components (`Hero`, `HowWeTeach`, `Programs`, `Gallery`, `Testimonials`, `CTA`) were updated to consume the `data` and `lang` props dynamically.

## How to access
*   **Admin Dashboard:** Run `pnpm run dev` and navigate to `http://127.0.0.1:4321/little-da-vinci/keystatic`.
*   **Content Location:** All managed content is saved as JSON/Markdown files directly inside the `src/content/` directory.