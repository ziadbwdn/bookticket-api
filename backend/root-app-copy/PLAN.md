# **Back-End System: Study Case Ticketing** - Planning Set of task

## **Progress**

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

**We will review, correct few lines of code, and refactor if necessary**. try our best at practicing idiomatic go code structure and clean code, with minimalize or without changing (more preferred) the existing code.

Try Look at README.md for what this project purposed to check each functionality of each code previously constructed and the reason why

**Step by step**

### 1. **Framework and Libraries to Use**:

- Web Framework: Gin
- ORM: GORM
- Database: MySQL
- Authentication: JWT

### 2. **Project Structure:**

(current progress)

root-app
├── PLAN.md
├── README.md
├── cmd
│   └── server
│       └── main.go
├── go.mod
├── go.sum
├── internal
│   ├── api
│   │   ├── dto
│   │   │   ├── activity_dto.go
│   │   │   ├── error_response.go
│   │   │   ├── event_dto.go
│   │   │   ├── profile_dto.go
│   │   │   ├── ticket_dto.go
│   │   │   └── user_dto.go
│   │   ├── handler
│   │   │   ├── activity_handler.go
│   │   │   ├── report_handler.go
│   │   │   ├── ticket_handler.go
│   │   │   └── user_handler.go
│   │   └── router
│   │       ├── activity_route.go
│   │       ├── event_route.go
│   │       ├── router.go
│   │       ├── ticket_route.go
│   │       └── user_route.go
│   ├── config
│   │   └── config.go
│   ├── contract
│   │   ├── activity.go
│   │   ├── event.go
│   │   ├── report.go
│   │   ├── ticket.go
│   │   └── user.go
│   ├── database
│   │   ├── connection.go
│   │   └── migration.go
│   ├── entities
│   │   ├── activities.go
│   │   ├── event.go
│   │   ├── ticket.go
│   │   └── user.go
│   ├── exception
│   │   └── exception.go
│   ├── logger
│   │   └── logger.go
│   ├── middleware
│   │   └── auth.go
│   ├── repository
│   │   ├── activity_repository.go
│   │   ├── event_repository.go
│   │   ├── report_repository.go
│   │   ├── ticket_repository.go
│   │   └── user_repository.go
│   ├── service
│   │   ├── activity_service.go
│   │   ├── event_service.go
│   │   ├── report_service.go
│   │   ├── ticket_service.go
│   │   └── user_auth_service.go
│   └── utils
│       ├── decimal.go
│       ├── pass_validator.go
│       └── uuid.go
└── pkg
    ├── gin_helper
    │   └── gin_helper.go
    ├── jwt
    │   └── jwt_helper.go
    ├── report_generator
    ├── role_validator
    │   └── role_validator.go
    └── web_response
        └── web_response.go

### 3. Plan and Things to Understand

after look at README.md, this is for the instruction:

* Any code from entities, exception, utils, config, database, middleware, logger, pkg, is untouchable
* certain contract, repository, and services which related to user and user_activity is untouchable
* contract > repository > services packages must be aligned each other and well
* after our task is done, we will be compose handlers and router package left. handler of ticket and report were empty, and event_handler were not made yet. router of event, ticket were empty, and report router were not made yet too