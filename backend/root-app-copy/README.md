# **Back-End System: Study Case Ticketing**

## **Objectives**

1. Design and develop RESTful API using backend development best practices.

2. Implement clean project structure with layer separation (controller, service, repository, etc.).

3. Apply authentication and authorization to maintain API security.

4. Create input data validation to ensure system integrity.

5. Build complex CRUD features according to project scenario requirements.

6. Manage data relationships in relational databases (1-to-Many, Many-to-Many).

7. Create simple reports and data analysis from the built system.

8. Create clear and easily understandable API documentation.

9. Design code testing to ensure system stability and reliability.

## **Scenario**

**We are asked to build a RESTful API** that has complex features and covers several interconnected data management modules. The project theme focuses on medium to large scale systems, namely Ticketing System with strict authentication, authorization, and validation requirements.

This API aims to manage tickets for events, transportation, or other services. This system is often used in the entertainment industry (concerts, cinemas), transportation (planes, trains), or even internal company systems for issue tracking.

**Step by step**

### 1. **Framework and Libraries to Use**:

- Web Framework: Gin
- ORM: GORM
- Database: MySQL
- Authentication: JWT

### 2. **Project Structure:**

Use clean project structure with layer separation (e.g., controller, service, repository, model).

(current project structure)

 root-app
    ├── cmd
    │   └── server/
    ├── internal
    │   ├── api
    │   │   ├── dto/   # request and response file
    │   │   ├── handler/   # handler / controller files
    │   │   └── router/  # router files 
    │   ├── config/   # configuration setup files
    │   ├── contract/
    │   ├── database/
    │   ├── entities/
    │   ├── exception/
    │   ├── logger/
    │   ├── middleware/
    │   ├── repository/
    │   ├── service/
    │   └── utils/
    └── pkg
        ├── gin_helper/
        ├── jwt/
        ├── report_generator/
        ├── role_validator/
        └── web_response/

### 3. **API Endpoints**

a. **User Management:**
- POST /register → User registration.
- POST /login → User login.

b. **Event Management:**
- GET /events → View list of events.
- POST /events → Add new event.
- PUT /events/:id → Update specific event data.
- DELETE /events/:id → Delete specific event.

c. **Ticket Management:**
- GET /tickets → View list of tickets.
- POST /tickets → Buy or book tickets.
- GET /tickets/:id → View specific ticket details.
- PATCH /tickets/:id → Update ticket status to 'cancelled'.

d. **Reports or Analysis (Admin):**
- GET /reports/summary → Summary report of tickets sold and revenue.
- GET /reports/event/:id → Ticket report based on event.

### 4. **Validation Rules:**

a. Input data validation, for example:
- Event name must be unique.
- Event capacity cannot be negative.
- Ticket price must be valid (not negative).

b. Data relationship validation, such as:
- Sold tickets cannot be deleted.
- Events that have already taken place cannot be changed.

c. Attribute validation:
- Event capacity must be sufficient for ticket purchases.
- Ticket status (available, sold out, cancelled) must be consistent.

### 5. **API Documentation:**

- Create API documentation using tools like Swagger or Postman.

### 6. **Role-Based Access Control (RBAC):**

RBAC ensures that each user can only access data and features appropriate to their role. For ticketing systems, commonly used roles are:

- Admin: Has full access (CRUD all data).
- User: Can only buy tickets, view events, and cancel their own tickets.

## **Main Features**

### 1. **User Management**
- User registration and login.
- Authentication using JWT to ensure security.
- Role-based authorization to differentiate Admin and User access.

### 2. **Event Management:**
- CRUD for events with capacity and price validation.
- Event status (Active, Ongoing, Finished).

### 3. **Ticket Management:**
- Buy tickets (POST /tickets).
- View tickets by user (GET /tickets).
- Cancel tickets with cancellation rule validation.

### 4. **Reports and Analysis:**
- Report of tickets sold per event.
- Revenue summary based on events or time periods.

### 5. **System Security**
- Implementation of middleware for authentication and authorization.
- Input data validation to prevent errors and maintain system integrity.

### 6. **API Documentation**
- Complete documentation using Swagger or Postman that includes:
  - Endpoint explanations.
  - Payload and response examples.
  - Explanation of data relationships.

### 7. **Advanced Data Search and Filter**
- Provide data search capabilities using specific keywords in certain attributes (example: product name, description, or category).
- Add data filter features based on criteria such as:
  - Category: For example, product category, service, or class.
  - Status: Example: "Active", "Inactive", "Available".
  - Time: Filter based on specific dates (for example, transactions this week or this month).

### 8. **Pagination for Data Lists**
- Display data in small amounts per page to improve application performance and ease user navigation.
- Provide page and limit parameters to determine:
  - The page the user wants (example: page 2 of 10).
  - Amount of data per page (example: 10 or 20 data per page).
- API response includes pagination metadata such as:
  - current_page: Current page.
  - total_pages: Total number of pages.
  - total_items: Total amount of available data.

### 9. **Bonus Features (optional):**
- Unit Testing for Code
- File Upload/Download
- Audit Trail
- Caching
- GRPC

## **Development Guidelines**

- Do it gradually, focus on main features first.
- Make sure your code is clean and use relevant comments.
- Use database with dummy data for testing.