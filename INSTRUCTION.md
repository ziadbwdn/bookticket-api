# Frontend Refactoring and Implementation Guide

This document provides explicit instructions to refactor the React frontend based on the definitive API contract provided by the Postman documentation. The goal is to resolve blocking issues and establish correct patterns for future development.

## 1. Core Architectural Correction: Landing Page and Routing

The Postman documentation confirms that `GET /api/events` is a **protected endpoint** requiring a Bearer Token. The initial concept of a public landing page showing events is therefore invalid.

**The primary public route for unauthenticated users must be the login page.**

### Action Plan: Reconfigure Application Routing

In your main router file (`src/App.tsx` or similar), you must implement the following logic:

1.  **Set `/login` as the primary public route.** A route for `/register` should also be public.
2.  The default path `/` should redirect to `/login` or render the `LoginPage`.
3.  Create a **protected route**, for example `/dashboard`, which will render the component that fetches and displays the event list.
4.  Implement a `<ProtectedRoute>` component that checks for an authentication token. If no token exists, it must redirect the user to `/login`.
5.  After a successful login, the application must programmatically redirect the user from `/login` to `/dashboard`.

**Example `App.tsx` Router Setup:**

```tsx
import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import { LoginPage } from './pages/LoginPage';
import { RegisterPage } from './pages/RegisterPage';
import { DashboardPage } from './pages/DashboardPage'; // This page will fetch and show events
import { ProtectedRoute } from './components/ProtectedRoute';

function App() {
  return (
    <BrowserRouter>
      <Routes>
        {/* Public Routes */}
        <Route path="/login" element={<LoginPage />} />
        <Route path="/register" element={<RegisterPage />} />

        {/* Protected Routes */}
        <Route 
          path="/dashboard"
          element={
            <ProtectedRoute>
              <DashboardPage />
            </ProtectedRoute>
          } 
        />

        {/* Default route redirects to login */}
        <Route path="*" element={<Navigate to="/login" />} />
      </Routes>
    </BrowserRouter>
  );
}
```

## 2. Resolving the `400 Bad Request` on Registration

The Postman documentation for `POST /api/auth/register` shows that the request body **requires** four fields: `username`, `email`, `password`, and `role`. The `400` error occurs because the `role` field is missing from the frontend payload.

### Action Plan: Add the `role` Field to the Registration Payload

Modify the `handleSubmit` function in `src/pages/RegisterPage.tsx`.

**Current (Incorrect) Code:**```tsx
const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    const payload = { username, email, password }; // Missing 'role'
    // ...
};
```

**Corrected Code:**
```tsx
// src/pages/RegisterPage.tsx

const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    
    // This payload matches the API contract.
    const payload = {
        username: username,
        email: email,
        password: password,
        role: "user" // Add the required 'role' field
    };

    try {
        await register(payload);
        // On success, redirect to login page
    } catch (error) {
        console.error("Registration failed:", error);
        // Display error message to the user
    }
};
```

## 3. Implementation Guide for All API Interaction

To ensure all other API calls work correctly, follow these patterns.

### Step 1: Centralize API Client and Authentication Logic

Your `apiClient` (Axios instance) must automatically attach the authentication token to all requests. This is done with an interceptor.

**`src/lib/apiClient.ts`:**
```typescript
import axios from 'axios';

// Assume you have a way to get the token, e.g., from localStorage or a state manager.
const getAuthToken = () => localStorage.getItem('authToken');

const apiClient = axios.create({
    baseURL: import.meta.env.VITE_API_BASE_URL,
});

// This interceptor adds the token to every request if it exists.
apiClient.interceptors.request.use(
    (config) => {
        const token = getAuthToken();
        if (token) {
            config.headers['Authorization'] = `Bearer ${token}`;
        }
        return config;
    },
    (error) => {
        return Promise.reject(error);
    }
);

export { apiClient };
```

### Step 2: Implement Data Fetching for Protected Routes

Use `TanStack Query` and your `apiClient` to fetch data for authenticated users. The `apiClient` will handle the token automatically.

**Example: Fetching Events for the Dashboard**

**`src/services/eventService.ts`:**
```typescript
import { apiClient } from '../lib/apiClient';
import { Event } from '../types/event'; // Assuming you have this type defined

export const getEvents = async (): Promise<Event[]> => {
    // This call will automatically include the Authorization header.
    const response = await apiClient.get('/events');
    return response.data;
};
```

**`src/pages/DashboardPage.tsx`:**
```tsx
import { useQuery } from '@tanstack/react-query';
import { getEvents } from '../services/eventService';

export function DashboardPage() {
    const { data: events, isLoading, error } = useQuery({
        queryKey: ['events'],
        queryFn: getEvents,
    });

    if (isLoading) return <div>Loading events...</div>;
    
    // TanStack Query will catch the 401 error if the token is bad/expired
    if (error) return <div>Error fetching events. Please try logging in again.</div>;

    return (
        <div>
            <h1>Your Events</h1>
            {/* ... render the events list ... */}
        </div>
    );
}
```

Follow this instruction
