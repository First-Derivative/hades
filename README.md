# Go HTTP Server with Gin and JWT Authentication

This Go HTTP server is built using the Gin framework and employs JWT (JSON Web Tokens) for authentication. It allows you to secure your endpoints and control access to your resources. To run this server, you need to set up specific environment properties defined in a `.env` file. Also uses SQL (I recommend [planetscale](https://planetscale.com/) for hobby devs)

## Prerequisites

Before running the server, make sure you have the following prerequisites installed and configured:

- **Go**: Make sure you have Go installed on your system. You can download it from [here](https://golang.org/dl/).

## Getting Started

1. **Clone Repository**:

```bash
git clone https://github.com/First-Derivative/hades.git
cd hades
```

2. **Create .env File**:

Create a .env file in the root directory of your project with the following properties:

```bash
PORT=8080 # Port on which the server will run
DNS="some sql dns connection string"      # Domain or hostname for sql server
JWT_SECRET=your-secret-key # Secret key for JWT token generation and validation
```

Build and Run:

`go run main.go` or 
`go build main.go && ./main`

3. **Access the Server**:
   The server will start, and you can access it at <http://127.0.0.1:8080> (or the DNS and port you specified in the .env file).
   Example:
   GET /validate: User authentication. Send a POST request with credentials to obtain a JWT token.

## API Endpoints

### /signup (POST)

<details>

- **Description**: Allows users to sign up by providing valid credentials.
- **URL**: `/signup`
- **Method**: `POST`
- **Request Body**:

  - Format: JSON
  - Fields:
    - `user` (SignupRequest) - User registration details

- **Success Response (201 Created)**:

  - Response Body: JSON
    - `token` (TokenResponse) - User successfully registered

- **Error Response (400 Bad Request)**:
  - Response Body: JSON - `error` (ErrorResponse) - Invalid input data
  </details>

### /login (POST)

<details>

- **Description**: Allows users to authenticate by providing valid credentials.
- **URL**: `/login`
- **Method**: `POST`
- **Request Body**:

  - Format: JSON
  - Fields:
    - `email`: email as string
    - `password`: password as string

- **Success Response (200 OK)**:

  - Response Body: JSON
    - `status` : `success`

- **Error Response (401 Unauthorized)**:
  - Response Body: JSON - `error` (ErrorResponse) - Invalid credentials
  </details>

### /logout (GET)

<details>

- **Description**: Logs out the currently authenticated user.
- **URL**: `/logout`
- **Method**: `GET`
- **Authentication**: Requires a valid JWT token (ApiKeyAuth)

- **Success Response (204 No Content)**:

  - User successfully logged out

- **Error Response (401 Unauthorized)**:
  - Response Body: JSON - `error` (ErrorResponse) - Authentication required
  </details>

### /validate (GET)

  <details>

- **Description**: Validates the provided JWT token.
- **URL**: `/validate`
- **Method**: `GET`
- **Authentication**: Requires a valid JWT token (ApiKeyAuth)

- **Success Response (200 OK)**:

  - Response Body: JSON
    - `status` : `authenticated`
    - `success`: `true`,
    - `user`: `user email`

- **Error Response (401 Unauthorized)**:
  - Response Body: JSON
    - `error` (ErrorResponse) - Invalid token or authentication required
    </details>

### Request and Response Data Structures

```json

// POST body for signup
{
  "email": "string",
  "password": "string",
  "firstName": "string" | null,
  "lastName": "string" | null
}

// POST body for login
{
  "email": "string",
  "password": "string"
}

// Error response
{
    "error": "string"
}

```
