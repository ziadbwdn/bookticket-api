
# API Guideline

This document provides a summary of the API routes for the application.

## Authentication

- `POST /api/auth/register`: Register a new user.
- `POST /api/auth/login`: Login a user.
- `POST /api/auth/refresh`: Refresh the authentication token.
- `POST /api/auth/logout`: Logout a user.
- `GET /api/auth/profile`: Get the user's profile.
- `PUT /api/auth/profile`: Update the user's profile.

## Events

- `POST /api/events`: Create a new event. (Admin only)
- `GET /api/events/:id`: Get an event by its ID.
- `GET /api/events`: Get all events.
- `PUT /api/events/:id`: Update an event. (Admin only)
- `DELETE /api/events/:id`: Delete an event. (Admin only)

## Tickets

- `POST /api/tickets`: Purchase a ticket.
- `GET /api/tickets/:id`: Get a ticket by its ID.
- `PATCH /api/tickets/:id/status`: Update the status of a ticket.
- `GET /api/tickets`: Get all tickets. (Admin only)
- `DELETE /api/tickets/:id`: Delete a ticket. (Admin only)

## Reports

- `GET /api/reports/summary`: Get a summary report. (Admin only)
- `GET /api/reports/events/:id`: Get a ticket event report. (Admin only)

## User Activities

- `POST /api/activities`: Log a user activity.
- `GET /api/activities`: List all user activities. (Admin only)
- `GET /api/activities/summary/:userID`: Get a summary of user activities. (Admin only)
- `GET /api/activities/alerts`: Get security alerts. (Admin only)
- `DELETE /api/activities`: Clean old user activities. (Admin only)

