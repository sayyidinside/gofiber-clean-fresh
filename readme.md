# GoFiber Clean Fresh

GoFiber Clean Fresh is a base template for Go projects, structured using Clean Architecture principles with some modifications to suit organizational requirements and knowledge from previous projects (Node.js Express-based). This template integrates essential packages and comes with user management, authentication-related middleware, and role/permission management.

## Features / Technologies Used

- **GoFiber**: Web framework for building fast and scalable APIs.
- **GORM**: Object-Relational Mapper (ORM) for MySQL database operations, utilizing GORM datatypes.
- **Air**: Live reload for Go applications during development.
- **Zap**: Fast, structured logging.
- **Validator V10**: Validation of incoming data for requests.
- **Gocron**: Job scheduling for recurring tasks.
- **Bluemonday**: HTML sanitizer for handling user-generated content securely.
- **Viper**: For configuration management, with support for environment variables and multiple file formats (YAML, JSON, etc.), including auto-reloading of configuration files.
- **Clean Architecture**: A layered approach to structure the codebase for maintainability and scalability.

## Project Structure

```bash
├── cmd/                     # Main application entry points
│   ├── server/              # HTTP server setup
│   ├── worker/              # Background worker setup
├── infrastructure/          # Infrastructure-specific code (frameworks, DB, etc.)
│   ├── config/              # Configuration files
│   ├── database/            # Database setup and implementations (GORM)
│   ├── logger/              # Logging setup (zap)
│   ├── scheduler/           # Scheduling logic (gocron)
├── interfaces/              # Interface adapters (Delivery layer)
│   ├── http/                # HTTP delivery (GoFiber routes)
│   │   ├── routes/          # Route definitions
│   │   ├── middleware/      # HTTP middleware (auth, JWT, role-based)
│   │   ├── auth/            # HTTP handlers for auth-related routes
│   │   ├── user/            # HTTP handlers for user-related routes
│   │   ├── role/            # HTTP handlers for role-related routes
│   │   └── permission/      # HTTP handlers for permission-related routes
│   ├── model/               # Data transfer objects (DTOs) for mapping HTTP <-> domain
├── domain/                  # Core business logic
│   ├── auth/                # Auth domain (user, role, permission)
│   │   ├── entity/          # Domain entities/model for auth
│   │   ├── service/         # Business logic and use cases for auth
│   │   └── repository/      # Repository interfaces for auth
│   ├── user/                # User domain logic
│   │   ├── entity/          # User domain entities/model
│   │   ├── service/         # Business logic and use cases for user
│   │   └── repository/      # Repository interfaces for user
│   ├── role/                # Role domain logic
│   │   ├── entity/          # Role domain entities/model
│   │   ├── service/         # Business logic and use cases for role
│   │   └── repository/      # Repository interfaces for role
│   └── permission/          # Permission domain logic
│       ├── entity/          # Permission domain entities/model
│       ├── service/         # Business logic and use cases for permission
│       └── repository/      # Repository interfaces for permission
├── pkg/                     # Shared libraries and utilities
│   └── utils/               # Generic utility functions (not domain-specific)
├── tests/                   # Unit and integration tests
└── .env                     # Environment variables
```

## Installation

1. **Clone the repository:**

```bash
git clone https://github.com/your-org/gofiber-clean-fresh.git
cd gofiber-clean-fresh
```

2. **Set up environment variables:**

Create a `.env` file based on `.env.example` and update the configuration as needed.

3. **Install dependencies:**

```bash
go mod tidy
```

4. **Run the application (with live reload):**

```bash
air
```

## User Management

- **Users**: Manage user accounts.
- **Roles**: Assign different roles to users.
- **Permissions**: Define and assign permissions to roles.

## Auth Middleware

The project includes JWT-based authentication, as well as role-based access control middleware. You can extend the authentication middleware as needed.

## Contributing

Feel free to submit issues or pull requests to improve this project. Make sure to follow the contribution guidelines.

## License

This project is licensed under the MIT License.
