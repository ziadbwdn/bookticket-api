import React, { createContext, useState, useContext, useEffect, useCallback } from 'react';
import api from '../api/api';

const AuthContext = createContext(null);

export const AuthProvider = ({ children }) => {
  // CORRECTED: This is now a proper useState hook call.
  const [user, setUser] = useState(null); 
  const [tokens, setTokens] = useState(() => {
    try {
      const storedTokens = localStorage.getItem('tokens');
      return storedTokens ? JSON.parse(storedTokens) : null;
    } catch (error) {
      console.error("Error parsing tokens from localStorage", error);
      return null;
    }
  });

  const logout = useCallback(() => {
    if (tokens?.refresh_token) {
      api.post('/auth/logout', { refresh_token: tokens.refresh_token })
         .catch(err => console.error("Logout API call failed:", err));
    }
    
    setUser(null);
    setTokens(null);
    localStorage.removeItem('tokens');
    delete api.defaults.headers.common['Authorization'];
  }, [tokens]);

  useEffect(() => {
    if (tokens?.access_token && !user) { // Only fetch if we have a token but no user object yet
      api.defaults.headers.common['Authorization'] = `Bearer ${tokens.access_token}`;
      api.get('/auth/profile')
        .then(response => {
          setUser(response.data.data);
        })
        .catch(() => {
          console.error("Failed to fetch profile with stored token, logging out.");
          logout();
        });
    }
  }, [tokens, user, logout]);

  const login = async (username, password) => {
    const response = await api.post('/auth/login', { username, password });
    const { access_token, refresh_token } = response.data.data;
    
    const newTokens = { access_token, refresh_token };
    api.defaults.headers.common['Authorization'] = `Bearer ${access_token}`;
    
    // Fetch profile right after login to populate user object
    const profileResponse = await api.get('/auth/profile');
    setUser(profileResponse.data.data);

    // Set tokens in state and storage *after* everything succeeds
    setTokens(newTokens);
    localStorage.setItem('tokens', JSON.stringify(newTokens));
  };

  const register = async (fullName, username, email, password, role) => {
    // This function will simply make the API call. 
    // It will throw an error on failure, which the UI will catch.
    // On success, it does nothing; the UI will handle redirecting the user.
    await api.post('/auth/register', {
      fullName,
      username,
      email,
      password,
      role,
    });
  };

  const authContextValue = {
    user,
    tokens,
    login,
    logout,
    register, // <-- Add register here
    isAuthenticated: !!user,
  };

  return (
    <AuthContext.Provider value={authContextValue}>
      {children}
    </AuthContext.Provider>
  );
};

export const useAuth = () => {
  return useContext(AuthContext);
};