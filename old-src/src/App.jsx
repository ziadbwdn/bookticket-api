import React from 'react';
import { BrowserRouter, Routes, Route } from 'react-router-dom';
import { AuthProvider } from './context/AuthContext.jsx';
import ProtectedRoute from './components/layout/ProtectedRoute.jsx';

// Import our page components
import LoginPage from './pages/LoginPage.jsx';
import DashboardPage from './pages/DashboardPage.jsx';
// --> 1. MAKE SURE THIS IMPORT IS HERE <--
import RegisterPage from './pages/RegisterPage.jsx'; 

function App() {
  return (
    <AuthProvider>
      <BrowserRouter>
        <Routes>
          {/* Public Routes */}
          <Route path="/login" element={<LoginPage />} />
          
          {/* --> 2. MAKE SURE THIS ROUTE DEFINITION IS HERE <-- */}
          <Route path="/register" element={<RegisterPage />} />
          
          {/* Protected Routes - all wrapped in the ProtectedRoute component */}
          <Route element={<ProtectedRoute />}>
            <Route path="/" element={<DashboardPage />} />
          </Route>
          
          {/* A catch-all route for 404 Not Found pages */}
          <Route path="*" element={<h1>404 - Page Not Found</h1>} />
        </Routes>
      </BrowserRouter>
    </AuthProvider>
  );
}

export default App;