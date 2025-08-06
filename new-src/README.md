# Borehole Data Frontend Application

A modern, responsive frontend application for managing borehole data, built with React and TypeScript. This application handles complex business logic on the client side while interfacing with a RESTful API backend.

## 🚀 Features

### 👤 User Management
- **Authentication System**: Secure login/register with JWT tokens
- **Profile Management**: User profile editing and account management
- **Session Management**: Auto-refresh tokens and session timeout handling

### 📊 Dashboard
- Displays a personalized dashboard for authenticated users.

## 🛠️ Technology Stack

### Core Technologies
- **[React](https://react.dev/)** - A JavaScript library for building user interfaces
- **[TypeScript](https://www.typescriptlang.org/)** - Type-safe JavaScript
- **[Vite](https://vitejs.dev/)** - Fast build tool and dev server

### Libraries & Tools
- **[React Router](https://reactrouter.com/)** - Declarative routing for React
- **[Axios](https://axios-http.com/)** - Promise-based HTTP client for the browser and Node.js
- **[ESLint](https://eslint.org/)** - Pluggable JavaScript linter

## 📁 Project Structure

```
new-src/
├── public/
│   └── vite.svg
├── src/
│   ├── api/                 # API service integrations (TypeScript)
│   │   ├── api.ts
│   │   ├── auth.ts
│   │   ├── event.ts
│   │   └── ticket.ts
│   ├── assets/
│   │   └── react.svg
│   ├── components/
│   │   ├── common/
│   │   └── layout/
│   │       └── ProtectedRoute.tsx
│   ├── context/
│   │   └── AuthContext.tsx
│   ├── features/
│   ├── hooks/
│   ├── pages/               # Page-level components (TypeScript)
│   │   ├── DashboardPage.tsx
│   │   ├── LoginPage.tsx
│   │   ├── ProjectDetailPage.jsx
│   │   ├── ProjectsListPage.jsx
│   │   └── RegisterPage.tsx
│   ├── types/               # TypeScript type definitions
│   │   └── api.ts
│   ├── App.css
│   ├── App.tsx
│   ├── index.css
│   └── main.tsx
├── .gitignore
├── eslint.config.js
├── index.html
├── package.json
├── package-lock.json
├── README.md
├── SUMMARY.md
├── tsconfig.json
├── tsconfig.node.json
└── vite.config.js
```

## 🚀 Getting Started

### Prerequisites
- Node.js (v18 or higher)
- npm or yarn
- Backend API server running (refer to `APIGuideline.md` for API endpoints)

### Installation

1.  **Navigate to the `new-src` directory**
    ```bash
    cd /home/ziad_bwdn/gemini-cli/app-ok/client/new-src
    ```

2.  **Install dependencies**
    ```bash
    npm install
    ```

3.  **Set up environment variables**
    Create a `.env` file in the `new-src` directory based on your backend API URL. For example:
    ```env
    VITE_API_URL=http://localhost:3000/api
    ```

4.  **Start development server**
    ```bash
    npm run dev
    ```

5.  **Open your browser**
    Navigate to `http://localhost:5173` (or the port indicated by Vite).

## 📖 API Documentation

Refer to `APIGuideline.md` in the parent directory for a summary of the API routes and their specifications.

## ✅ TypeScript Conversion

This project has been refactored to use TypeScript for improved code quality, maintainability, and type safety, as outlined in `PLAN.md` and `TASK.md`.