# Refactoring Task List: Landing Page & Dashboard

This task list is ordered to establish a stable foundation before building UI components.

## Phase 1: Foundation & API Contract

-   [ ] **Task 1.1:** Create a new `refactor` git branch and ensure the dev environment (`npm run dev`) runs both services correctly.
-   [ ] **Task 1.2: [CRITICAL]** Based on the Go API, create the following TypeScript files:
    -   `src/types/event.ts`: Defines the `Event` interface.
    -   `src/types/ticket.ts`: Defines the `Ticket` interface.
    -   `src/types/user.ts`: Defines the `UserProfile` interface.
    -   `src/types/report.ts`: Defines the `SummaryReport` interface.
-   [ ] **Task 1.3:** Install dependencies: `axios @tanstack/react-query`.
-   [ ] **Task 1.4:** Create the central `apiClient` (Axios) instance and the `AuthContext`.

## Phase 2: API Service & Authentication

-   [ ] **Task 2.1:** Create `src/services/authService.ts` with a `login` function that calls `POST /api/auth/login`.
-   [ ] **Task 2.2:** Refactor the `LoginPage.tsx` to use the `login` function. On success, store the token in `AuthContext` and redirect to `/dashboard`.
-   [ ] **Task 2.3:** Implement a `<ProtectedRoute>` component that checks `AuthContext` and redirects to `/login` if no user is authenticated.

## Phase 3: View Implementation (Landing Page & Dashboard)

-   [ ] **Task 3.1 (Landing Page):**
    -   [ ] Create `src/services/eventService.ts` with a `getEvents` function for `GET /api/events`.
    -   [ ] In `LandingPage.tsx`, use a `useQuery({ queryKey: ['events'], queryFn: getEvents })` hook.
    -   [ ] Render the list of events from the hook's data. Handle `isLoading` and `error` states.
-   [ ] **Task 3.2 (Dashboard):**
    -   [ ] Create `src/services/userService.ts` with `getUserProfile` (`GET /api/auth/profile`) and `getUserTickets` (`GET /api/tickets`).
    -   [ ] Create `src/services/reportService.ts` with `getSummaryReport` (`GET /api/reports/summary`).
    -   [ ] In `DashboardPage.tsx`, use **parallel `useQuery` hooks** to fetch profile, tickets, and summary report data.
    -   [ ] Create and assemble dashboard components (`<WelcomeHeader>`, `<MyTicketsList>`, `<SummaryMetrics>`) that consume the data from these hooks.

## Phase 4: Finalization & Review

-   [ ] **Task 4.1:** Manually test the full user flow: Visit landing page -> Login -> View Dashboard -> Logout.
-   [ ] **Task 4.2:** Review browser console and network tab to ensure no errors and that API calls are being cached correctly by TanStack Query.
-   [ ] **Task 4.3:** Code review and merge.