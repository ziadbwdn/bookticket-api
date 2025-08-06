---

### 3. `planning.md` (Adjusted)

```markdown
# Ticketing App Refactoring Plan

## 1. Current State Analysis

### Problem Statement

The React client and Go backend are failing to communicate correctly. The frontend does not reliably consume the API defined by the Gin router, leading to broken features on the landing page and dashboard.

### Known Issues & Root Causes

- **API Contract Mismatch:** Frontend API calls do not correctly match the required paths, methods, or JSON structures of the Go backend (e.g., `GET /api/events`).
- **Missing or Incorrect Authentication:** The frontend fails to properly store and send the JWT, causing requests to protected endpoints (`/profile`, `/tickets`, `/reports`) to fail with 401 Unauthorized errors.
- **Lack of Type Safety:** The frontend treats API responses as `any`, preventing TypeScript from catching bugs related to mismatched Go structs and client interfaces.

## 2. Refactoring Goals

1.  **Achieve Full API Compliance:** The primary goal is to make the React client a **perfectly compliant consumer** of the provided Go API.
2.  **Implement the App Design Context:** The refactor will strictly follow the principles and user journeys outlined in `context.md`.
3.  **Deliver a Functional Landing Page & Dashboard:** These two views must be fully functional and stable by the end of the refactor.

### Non-Goals

-   Changing any backend API logic or database schemas.
-   Implementing features beyond the landing page, dashboard, and their supporting details (login, event details).

## 3. Technical Refactoring Strategy

1.  **API Contract:** We will create TypeScript `type` definitions in `src/types/` that **exactly mirror the JSON output of the Go backend's structs**.
2.  **API Service Layer:** We will use **Axios** to create a central `apiClient`. This client will be configured to read the JWT from an `AuthContext`.
3.  **Server State Management:** We will use **TanStack Query** to manage all data fetching, caching, and mutation, eliminating manual loading/error state management in components.