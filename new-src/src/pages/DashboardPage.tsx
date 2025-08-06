import React from 'react';
import { useAuth } from '../context/AuthContext';
import { useNavigate } from 'react-router-dom';

const DashboardPage = () => {
  const { user, logout } = useAuth();
  const navigate = useNavigate();

  const handleLogout = () => {
    logout();
    // After logging out, redirect the user back to the login page
    navigate('/login');
  };

  // The ProtectedRoute ensures the user object is available,
  // but it's good practice to handle the brief moment it might be null.
  if (!user) {
    return <div>Loading user data...</div>;
  }

  return (
    <div style={{ padding: '20px' }}>
      <h1>Dashboard</h1>
      {/* The user object from /api/auth/profile might have 'name' or 'username' */}
      <h2>Welcome, {user.name || user.username}!</h2>
      
      
      <button onClick={handleLogout} style={{ padding: '10px 20px', cursor: 'pointer' }}>
        Logout
      </button>
    </div>
  );
};

export default DashboardPage;
