# User Management System

A robust, modular web application built with Go, adhering to **Clean Architecture** principles. This project provides a solid foundation for user management with role-based access control (User & Admin).

## 🚀 Features

*   **User Authentication**: Secure Signup, Login, and Logout functionality.
*   **Role-Based Access Control (RBAC)**: Distinct access levels for Users and Admins.
*   **Admin Dashboard**: comprehensive dashboard for admins to manage users (Create, Update, Delete, Search).
*   **Session Management**: Secure cookie-based sessions using `gin-contrib/sessions`.
*   **Clean Architecture**: Separation of concerns into `Delivery`, `Domain`, `UseCase`, and `Repository` layers.
*   **Logging**: Custom structured logging using `woodpecker`.

## 🛠️ Tech Stack

*   **Language**: Go (Golang)
*   **Web Framework**: [Gin Web Framework](https://github.com/gin-gonic/gin)
*   **Database**: PostgreSQL
*   **ORM**: [GORM](https://gorm.io/)
*   **Templating**: Go HTML Templates

## 📂 Project Structure

The project follows the standard Go project layout:

```
├── cmd
│   └── server          # Application entry point
├── internal
│   ├── delivery        # HTTP handlers, Router, Middleware
│   ├── domain          # Domain entities and interface definitions
│   ├── repository      # Database implementations (Postgres/GORM)
│   └── usecase         # Business logic and application services
├── pkg
│   └── woodpecker      # Custom logger functionality
├── static              # Static assets (CSS, JS, Images)
└── template            # HTML templates
```

## ⚡ Getting Started

### Prerequisites

*   Go 1.20+ installed
*   PostgreSQL running

### Installation

1.  **Clone the repository**
    ```bash
    git clone <repository-url>
    cd full-domain
    ```

2.  **Environment Setup**
    Create a `.env` file in the root directory (refer to `.env.example` if available) with the following variables:
    ```env
    DB_HOST=localhost
    DB_USER=your_user
    DB_PASSWORD=your_password
    DB_NAME=your_dbname
    DB_PORT=5432
    ```

3.  **Run the Application**
    ```bash
    go run cmd/server/main.go
    ```

4.  **Access the App**
    *   Home: `http://localhost:8080`
    *   Admin Login: `http://localhost:8080/admin/login`

## 📝 API Endpoints

*   `POST /api/signup` - Register a new user
*   `POST /api/login` - User login
*   `POST /api/logout` - User logout
*   `POST /api/users/password` - Update password
*   `POST /api/admin/login` - Admin login