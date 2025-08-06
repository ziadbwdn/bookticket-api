
***

### 2. `planning.md` - The Refactoring Strategy

This document outlines the "what" and "why" of the refactor. It defines the problems, sets clear goals, and establishes the technical principles you will follow.

**Template for `planning.md`:**

```markdown
# Ticketing App Refactoring Plan

## 1. Current State Analysis

### Problem Statement

The React/TypeScript frontend application can be run, but it fails to interact correctly with the Go backend. Key functionalities such as user login, ticket fetching, and ticket creation are broken. The current state management is inconsistent, and there is no centralized logic for handling API requests, leading to duplicated code and poor error handling.

### Known Issues

- **[Issue 1, e.g., Login Failure]:** Submitting the login form results in a network error or an unhandled exception. The client does not correctly process the JWT from the backend.
- **[Issue 2, e.g., Ticket List Empty]:** The ticket list page loads but never displays tickets, likely due to a failed API call or incorrect state update.
- **[Issue 3, e.g., Inconsistent Typing]:** API responses are treated as `any`, nullifying the benefits of TypeScript and causing potential runtime errors.
- **[Issue 4, e.g., No Error Feedback]:** When API calls fail, the user is not shown any error message; the UI simply hangs or shows an empty state.

## 2. Refactoring Goals (The "To-Be" State)

### Primary Objectives

1.  **Establish Reliable API Communication:** All frontend API calls to the Go backend must be reliable, type-safe, and handle success and error states gracefully.
2.  **Implement Predictable State Management:** Centralize application state. Server cache and remote data (e.g., tickets, user info) should be managed separately from transient UI state.
3.  **Enforce Full Type Safety:** Ensure all data flowing from the backend, through state management, and into components is fully typed with TypeScript. Eliminate all uses of `any`.
4.  **Create a Clear, Maintainable Code Structure:** Organize files and components logically. Isolate side effects and create a clear data flow.

### Non-Goals (Out of Scope)

-   Adding any new features (e.g., ticket filtering, user profiles).
-   Making significant changes to the visual UI/UX design.
-   Migrating the database or changing the backend language/framework.

## 3. Refactoring Strategy & Key Decisions

1.  **API Layer:**
    -   We will use **Axios** for all HTTP requests.
    -   A single, centralized API client instance will be created. This instance will be configured with the base URL from environment variables and will handle attaching the authorization token (JWT) to headers for all protected requests.
    -   A structured error handling interceptor will be implemented to standardize API error management.

2.  **Server State Management:**
    -   We will use **TanStack Query (React Query)** to manage server state. This will handle all data fetching, caching, re-fetching, and mutation logic for tickets and other backend resources. This eliminates manual loading/error state management in components.

3.  **Client State Management:**
    -   We will use **React Context API with `useReducer`** for managing global client state, specifically user authentication status. This is lightweight and sufficient for our needs.

4.  **Code Style & Quality:**
    -   We will enforce **ESLint** and **Prettier** with a strict ruleset. A pre-commit hook will be added to ensure all committed code is formatted and linted.
