
// src/types/api.ts

export interface AuthTokens {
  access_token: string;
  refresh_token: string;
}

export interface UserProfile {
  id: string;
  username: string;
  email: string;
  // Add other user profile fields as per your backend
}

export interface RegisterPayload {
  username: string;
  email: string;
  password: string;
}

export interface LoginPayload {
  email: string;
  password: string;
}

export interface RefreshTokenPayload {
  refresh_token: string;
}

export interface Event {
  id: string;
  name: string;
  description: string;
  date: string;
  location: string;
  // Add other event fields
}

export interface Ticket {
  id: string;
  eventId: string;
  userId: string;
  status: 'purchased' | 'cancelled' | 'used';
  // Add other ticket fields
}

export interface UserActivity {
  id: string;
  userId: string;
  action: string;
  timestamp: string;
  // Add other activity fields
}

export interface ApiResponse<T> {
  data: T;
  message: string;
  status: string;
}
