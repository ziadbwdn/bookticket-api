
import { Event } from './event';

export interface Ticket {
  id: string;
  eventId: string;
  userId: string;
  purchaseDate: string;
  event?: Event;
}
