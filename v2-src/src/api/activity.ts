
import api from './api';
import { ApiResponse, UserActivity } from '../types/api';

export const activityService = {
  logActivity: async (activity: Omit<UserActivity, 'id' | 'timestamp'>): Promise<ApiResponse<UserActivity>> => {
    const response = await api.post<ApiResponse<UserActivity>>('/activities', activity);
    return response.data;
  },

  listActivities: async (): Promise<ApiResponse<UserActivity[]>> => {
    const response = await api.get<ApiResponse<UserActivity[]>>('/activities');
    return response.data;
  },

  getActivitySummary: async (userId: string): Promise<ApiResponse<any>> => {
    const response = await api.get<ApiResponse<any>>(`/activities/summary/${userId}`);
    return response.data;
  },

  getSecurityAlerts: async (): Promise<ApiResponse<any>> => {
    const response = await api.get<ApiResponse<any>>('/activities/alerts');
    return response.data;
  },

  cleanOldActivities: async (): Promise<ApiResponse<any>> => {
    const response = await api.delete<ApiResponse<any>>('/activities');
    return response.data;
  },
};
