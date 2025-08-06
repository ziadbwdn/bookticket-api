
import api from './api';
import { AuthTokens, LoginPayload, RegisterPayload, UserProfile, ApiResponse } from '../types/api';

export const authService = {
  register: async (payload: RegisterPayload): Promise<ApiResponse<AuthTokens>> => {
    const response = await api.post<ApiResponse<AuthTokens>>('/auth/register', payload);
    return response.data;
  },

  login: async (payload: LoginPayload): Promise<ApiResponse<AuthTokens>> => {
    const response = await api.post<ApiResponse<AuthTokens>>('/auth/login', payload);
    return response.data;
  },

  refreshToken: async (refresh_token: string): Promise<ApiResponse<AuthTokens>> => {
    const response = await api.post<ApiResponse<AuthTokens>>('/auth/refresh', { refresh_token });
    return response.data;
  },

  logout: async (): Promise<ApiResponse<any>> => {
    const response = await api.post<ApiResponse<any>>('/auth/logout');
    return response.data;
  },

  getProfile: async (): Promise<ApiResponse<UserProfile>> => {
    const response = await api.get<ApiResponse<UserProfile>>('/auth/profile');
    return response.data;
  },

  updateProfile: async (profile: Partial<UserProfile>): Promise<ApiResponse<UserProfile>> => {
    const response = await api.put<ApiResponse<UserProfile>>('/auth/profile', profile);
    return response.data;
  },
};
