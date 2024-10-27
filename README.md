# Movies-Reviews

An API for managing movies reviews.

### Usage

You can use this API to manage movies reviews. You can create, read, update and delete movies and reviews.

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