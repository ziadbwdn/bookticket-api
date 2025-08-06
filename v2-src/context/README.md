# Book Ticketing API

**Project Status: UNDERGOING REFACTORING**

> **Warning:** The application is currently in a non-working state and is being actively refactored. The primary goal is to stabilize the communication between the React front-end and the Go backend and implement a predictable state management pattern.

## 1. Project Description

This is a ticketing application designed to [briefly describe what the app does, e.g., allow users to create, track, and resolve support tickets]. The front-end is built with React/TypeScript, and the backend service is written in Go.

## 2. Tech Stack

- **Frontend:** React, TypeScript, [e.g., Vite, Create React App], [e.g., Axios for data fetching], [e.g., Zustand/Context API for state management]
- **Backend:** Go (Golang), [e.g., Gin/Echo framework], [e.g., GORM for database interaction]
- **Database:** [e.g., PostgreSQL, MySQL, SQLite]
- **Development:** Concurrently (for running both servers at once)

## 3. Local Development Setup

### Prerequisites

- Node.js (v18 or higher)
- Go (v1.19 or higher)
- Docker (for database) OR a locally installed instance of [Your DB]

### Installation & Setup

1.  **Clone the repository:**
    ```bash
    git clone [your-repo-url]
    cd [your-repo-folder]
    ```

2.  **Backend Setup (Go):**
    ```bash
    # Navigate to the backend directory
    cd server/

    # Create a .env file from the example
    cp .env.example .env

    # Fill in the .env file with your DB credentials, server port, etc.
    # e.g., DB_HOST=localhost, DB_PORT=5432, GO_SERVER_PORT=8080

    # Install dependencies
    go mod tidy

    # Run database migrations (if any)
    go run cmd/migrate/main.go
    ```

3.  **Frontend Setup (React/TypeScript):**
    ```bash
    # Navigate to the frontend directory
    cd client/

    # Create a .env file from the example
    cp .env.example .env

    # Fill in the .env file with the Go backend API URL
    # e.g., VITE_API_BASE_URL=http://localhost:8080/api

    # Install dependencies
    npm install
    ```

### Running the Application

This project uses `concurrently` to run both the front-end and backend servers with a single command from the root directory.

```bash
# From the project root
npm run dev
