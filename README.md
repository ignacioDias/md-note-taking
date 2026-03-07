# Markdown Note-Taking App

A full-stack web application for creating, editing, and managing markdown notes with user authentication and session management.

Project based on: https://roadmap.sh/projects/markdown-note-taking-app

## Features

- User authentication and authorization with session-based auth
- Create, read, update, and delete markdown notes
- User profile management
- Customizable theme settings (background color)
- Responsive design for desktop and mobile
- Secure session management with Redis
- PostgreSQL database for persistent storage

## Tech Stack

### Backend
- Go (Golang)
- PostgreSQL - Main database
- Redis - Session cache and management
- Docker & Docker Compose - Containerization

### Frontend
- HTML5
- CSS3 (Vanilla CSS with CSS Variables)
- JavaScript (Vanilla JS)

## Prerequisites

- Go 1.21 or higher
- Docker and Docker Compose
- Git

## Installation

### 1. Clone the repository

```bash
git clone <repository-url>
cd md-editor
```

### 2. Start the infrastructure

Start PostgreSQL and Redis using Docker Compose:

```bash
docker-compose up -d
```

This will start:
- PostgreSQL on port 5432
- Redis on port 6379

### 3. Build the application

```bash
go build ./cmd/mdeditor
```

### 4. Run the application

```bash
./mdeditor
```

The application will be available at `http://localhost:8080` (or the configured port).

## Usage

### Getting Started

1. Navigate to the home page
2. Click "Sign up" to create a new account
3. Log in with your credentials
4. Start creating and managing your markdown notes

### Available Routes

- `/` - Home page
- `/register` - User registration
- `/login` - User login
- `/dashboard` - Main dashboard (requires authentication)
- `/settings` - User settings and theme customization (requires authentication)
- `/me` - User profile (requires authentication)
- `/notes/{id}` - Note editor + viewer (requires authentication)

### API Endpoints

#### Authentication
- `POST /api/auth/register` - Register new user
- `POST /api/auth/login` - User login
- `DELETE /api/auth/logout` - User logout

#### User
- `GET /api/me` - Get current user info
- `PUT /api/me` - Update user info
- `DELETE /api/me` - Delete user account

#### Notes
- `POST /api/notes` - Create new note
- `GET /api/notes/{id}` - Get specific note
- `PUT /api/notes/{id}` - Update note
- `DELETE /api/notes/{id}` - Delete note
- `GET /api/me/notes` - Get all notes for current user

## Development

### Database Migrations

Database schema is managed through the application's repository layer. Ensure PostgreSQL is running before starting the application.

### Session Management

Sessions are stored in Redis with automatic expiration. Each session is tied to a unique session ID stored in HTTP-only cookies.

### Theme Customization

Users can customize their background color preference through the Settings page. The preference is stored in browser localStorage and persists across sessions.

## Environment Variables

Configure the following environment variables as needed:

- `DATABASE_URL` - PostgreSQL connection string
- `REDIS_URL` - Redis connection string
- `PORT` - Application port (default: 8080)
- `SESSION_TIMEOUT` - Session timeout duration
- `ENV` - notProduction if it's in development

## Security Features

- Session-based authentication
- HTTP-only cookies for session management
- Password hashing
- Protected API endpoints with authentication middleware
- CSRF protection considerations
