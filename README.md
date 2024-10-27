<div align="center">

[![go.mod Go version](https://img.shields.io/badge/Go-v1.23.0-blue)](https://github.com/DavidMovas/Movies-Reviews)
[![codecov](https://codecov.io/gh/DavidMovas/Movies-Reviews/graph/badge.svg?token=RI6OY6VZC3)](https://codecov.io/gh/DavidMovas/Movies-Reviews)
[![Go Report Card](https://goreportcard.com/badge/github.com/DavidMovas/Movies-Reviews)](https://goreportcard.com/report/github.com/DavidMovas/Movies-Reviews)
[![Build Status](https://img.shields.io/badge/build-passing-brightgreen)](https://github.com/DavidMovas/Movies-Reviews)
[![Go Modules](https://img.shields.io/badge/go--modules-enabled-brightgreen)](https://blog.golang.org/using-go-modules)
[![License](https://img.shields.io/badge/license-Apache%20License%202.0-E91E63.svg?style=flat-square)](LICENSE)

</div>

# Movies-Reviews Service API

I have developed a service in Go using the [Echo](https://echo.labstack.com/) framework that serves as an API for handling requests. 
This is a movie service where users can authenticate and authorize using JWT tokens. 
OpenAPI documentation is available at the special endpoint [OpenAPI](#OpenAPI).

The service operates with a PostgresSQL database and is organized through Docker, allowing it to be deployed in isolated containers for both the server and the database. 
This setup ensures ease of managing dependencies and the ability to scale the application effortlessly.

Key features of the service include user management: retrieving, creating, editing, and deleting user profiles.
Additionally, it supports working with users, genres, movies, movie stars and movie reviews. 

Moreover, the application includes a substantial number of integration tests that ensure its stability and reliability. 
These tests help identify potential issues and provide confidence in the correct operation of the service after code changes.

------------------------------------------------------------------------------------------------
### Routes

##### Auth API:
| Method | Endpoint       | Description                                     | Auth |
|--------|----------------|-------------------------------------------------|------|
| POST   | /auth/register | Register a new user (Create user)               | -    |
| POST   | /auth/login    | Login a user. Returns access and refresh tokens | -    |

##### Users API:
| Method | Endpoint                        | Description                   | Auth  |
|--------|---------------------------------|-------------------------------|-------|
| GET    | /api/users/{userId}             | Get existing user by id       | any   |
| GET    | /api/users/{username}           | Get existing user by username | any   |
| PUT    | /api/users/{userId}             | Update existing user by id    | self  |
| PUT    | /api/users/{userId}/role/{role} | Update role by user id        | admin |
| DELETE | /api/users/{userId}             | Delete existing user (soft)   | admin |

##### Genres API:
| Method | Endpoint              | Description        | Auth   |
|--------|-----------------------|--------------------|--------|
| GET    | /api/genres           | Get all genres     | any    |
| GET    | /api/genres/{genreId} | Get genre by id    | any    |
| POST   | /api/genres           | Create new genre   | editor |
| PUT    | /api/genres/{genreId} | Update genre by id | editor |
| DELETE | /api/genres/{genreId} | Delete genre by id | editor |

##### Stars API:
| Method | Endpoint            | Description                                  | Auth   |
|--------|---------------------|----------------------------------------------|--------|
| GET    | /api/stars          | Get all stars (paginated, filtered, ordered) | any    |
| GET    | /api/stars/{starId} | Get star by id                               | any    |
| POST   | /api/stars          | Create new star                              | editor |
| PUT    | /api/stars/{starId} | Update star by id                            | editor |
| DELETE | /api/stars/{starId} | Delete star by id                            | editor |

##### Movies API:
| Method | Endpoint              | Description                                   | Auth   | 
|--------|-----------------------|-----------------------------------------------|--------|
| GET    | /api/movies           | Get all movies (paginated, filtered, ordered) | any    |
| GET    | /api/movies/{movieId} | Get movie by id                               | any    |
| POST   | /api/movies           | Create a new movie                            | editor |
| PUT    | /api/movies/{movieId} | Update movie by id                            | editor |
| DELETE | /api/movies/{movieId} | Delete movie by id                            | editor |

##### OpenAPI API:
| Method | Endpoint | Description  | Auth  |
|--------|----------|--------------|-------|
| GET    | /swagger | OpenAPI spec | admin |

------------------------------------------------------------------------------------------------
### Environment Variables

For the service to function correctly, you need to set environment variables in a `.env` file. Below is a list of available variables and their descriptions:
##### Database Configuration

- `DB_USER=user`  # Database username
- `DB_PASSWORD=password` # Database user password
- `DB_NAME=dbname` # Name of the database used by the service

##### Server Configuration

- `LOCAL=true` # Indicates whether the server runs in local mode
- `PORT=port` # Port on which the server will run (Default: 8000)
- `EXTERNAL_PORT=port` # Port used for external requests
- `DB_URL=postgres://${DB_USER}:${DB_PASSWORD}@db:5432/${DB_NAME}` # URL for connecting to the database

##### JWT Configuration

- `JWT_SECRET=secret` # Secret key for signing JWT
- `JWT_ACCESS_EXPIRATION=5m` # Access token lifetime
- `JWT_REFRESH_EXPIRATION=10m` # Refresh token lifetime

##### Logging Configuration

- `LOG_LEVEL=info` # Logging level (info, debug, warn, error)

##### Admin Configuration (initial creation)

- `ADMIN_USERNAME=name` # Admin username
- `ADMIN_EMAIL=email` # Admin email
- `ADMIN_PASSWORD=password` # Admin password

##### Pagination Configuration (optional)

- `PAGINATION_DEFAULT_SIZE=10` # Size of the default page (amount of items per page) (Default: 10) 
- `PAGINATION_MAX_SIZE=20` # Default page size limit (amount of items per page) (Default: 20)

------------------------------------------------------------------------------------------------
### OpenAPI

If you want to use OpenAPI, you can download the OpenAPI spec file from
[JSON](docs/swagger.json)
|
[YAML](docs/swagger.yaml)

Also, if you want to use Swagger UI please go next route: 
`http://example-host/swagger/index.html`

------------------------------------------------------------------------------------------------
### License

This project is licensed under the [Apache 2.0 License](LICENSE).
