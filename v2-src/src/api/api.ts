import axios, { AxiosRequestConfig, AxiosResponse } from 'axios';
import { AuthTokens, ApiResponse, LoginPayload, RegisterPayload, RefreshTokenPayload, UserProfile, Event, Ticket, UserActivity } from '../types/api';
import { authService } from './auth';


const api = axios.create({
  baseURL: import.meta.env.VITE_API_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Request Interceptor: Adds the auth token to every outgoing request
api.interceptors.request.use(
  (config: AxiosRequestConfig) => {
    const tokens: AuthTokens | null = JSON.parse(localStorage.getItem('tokens') || 'null');
    if (tokens?.access_token) {
      if (!config.headers) {
        config.headers = {};
      }
      config.headers['Authorization'] = `Bearer ${tokens.access_token}`;
    }
    return config;
  },
  (error) => Promise.reject(error)
);

// Response Interceptor: Handles token refresh on 401 errors
api.interceptors.response.use(
  (response: AxiosResponse) => response, // Simply return the response if it's successful
  async (error) => {
    const originalRequest = error.config;
    
    // Check if the error is a 401 and we haven't already retried the request
    if (error.response.status === 401 && !originalRequest._retry) {
      originalRequest._retry = true;

      const tokens: AuthTokens | null = JSON.parse(localStorage.getItem('tokens') || 'null');

      if (!tokens?.refresh_token) {
        window.location.href = '/login';
        return Promise.reject(error);
      }
      originalRequest._retry = true; // Mark that we've retried this request

      try {
        const localTokens = JSON.parse(localStorage.getItem('tokens'));
        if (!localTokens?.refresh_token) {
          // If no refresh token, we can't do anything.
          window.location.href = '/login';
          return Promise.reject(error);
        }

        // Call the refresh token endpoint
        const { data: newTokens } = await authService.refreshToken(localTokens.refresh_token);
        localStorage.setItem('tokens', JSON.stringify(newTokens));

        // Update the default header for subsequent requests
        api.defaults.headers.common['Authorization'] = `Bearer ${newTokens.access_token}`;
        // Update the header of the original failed request
        originalRequest.headers['Authorization'] = `Bearer ${newTokens.access_token}`;

        // Retry the original request with the new token
        return api(originalRequest);

      } catch (refreshError) {
        // If refresh fails, clear storage and redirect to login
        console.error("Session refresh failed. Redirecting to login.", refreshError);
        localStorage.removeItem('tokens');
        window.location.href = '/login';
        return Promise.reject(refreshError);
      }
    }

    // For any other errors, just reject the promise
    return Promise.reject(error);
  }
);

export default api;
