# Project Color Scheme & Design System

This document outlines the semantic color tokens used in the ComputerXplorers project. These tokens are defined in `src/styles/global.css` within the Tailwind `@theme` block.

## Core Brand Colors

| Token Name | Color Code | Description |
| :--- | :--- | :--- |
| `primary` | `#00696c` | Main brand teal, used for headers, primary icons, and active states. |
| `secondary` | `#5c6300` | Olive green accent. |
| `tertiary` | `#466800` | Bright green accent. |
| `secondary-container` | `#Dbeb00` | Used for primary "Book Now" buttons. |

## Surface & UI Tokens

Use these tokens instead of generic colors (like `bg-white` or `bg-slate-x`) to ensure theme consistency.

### Backgrounds & Surfaces
- `surface-container-lowest` (`#ffffff`): Cleanest white for cards and dropdown backgrounds.
- `surface-container-low` (`#f5f3f3`): Secondary background / Hover states.
- `surface` (`#fbf9f8`): Main body/app background.
- `inverse-surface` (`#303030`): Dark mode surface background for components.

### Text & Icons
- `on-surface` (`#1b1c1c`): Primary text color.
- `on-surface-variant` (`#3d4949`): Secondary/de-emphasized text.
- `inverse-on-surface` (`#f2f0f0`): Text color for dark surfaces.
- `outline-variant` (`#bcc9c9`): Standard border color for separators/cards.

## Implementation Pattern in Astro/Tailwind

Always prefer semantic Tailwind classes provided by the theme:

```html
<!-- INCORRECT -->
<div class="bg-white text-slate-800 border-slate-100">...</div>

<!-- CORRECT -->
<div class="bg-surface-container-lowest text-on-surface border-outline-variant">...</div>
```
