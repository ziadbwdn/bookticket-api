# BookTicket API

This is a comprehensive ticketing application that allows users to book tickets for various events. The system provides functionalities for event management, user authentication, ticket purchasing, and activity logging.

## Features

- **Event Management:** Create, read, update, and delete events.
- **User Authentication:** User registration, login, and profile management with JWT-based authentication.
- **Ticket Booking:** Purchase tickets for events.
- **Activity Logging:** Track user activities for security and auditing purposes.
- **Reporting:** Generate summary reports for system activities and event-specific ticket reports.

## Getting Started

### Prerequisites

- Docker
- Docker Compose

### Installation

1. **Clone the repository:**
   ```bash
   git clone https://github.com/your-username/your-repository.git
   cd your-repository/app-ok
   ```

2. **Run the application:**
   ```bash
   docker-compose up --build
   ```

This will start the application and the database. The API will be accessible at `http://localhost:8080`.

## API Documentation

The API is documented using Swagger. You can access the Swagger UI at:

`http://localhost:8080/swagger/index.html`

For detailed information about each endpoint, please refer to the following documentation files:

- [Activity API](./docs/activity_api.md)
- [Event API](./docs/event_api.md)
- [Report API](./docs/report_api.md)
- [Ticket API](./docs/ticket_api.md)
- [User and Authentication API](./docs/user_auth_api.md)

## Environment Variables

The following environment variables can be set to configure the application:

- `DB_USER`: The username for the database.
- `DB_PASSWORD`: The password for the database.
- `DB_HOST`: The host of the database.
- `DB_PORT`: The port of the database.
- `DB_NAME`: The name of the database.
- `JWT_SECRET`: The secret key for JWT.
- `PORT`: The port on which the application will run.