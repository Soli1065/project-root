dastanamon/
├── cmd/
│   └── main.go
├── internal/
│   ├── api_gateway/
│   │   ├── gateway.go
│   │   └── routes.go
│   ├── auth/
│   │   ├── auth.go
│   │   ├── handler.go
│   │   └── model.go
│   ├── user/
│   │   ├── user.go
│   │   ├── handler.go
│   │   └── model.go
│   ├── content/
│   │   ├── content.go
│   │   ├── handler.go
│   │   └── model.go
│   ├── category/
│   │   ├── category.go
│   │   ├── handler.go
│   │   └── model.go
│   ├── recommendation/
│   │   ├── recommendation.go
│   │   ├── handler.go
│   │   └── model.go
│   └── config/
│       └── config.go
├── pkg/
│   ├── database/
│   │   └── database.go
│   └── middleware/
│       └── jwt.go
└── go.mod


==================================================



API Documentation
Overview
This document provides documentation for the APIs available in the Dastanamon application.

Base URL
arduino
Copy code
https://api.example.com
Authentication
Login
bash
Copy code
POST /auth/login
Description: Endpoint to authenticate user login.
Request:
Body:
json
Copy code
{
  "username": "example",
  "password": "password"
}
Response:
Success:
json
Copy code
{
  "token": "JWT_TOKEN"
}
Error:
json
Copy code
{
  "error": "Invalid credentials"
}
Register
arduino
Copy code
POST /auth/register
Description: Endpoint to register a new user.
Request:
Body:
json
Copy code
{
  "username": "example",
  "email": "example@example.com",
  "password": "password"
}
Response:
Success:
json
Copy code
{
  "message": "User registered successfully"
}
Error:
json
Copy code
{
  "error": "Username already exists"
}
User Management
Get All Users
bash
Copy code
GET /users
Description: Endpoint to retrieve all users.
Get User by ID
bash
Copy code
GET /users/{id}
Description: Endpoint to retrieve a user by ID.
Create User
bash
Copy code
POST /users
Description: Endpoint to create a new user.
Update User
bash
Copy code
PUT /users/{id}
Description: Endpoint to update a user by ID.
Delete User
bash
Copy code
DELETE /users/{id}
Description: Endpoint to delete a user by ID.
Content Management
Get All Content
bash
Copy code
GET /content
Description: Endpoint to retrieve all content.
Get Content by ID
bash
Copy code
GET /content/{id}
Description: Endpoint to retrieve content by ID.
Create Content
bash
Copy code
POST /content
Description: Endpoint to create new content.
Update Content
bash
Copy code
PUT /content/{id}
Description: Endpoint to update content by ID.
Delete Content
bash
Copy code
DELETE /content/{id}
Description: Endpoint to delete content by ID.
Category Management
Get All Categories
bash
Copy code
GET /categories
Description: Endpoint to retrieve all categories.
Get Category by ID
bash
Copy code
GET /categories/{id}
Description: Endpoint to retrieve a category by ID.
Create Category
bash
Copy code
POST /categories
Description: Endpoint to create a new category.
Update Category
bash
Copy code
PUT /categories/{id}
Description: Endpoint to update a category by ID.
Delete Category
bash
Copy code
DELETE /categories/{id}
Description: Endpoint to delete a category by ID.