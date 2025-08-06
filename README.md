---

# 🎫 React Ticketing System - Frontend Application

**Project Status: UNDERGOING MAJOR REFACTORING**

> **Note:** This project is being actively refactored. The original implementation was built with a different framework, and this documentation describes the **target architecture** using **React, TypeScript, and TanStack Query**. The primary goal is to create a stable, maintainable, and scalable frontend that communicates reliably with the Go backend API.

## 🚀 Target Features

This application provides a comprehensive ticketing system. The feature set remains the same, but the implementation is being rebuilt on a new foundation.

### 👤 User Management
- **Authentication System**: Secure login/register with JWT tokens via the Go backend.
- **Role-Based Access**: Admin and User roles with different permissions, enforced by the backend.
- **Profile Management**: Users can edit their own profile information.
- **Session Management**: Secure token handling with refresh logic via the `AuthContext`.

### 🎭 Event Management
- **Event Browsing**: View all available events fetched from the server.
- **Event Details**: Comprehensive event information pages.
- **Real-time Updates**: Data is kept fresh via TanStack Query's caching and re-fetching mechanisms.

### 🎟️ Ticket Management
- **Ticket Purchase**: Streamlined buying process that communicates with the `/api/tickets` endpoint.
- **Ticket History**: Users can view and manage their purchased tickets.

### 📊 Admin Dashboard
- **Event CRUD**: A management interface for admins to create, update, and delete events.
- **User Management**: Monitor and manage user accounts (Admin).
- **Analytics Dashboard**: Sales reports and performance metrics fetched from the `/api/reports` endpoints.

## 🛠️ Target Technology Stack

This refactor establishes a modern, consistent React technology stack.

### Core Technologies
- **[React](https://react.dev/)** - The core UI library for building the application.
- **[TypeScript](https://www.typescriptlang.org/)** - For a fully type-safe codebase, from API to components.
- **[Vite](https://vitejs.dev/)** - The build tool and development server, providing a fast developer experience.
- **[Tailwind CSS](https://tailwindcss.com/)** - A utility-first CSS framework for rapid and consistent styling.

### Libraries & Tools
- **[React Router](https://reactrouter.com/)** - For all client-side routing and navigation.
- **[TanStack Query (React Query)](https://tanstack.com/query/latest)** - For all server state management: data fetching, caching, and mutations.
- **[Axios](https://axios-http.com/)** - The primary HTTP client for interacting with the Go backend API.
- **[React Context API](https://react.dev/learn/passing-data-deeply-with-context)** - For managing global client state like authentication status.
- **[Lucide React](https://lucide.dev/)** - For a clean and consistent icon library.
- **[Chart.js](https://www.chartjs.org/)** - For data visualization on the admin dashboard.
- **[date-fns](https://date-fns.org/)** - For robust date manipulation.

### Development Tools
- **[ESLint](https://eslint.org/)** - For enforcing code quality and style.
- **[Prettier](https://prettier.io/)** - For automated code formatting.
- **[Vitest](https://vitest.dev/)** - For unit and component testing within the Vite ecosystem.
- **[Playwright](https://playwright.dev/)** - For robust end-to-end testing of user flows.

## 📁 Target Project Structure

The project will follow a standard feature-oriented structure for React applications.

```
src/
├── components/              # Shared, reusable UI components (Button, Input, Modal)
├── config/                  # Application configuration (e.g., environment variables)
├── contexts/                # React Context providers (e.g., AuthProvider.tsx)
├── features/                # Feature-based modules (e.g., events, dashboard)
│   ├── auth/                # Login/Register components, hooks, services
│   └── events/              # Event list, event details, etc.
├── hooks/                   # Custom, reusable React hooks
├── lib/                     # External library configurations (e.g., apiClient.ts for Axios)
├── pages/                   # Top-level page components for routing
├── services/                # Type-safe functions that call the API (e.g., eventService.ts)
├── types/                   # All shared TypeScript type definitions (event.ts, user.ts)
├── utils/                   # General utility functions
├── App.tsx                  # Main application component with router setup
└── main.tsx                 # Application entry point
```

## 🚀 Getting Started (Refactor Guide)

### Prerequisites
- Node.js (v18 or higher)
- A running instance of the Go backend API server.

### Installation

1.  **Clone the repository**
    ```bash
    git clone <repository-url>
    cd <project-folder>
    ```

2.  **Install dependencies**
    ```bash
    npm install
    ```

3.  **Set up environment variables**
    Create a `.env` file from the example:
    ```bash
    cp .env.example .env
    ```
    Configure the `.env` file to point to your backend:
    ```env
    VITE_API_BASE_URL=http://localhost:8080/api
    VITE_APP_NAME="Ticketing System"
    ```

4.  **Start the development server**
    ```bash
    npm run dev
    ```
    Navigate to `http://localhost:5173` in your browser.

## 🔧 Available Scripts

- `npm run dev` - Starts the Vite development server.
- `npm run build` - Builds the application for production.
- `npm run preview` - Serves the production build locally for testing.
- `npm run test` - Runs unit tests with Vitest.
- `npm run lint` - Runs ESLint across the codebase.
- `npm run format` - Formats all code with Prettier.

## 🏗️ Target Architecture

### State Management
The application uses a clear separation of state concerns:

- **Server State:** Managed by **TanStack Query**. This handles all fetching, caching, synchronization, and updating of data from the Go backend. All API data (events, tickets, reports) lives here. This eliminates the need for manual loading and error state management in components.
- **Client State:** Managed by **React Context** and local component state (`useState`). This is used for global state like user authentication (`AuthContext`) and transient UI state (e.g., whether a modal is open).

### API Integration
All interaction with the RESTful API is centralized through a service layer that uses a shared **Axios** client instance. This ensures that setting the base URL and attaching JWT authorization headers is done consistently.

### Routing
**React Router** is used to manage all client-side routes, including public-only routes, protected routes that require authentication, and dynamic routes (e.g., `/events/:id`).

## 🔐 Security
- **Route Protection**: A custom `<ProtectedRoute>` component uses the `AuthContext` to prevent unauthenticated access to pages like the dashboard.
- **Secure Token Handling**: JWTs received from the backend are stored securely in memory via the `AuthContext` and are never exposed insecurely.
- **XSS Prevention**: React's JSX automatically escapes content, providing strong protection against XSS attacks.

## 📖 API Documentation

The frontend consumes the following endpoints from the Go backend:

### Authentication
- `POST /api/auth/login`
- `POST /api/auth/register`
- `GET /api/auth/profile`

### Events
- `GET /api/events`
- `GET /api/events/:id`

### Tickets (Protected)
- `POST /api/tickets`
- `GET /api/tickets`
- `GET /api/tickets/:id`

### Reports (Admin)
- `GET /api/reports/summary`

---

**This project is being refactored with ❤️ using React, TypeScript, and Vite.**