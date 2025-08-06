
import api from './api';
import { ApiResponse } from '../types/api';
import { Event } from '../types/event';

export const eventService = {
  createEvent: async (event: Omit<Event, 'id'>): Promise<ApiResponse<Event>> => {
    const response = await api.post<ApiResponse<Event>>('/events', event);
    return response.data;
  },

  getEventById: async (id: string): Promise<ApiResponse<Event>> => {
    const response = await api.get<ApiResponse<Event>>(`/events/${id}`);
    return response.data;
  },

  getAllEvents: async (): Promise<ApiResponse<Event[]>> => {
    const response = await api.get<ApiResponse<Event[]>>('/events');
    return response.data;
  },

  updateEvent: async (id: string, event: Partial<Event>): Promise<ApiResponse<Event>> => {
    const response = await api.put<ApiResponse<Event>>(`/events/${id}`, event);
    return response.data;
  },

  deleteEvent: async (id: string): Promise<ApiResponse<any>> => {
    const response = await api.delete<ApiResponse<any>>(`/events/${id}`);
    return response.data;
  },
};
