
import { useQuery } from '@tanstack/react-query';
import { getEvents } from '../services/eventService';
import { useAuth } from '../context/AuthContext';

export function DashboardPage() {
    const { data: events, isLoading, error } = useQuery({
        queryKey: ['events'],
        queryFn: getEvents,
    });
    const { logout } = useAuth();

    if (isLoading) return <div>Loading events...</div>;
    
    if (error) return <div>Error fetching events. Please try logging in again.</div>;

    return (
        <div style={{ padding: '20px' }}>
            <button onClick={logout} style={{ padding: '10px 20px', cursor: 'pointer' }}>
              Logout
            </button>
            <h1>Your Events</h1>
            <ul>
                {events?.map((event: any) => (
                    <li key={event.id}>{event.name}</li>
                ))}
            </ul>
        </div>
    );
}

export default DashboardPage;

