
import api from './api';
import { ApiResponse, Ticket } from '../types/api';

export const ticketService = {
  purchaseTicket: async (ticket: Omit<Ticket, 'id' | 'status'>): Promise<ApiResponse<Ticket>> => {
    const response = await api.post<ApiResponse<Ticket>>('/tickets', ticket);
    return response.data;
  },

  getTicketById: async (id: string): Promise<ApiResponse<Ticket>> => {
    const response = await api.get<ApiResponse<Ticket>>(`/tickets/${id}`);
    return response.data;
  },

  updateTicketStatus: async (id: string, status: 'purchased' | 'cancelled' | 'used'): Promise<ApiResponse<Ticket>> => {
    const response = await api.patch<ApiResponse<Ticket>>(`/tickets/${id}/status`, { status });
    return response.data;
  },

  getAllTickets: async (): Promise<ApiResponse<Ticket[]>> => {
    const response = await api.get<ApiResponse<Ticket[]>>('/tickets');
    return response.data;
  },

  deleteTicket: async (id: string): Promise<ApiResponse<any>> => {
    const response = await api.delete<ApiResponse<any>>(`/tickets/${id}`);
    return response.data;
  },
};
