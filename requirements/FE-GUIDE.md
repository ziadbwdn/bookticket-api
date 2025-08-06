# Application Design Context & Guidelines

This document outlines the architectural principles and user flow for the React frontend, ensuring it functions as a compliant and efficient client for the E-Commerce Go backend API.

## 1. Core Application Concept

The application is a client for the E-Commerce Go backend. It provides two primary user experiences based on roles defined by the backend.

-   **Unauthenticated Visitors:** Can browse all available menu items.
-   **Authenticated Customers:** Can manage their profile, add items to a shopping cart, place orders, and view their order history.
-   **Authenticated Admins:** Have all customer privileges, and can also manage menu items, manage all user orders, and view sales reports.

## 2. Guiding Principles

1.  **The Go Backend is the Single Source of Truth.** The React application **must not** contain its own business logic. Its role is to present data, collect user input, and send it to the correct API endpoint. All data validation (beyond simple form checks), permissions, and state calculations are handled by the server.

2.  **Server State vs. Client State.** We must strictly separate these two types of state:
    -   **Server State:** Data that lives on the server and is fetched via the API (e.g., menu items, user profiles, cart contents, orders, reports). This will be managed exclusively by a server state management library like **TanStack Query**.
    -   **Client State:** Data that exists only in the browser (e.g., contents of a form input, whether a modal is open, the current auth token). This will be managed by **React Context** and local component state (`useState`).

3.  **Data Flow is Unidirectional.** The flow of data for any feature should be clear and predictable:
    -   A React component dispatches an API call (via a TanStack Query hook).
    -   The API service sends the request to the Go backend.
    -   The backend processes the request and returns data or an error.
    -   TanStack Query caches the response and updates the component.
    -   The component re-renders with the new data, loading, or error state.

## 3. Key User Journeys & Required Endpoints

### Journey 1: Public Menu Browsing
-   **Goal:** Show a visitor all available food/drink items.
-   **View:** `MenuPage.tsx`
-   **Primary Endpoint:** `GET /api/menus`
-   **Data Flow:** The page loads, triggers a `useQuery` hook to fetch all menu items, and displays them.

### Journey 2: User Authentication
-   **Goal:** Authenticate a user and establish a persistent session.
-   **Views:** `LoginPage.tsx`, `RegisterPage.tsx`
-   **Primary Endpoints:** `POST /api/auth/login`, `POST /api/auth/register`
-   **Data Flow:** User submits credentials. A `useMutation` hook sends the data to the appropriate endpoint. On success, the returned JWT is stored in an `AuthContext`, and the `apiClient` will now automatically use this token for all future protected requests.

### Journey 3: Customer Shopping & Checkout
-   **Goal:** Allow a logged-in customer to add items to their cart and place an order.
-   **Views:** `MenuPage.tsx` (for adding items), `CartPage.tsx`, `CheckoutPage.tsx`
-   **Primary Endpoints:**
    -   `GET /api/cart`
    -   `POST /api/cart/items`
    -   `POST /api/orders/checkout`
-   **Data Flow:** From the menu, a customer triggers a `useMutation` to add an item to the cart. They can view their current cart, which uses a `useQuery` hook. Finally, they can "checkout," which triggers another `useMutation` to create a formal order from their cart.

### Journey 4: Admin Menu Management
-   **Goal:** Allow an admin to create, update, and delete menu items.
-   **View:** `AdminMenuPage.tsx` (Protected Admin Route)
-   **Primary Endpoints:**
    -   `POST /api/menus`
    -   `PUT /api/menus/:id`
    -   `DELETE /api/menus/:id`
-   **Data Flow:** The admin page will feature forms and buttons that trigger `useMutation` hooks for creating, updating, or deleting menu items. The list of menu items will be fetched with `useQuery`.

### Journey 5: Admin Order Management
-   **Goal:** Allow an admin to view all customer orders and update their statuses.
-   **View:** `AdminOrdersPage.tsx` (Protected Admin Route)
-   **Primary Endpoints:**
    -   `GET /api/orders`
    -   `PUT /api/orders/:id`
-   **Data Flow:** The admin fetches all orders with `useQuery`. They can then use a `useMutation` hook to change the status of any given order (e.g., "pending" to "shipped").
