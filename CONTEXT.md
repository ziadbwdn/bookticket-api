# Application Design Context & Guidelines

This document outlines the architectural principles and user flow for the React frontend, ensuring it functions as a compliant and efficient client for the Go backend API.

## 1. Core Application Concept

The application is an event ticketing platform. It serves two primary user roles:
- **Unauthenticated Visitors:** Can browse and view details of upcoming events.
- **Authenticated Users:** Can manage their profile, purchase tickets for events, view their purchased tickets, and (if they are organizers/admins) see summary reports.

## 2. Guiding Principles

1.  **The Go Backend is the Single Source of Truth.** The React application **must not** contain its own business logic. Its role is to present data, collect user input, and send it to the correct API endpoint. All data validation (beyond simple form checks), permissions, and state calculations are handled by the server.

2.  **Server State vs. Client State.** We must strictly separate these two types of state:
    - **Server State:** Data that lives on the server and is fetched via the API (e.g., events, tickets, user profiles, reports). This will be managed exclusively by **TanStack Query**.
    - **Client State:** Data that exists only in the browser (e.g., contents of a form input, whether a modal is open, the current auth token). This will be managed by **React Context** and local component state (`useState`).

3.  **Data Flow is Unidirectional.** The flow of data for any feature should be clear and predictable:
    *   A React component dispatches an API call (via a TanStack Query hook).
    *   The API service (e.g., Axios) sends the request to the Go backend.
    *   The backend processes the request and returns data or an error.
    *   TanStack Query caches the response and updates the component.
    *   The component re-renders with the new data, loading, or error state.

## 3. Key User Journeys & Required Endpoints

### Journey 1: The Public Landing Page

-   **Goal:** Show a visitor all available events.
-   **View:** `LandingPage.tsx`
-   **Primary Endpoint:** `GET /api/events`
-   **Data Flow:** The page loads, triggers a `useQuery` hook to fetch all events, and displays them in a grid or list. Each item links to its own detail page (`/events/:id`).

### Journey 2: User Login & Authentication

-   **Goal:** Authenticate a user and establish a persistent session.
-   **View:** `LoginPage.tsx`
-   **Primary Endpoint:** `POST /api/auth/login`
-   **Data Flow:** User submits credentials. A `useMutation` hook sends the data to the login endpoint. On success, the returned JWT is stored in the `AuthContext`, and the user is redirected to the dashboard. The `apiClient` will now automatically use this token for all future protected requests.

### Journey 3: The User Dashboard

-   **Goal:** Give a logged-in user a personalized, at-a-glance overview.
-   **View:** `DashboardPage.tsx` (Protected Route)
-   **Primary Endpoints (fetched in parallel):**
    -   `GET /api/auth/profile` to get the user's name/email for a welcome message.
    -   `GET /api/tickets` to get a list of the user's own tickets.
    -   `GET /api/reports/summary` to display high-level statistics (if the user is an admin/organizer).
    -   `GET /api/activities` to show a feed of recent relevant activities.
-   **Data Flow:** On page load, multiple `useQuery` hooks are triggered simultaneously. The dashboard displays loading states for each section and renders them as the data arrives.
