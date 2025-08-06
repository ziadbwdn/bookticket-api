# Refactoring Task List

This list is ordered to minimize conflicts and build upon a stable foundation.

## Phase 1: Foundation and Setup (The "Make it Run" Phase)

-   [ ] **Task 1.1:** Create a new `refactor-2025` git branch from `main`. All work will be done here.
-   [ ] **Task 1.2:** Validate, document, and fix the complete local development setup as described in the updated `readme.md`. Ensure `npm run dev` works reliably.
-   [ ] **Task 1.3:** Setup ESLint and Prettier with a strict, shared configuration. Run it across the entire frontend codebase to establish a consistent format baseline.
-   [ ] **Task 1.4:** Define TypeScript `interface` or `type` for all known API resources (e.g., `User`, `Ticket`, `AuthResponse`) in a dedicated `src/types` folder. Align these with the Go backend models.

## Phase 2: API Layer Refactor (The "Make it Talk" Phase)

-   [ ] **Task 2.1:** Install Axios and TanStack Query (`@tanstack/react-query`).
-   [ ] **Task 2.2:** Create a centralized Axios instance in `src/services/apiClient.ts` configured with the base URL.
-   [ ] **Task 2.3:** Implement the authentication flow (`login`, `logout`) in a dedicated `src/services/authService.ts` file using the new Axios client.
-   [ ] **Task 2.4:** Create an `AuthContext` to store the user's authentication state and JWT, and wrap the entire application in an `AuthProvider`.

## Phase 3: Component and View Refactor (The "Make it Work" Phase)

-   [ ] **Task 3.1 (Login Page):** Refactor the login page to use the new `authService`. On successful login, it should update the `AuthContext` and redirect the user. Implement proper loading and error feedback.
-   [ ] **Task 3.2 (Ticket List Page):**
    -   [ ] Create a TanStack Query hook (`useTickets`) in `src/features/tickets/hooks/useTickets.ts` that fetches data from the `GET /tickets` endpoint.
    -   [ ] Refactor the ticket list component to use this hook.
    -   [ ] Correctly display loading, error, and empty states based on the hook's return values.
-   [ ] **Task 3.3 (Create Ticket Page):**
    -   [ ] Create a TanStack Query mutation hook (`useCreateTicket`) for the `POST /tickets` endpoint.
    -   [ ] Refactor the "create ticket" form to use this mutation.
    -   [ ] On successful creation, invalidate the `useTickets` query to automatically refresh the ticket list.
-   [ ] **Task 3.4 (Protected Routes):** Implement a component that checks the `AuthContext` and redirects any unauthenticated users away from protected pages (like the ticket list) to the login page.

## Phase 4: Cleanup and Validation

-   [ ] **Task 4.1:** Search the codebase for any remaining direct `fetch()` calls and replace them.
-   [ ] **Task 4.2:** Search the codebase for any remaining uses of `any` and replace them with proper types.
-   [ ] **Task 4.3:** Perform a full manual test of all application features (login, logout, view tickets, create ticket).
-   [ ] **Task 4.4:** Review the browser's console and network tab during testing to ensure there are no errors.
-   [ ] **Task 4.5:** Merge the `refactor-2025` branch back into `main` after a final review.
