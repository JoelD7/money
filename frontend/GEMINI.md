# Project Gemini Context: Money Frontend

## Overview
This is the frontend application for the "Money" project, built with React, TypeScript, and Vite. It serves as a personal finance management tool.

## Tech Stack

### Core
- **Runtime/Build:** Vite, TypeScript
- **Framework:** React 18
- **Styling:** Tailwind CSS, Material UI (MUI), Emotion

### State Management & Data Fetching
- **Server State:** TanStack Query (React Query)
- **Client State:** Redux Toolkit (with Redux Persist)
- **Routing:** TanStack Router (File-based routing)

### Utilities
- **Forms:** Formik
- **Validation:** Yup, Zod
- **HTTP Client:** Axios (with `axios-auth-refresh`)
- **Date Handling:** Dayjs
- **Charts:** Recharts
- **Icons:** FontAwesome

## Architecture & Patterns

### Directory Structure (`src/`)
- **`api/`**: Raw API service definitions. Contains functions that directly call endpoints using Axios.
- **`assets/`**: Static assets, global styles, and color definitions.
- **`components/`**: UI components following **Atomic Design** principles:
  - **`atoms/`**: Basic building blocks (Buttons, Inputs).
  - **`molecules/`**: Simple combinations of atoms (Cards, Form fields).
  - **`organisms/`**: Complex sections (Tables, Charts, Navbars).
  - **`templates/`**: Page layouts.
- **`dev/`**: Development utilities (previews, palettes).
- **`pages/`**: Top-level page components (likely legacy or wrapped by routes).
- **`queries/`**: React Query hooks. This is the primary layer for data interaction in components.
- **`routes/`**: Route definitions for TanStack Router.
- **`store/`**: Redux store setup and slices (Auth, User).
- **`types/`**: TypeScript type definitions (Domain models, API responses).
- **`utils/`**: Helper functions and formatters.

### Key Conventions
1.  **Atomic Design**: Components are strictly categorized by complexity.
2.  **Query Separation**: API calls are defined in `api/` but consumed via hooks in `queries/`.
3.  **Routing**: The project uses TanStack Router. Routes are defined in `routes/` and generated into `routeTree.gen.ts`.
4.  **Styling**: Uses a mix of Tailwind CSS utility classes and MUI components.

## Development Commands

- **Start Dev Server:** `npm run dev` (or `bun run dev`)
- **Build:** `npm run build`
- **Lint:** `npm run lint`
- **Test:** `npm run test`

## Important Notes
- **Authentication**: Handled via `store/auth.ts` and `api/auth.ts`, likely using JWTs with refresh token logic (`axios-auth-refresh`).
- **File Naming**: CamelCase for non-component files, PascalCase for components.
