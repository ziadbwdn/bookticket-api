import React from 'react';
import { Navigate, Outlet } from 'react-router-dom';
import { useAuth } from '../../context/AuthContext';

const ProtectedRoute = () => {
  const { isAuthenticated, user } = useAuth();

  // The 'user' check is important. It means we've successfully fetched the user profile.
  // 'isAuthenticated' alone might be true for a split second before a bad token is rejected.
  const isAuthReady = user !== null;

  if (!isAuthenticated && !isAuthReady) {
    // Could also show a loading spinner here while the initial auth check is happening
    return <Navigate to="/login" replace />;
  }

  // Outlet is a placeholder that will be replaced by the actual component
  // for the matched child route (e.g., DashboardPage).
  return <Outlet />;
};

export default ProtectedRoute;
