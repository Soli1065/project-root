project-root/
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



# API Documentation

## Overview
This document provides documentation for the APIs available in the Dastanamon application.

### Base URL
http://172.16.170.47:8080


## Authentication

### Login
POST /auth/login


- **Description:** Endpoint to authenticate user login.
- **Request:**
  - Body:
    ```json
    {
      "username": "example",
      "password": "password"
    }
    ```
- **Response:**
  - Success:
    ```json
    {
      "token": "JWT_TOKEN"
    }
    ```
  - Error:
    ```json
    {
      "error": "Invalid credentials"
    }
    ```

### Register
POST /auth/register


- **Description:** Endpoint to register a new user.
- **Request:**
  - Body:
    ```json
    {
      "username": "example",
      "email": "example@example.com",
      "password": "password"
    }
    ```
- **Response:**
  - Success:
    ```json
    {
      "message": "User registered successfully"
    }
    ```
  - Error:
    ```json
    {
      "error": "Username already exists"
    }
    ```

## User Management

### Get All Users
GET /users


- **Description:** Endpoint to retrieve all users.

### Get User by ID
GET /users/{id}

- **Description:** Endpoint to retrieve a user by ID.

### Create User
POST /users


- **Description:** Endpoint to create a new user.

### Update User
PUT /users/{id}


- **Description:** Endpoint to update a user by ID.

### Delete User
DELETE /users/{id}


- **Description:** Endpoint to delete a user by ID.

## Content Management

### Get All Content
GET /content


- **Description:** Endpoint to retrieve all content.

### Get Content by ID
GET /content/{id}


- **Description:** Endpoint to retrieve content by ID.

### Create Content
POST /content


- **Description:** Endpoint to create new content.

### Update Content
PUT /content/{id}


- **Description:** Endpoint to update content by ID.

### Delete Content
DELETE /content/{id}


- **Description:** Endpoint to delete content by ID.

## Category Management

### Get All Categories
GET /categories


- **Description:** Endpoint to retrieve all categories.

### Get Category by ID
GET /categories/{id}


- **Description:** Endpoint to retrieve a category by ID.

### Create Category
POST /categories


- **Description:** Endpoint to create a new category.

### Update Category
PUT /categories/{id}

- **Description:** Endpoint to update a category by ID.

### Delete Category
DELETE /categories/{id}

- **Description:** Endpoint to delete a category by ID.
