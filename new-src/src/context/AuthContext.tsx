import React, { createContext, useState, useContext, useEffect, useCallback } from 'react';
import { authService } from '../api/auth';
import { AuthTokens, UserProfile, LoginPayload, RegisterPayload } from '../types/api';

interface AuthContextType {
  user: UserProfile | null;
  tokens: AuthTokens | null;
  login: (credentials: LoginPayload) => Promise<void>;
  logout: () => void;
  register: (payload: RegisterPayload) => Promise<void>;
  isAuthenticated: boolean;
}

const AuthContext = createContext<AuthContextType | null>(null);



export const AuthProvider = ({ children }) => {
  // CORRECTED: This is now a proper useState hook call.
  const [user, setUser] = useState<UserProfile | null>(null);
  const [tokens, setTokens] = useState<AuthTokens | null>(() => {
    try {
      const storedTokens = localStorage.getItem('tokens');
      return storedTokens ? JSON.parse(storedTokens) : null;
    } catch (error) {
      console.error("Error parsing tokens from localStorage", error);
      return null;
    }
  });

  const logout = useCallback(async () => {
    try {
      if (tokens?.refresh_token) {
        await authService.logout();
      }
    } catch (err) {
      console.error("Logout API call failed:", err);
    } finally {
      setUser(null);
      setTokens(null);
      localStorage.removeItem('tokens');
      // api.defaults.headers.common['Authorization'] is not directly accessible here, 
      // but the interceptor will handle the absence of tokens.
    }
  }, [tokens]);

  useEffect(() => {
    if (tokens?.access_token && !user) { // Only fetch if we have a token but no user object yet
      // api.defaults.headers.common['Authorization'] = `Bearer ${tokens.access_token}`;
      authService.getProfile()
        .then(response => {
          setUser(response.data);
        })
        .catch(() => {
          console.error("Failed to fetch profile with stored token, logging out.");
          logout();
        });
    }
  }, [tokens, user, logout]);

  const login = async (credentials: LoginPayload) => {
    const response = await authService.login(credentials);
    const { access_token, refresh_token } = response.data;
    
    const newTokens: AuthTokens = { access_token, refresh_token };
    
    // Fetch profile right after login to populate user object
    const profileResponse = await authService.getProfile();
    setUser(profileResponse.data);

    // Set tokens in state and storage *after* everything succeeds
    setTokens(newTokens);
    localStorage.setItem('tokens', JSON.stringify(newTokens));
  };

  const register = async (payload: RegisterPayload) => {
    // This function will simply make the API call. 
    // It will throw an error on failure, which the UI will catch.
    // On success, it does nothing; the UI will handle redirecting the user.
    await authService.register(payload);
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