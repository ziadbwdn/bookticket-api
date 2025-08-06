
import api from './api';
import { ApiResponse, SummaryReport, TicketEventReport } from '../types/api';

export const reportService = {
  getSummaryReport: async (): Promise<ApiResponse<SummaryReport>> => {
    const response = await api.get<ApiResponse<SummaryReport>>('/reports/summary');
    return response.data;
  },

  getTicketEventReport: async (eventId: string): Promise<ApiResponse<TicketEventReport>> => {
    const response = await api.get<ApiResponse<TicketEventReport>>(`/reports/events/${eventId}`);
    return response.data;
  },
};
