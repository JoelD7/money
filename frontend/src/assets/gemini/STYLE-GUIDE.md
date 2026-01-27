# Gemini Project Style Guide

This document defines the design system and styling standards for the Money Frontend project. It serves as the source of truth for the Gemini CLI when generating code, demos, or UI components.

## 1. Design Philosophy

*   **Atomic & Modular:** All UI elements should be built as reusable, atomic components.
*   **Clean & Trustworthy:** As a financial application, the design must prioritize clarity, readability, and a sense of security.
*   **Modern Hybrid:** The project leverages the utility-first speed of **Tailwind CSS** alongside the robust, accessible components of **Material UI (MUI)**.

## 2. Color Palette

The color system is defined in `tailwind.config.js` and `src/assets/colors/index.ts`.

### Primary Brand Colors (Green)
Used for primary actions, active states, and positive financial indicators (Income, Savings).

*   **Green 100 (Primary):** `#009821` (Tailwind: `bg-green-100`, `text-green-100`)
*   **Green 200 (Hover):** `#00851d`
*   **Green 300 (Deep/Text):** `#005212`
*   **Dark Green:** `#024511` (Used for high-contrast text or backgrounds)

### Secondary Colors (Blue)
Used for information, navigation, and secondary actions.

*   **Blue 100:** `#0088FE`
*   **Blue 200:** `#006dcc`
*   **Blue 300:** `#004d99`

### Functional Colors
*   **Error/Expense (Red):**
    *   `#D90707` (Red 100)
    *   `#ad0101` (Red 200)
*   **Warning (Orange):**
    *   `#FF8042` (Orange 100)

### Neutrals (Grays & Whites)
Used for backgrounds, borders, and text.

*   **Background (Gray):** `#F4F4F5` (Common app background)
*   **Surface (White):** `#FFFFFF`
*   **Text (Gray Dark):** `#6F6F6F`
*   **Text (Gray Darker):** `#4D4D4D`

## 3. Typography

*   **Font Family:** `Outfit`, sans-serif.
    *   *Note: Ensure this font is loaded in `index.html` or via CSS.*
*   **Weights:**
    *   Regular (400)
    *   Medium (500)
    *   Bold (700)

## 4. UI Components & Patterns

### Cards
Cards are the primary container for content.
*   **Background:** White (`bg-white`)
*   **Border:** Light Gray (`border border-slate-200` or `border-gray-200`)
*   **Radius:** Rounded XL (`rounded-xl` or `12px`)
*   **Shadow:** Soft shadow (`shadow-sm` or `shadow-md`)
*   **Padding:** Generous (`p-4` to `p-6`)

### Buttons
*   **Shape:** Rounded XL (`rounded-xl`).
*   **Primary:** Green background, White text.
*   **Secondary:** White background, Gray border, Gray text.
*   **Interaction:** Active scale effect (`active:scale-95`) and hover transitions (`transition-all`).

### Icons
*   **Library:** **FontAwesome** (`@fortawesome/react-fontawesome`) or **MUI Icons**.
*   *Note: Do not use Lucide icons in production code unless specifically added to `package.json`.*

## 5. Implementation Guidelines (Gemini CLI)

### Generating Demos (HTML/CSS)
When creating standalone HTML demos (e.g., in `src/assets/gemini/samples/`):
1.  Use **Tailwind CSS** (via CDN).
2.  Use the color hex codes defined above (e.g., `bg-[#009821]` or extend the tailwind config in the script).
3.  Use `Inter` or `Roboto` if `Outfit` is not easily available via CDN, but prefer `Outfit`.

### Generating React Components
1.  **Tailwind First:** Use Tailwind utility classes for layout, spacing, and typography.
2.  **MUI Integration:** Use MUI components for complex interactive elements (DataGrids, DatePickers) but style them to match the project's Tailwind theme.
3.  **Atomic Structure:** Place new components in `src/components/atoms`, `molecules`, or `organisms` appropriately.

## 6. CSS Variables Reference (Legacy Support)
If encountering legacy CSS, map these variables to the Tailwind equivalents:
*   `--green` -> `#009821`
*   `--gray` -> `#F4F4F5`
*   `--white` -> `#ffffff`
