
import { useQuery } from '@tanstack/react-query';
import { getEvents } from '../services/eventService';

const LandingPage = () => {
  const { data: events, isLoading, error } = useQuery({ queryKey: ['events'], queryFn: getEvents });

  if (isLoading) return <div>Loading...</div>;
  if (error) return <div>An error occurred: {error.message}</div>;

  return (
    <div>
      <h1>Upcoming Events</h1>
      <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fill, minmax(300px, 1fr))', gap: '20px' }}>
        {events?.map((event) => (
          <div key={event.id} style={{ border: '1px solid #ccc', borderRadius: '8px', padding: '20px' }}>
            <h2>{event.name}</h2>
            <p>{event.date}</p>
            <p>{event.description}</p>
          </div>
        ))}
      </div>
    </div>
  );
};

export default LandingPage;
