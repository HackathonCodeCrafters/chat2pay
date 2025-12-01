Simple REST API JWT Auth
============================

This app is very straight forward, build with golang and use  [Fiber](https://gofiber.io/) as the framework

Make sure you have mysql for running this app

### Prerequisite :

- Docker
- OS Linux Docker Image,
- Golang version 1.19 or latest

# Install
#### MySQL Query


Since im not providing endpoint for register you'll need to insert the user 
in your local mysql, simply run this command

    INSERT INTO dsi_technical.users (id, created_at, updated_at, deleted_at, email, full_name, age, mobile_number, password) VALUES('bf506362-4e81-43d0-8238-d377bc8dab80', now(), now(), NULL, 'iniemail@gmail.com', 'ini full name', 99, '081234567890', '$2a$10$wZIYJUPjmOYvI8aTie7Qd.Sw11X169/0yo0k17NnCIrEXFgDl38Pi');

we'll use this account for login

    {
        email : "iniemail@gmail.com"
        password: "inipassword"
    }

## Install localy
    go mod tidy

#### Config

    Make sure to change `app.yaml` variable value with your local db

#### Run the app

    go run main.go

#### Run the tests

    cd internal/service
    go test -verbose --cover

## Build with docker

_**note:**_ The docker-compose file contain mysql

    docker-compose -f deployment/docker-compose.yaml up -d

make sure to double check `app.yaml` file with the db


# REST API

There are 2 endpoint in this app

## Get Profile

### Request

you need to generate the bearer token from endpoint **`Login`**

`GET /api/profile`

    curl -i -H 'Accept: application/json' -H 'Authorization: Bearer {{Token}}'  http://localhost:{{port}}/api/profile

### Response

    HTTP/1.1 200 OK
    Content-Length: 179
    Content-Type: application/json
    Date: Sun, 05 Nov 2023 14:36:07 GMT

    {"data":{"id":"d13a939c-7518-431c-bb04-98701a3ab750","fullname":"ini fullname","email":"iniemail@gmail.com","age":99,"mobile_number":"081234567890"},"error":null,"status":true}%

## Login

### Request

`POST /api/login`

    curl -i -H 'Accept: application/json' -d 'email=iniemail@gmail.com&password=inipassword' http://localhost:{{port}}/api/login

### Response

    HTTP/1.1 200 OK
    Content-Length: 328
    Content-Type: application/json
    Date: Sun, 05 Nov 2023 14:33:07 GMT

    {"data":{"id":"d13a939c-7518-431c-bb04-98701a3ab750","email":"iniemail@gmail.com","access_token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6ImQxM2E5MzljLTc1MTgtNDMxYy1iYjA0LTk4NzAxYTNhYjc1MCIsImVtYWlsIjoiYXNlbG9sZUBnbWFpbC5jb20iLCJleHAiOjE2OTkxOTY1ODd9.0oGK-JueSveRYdr9iUz4jVxGL3ctS05k3q9cI7FT0-U"},"error":null,"status":true}%

Folder Structure
============================

### Top-level directory layout

    .
    ├── bootstrap                 # Connector with other depedency for example mysql, redis, kafka etc
    ├── config                    # Configuration folder that used by apps (.env like)
    ├── deployment                # Dockerfile and docker-compose.yaml
    ├── internal                  # Main directory for rest api
    ├── migrations                # Migrations for database
    ├── go.mod                  
    ├── main.go
    └── README.md

### Internal folder

The actual source files of this project are stored inside the
`internal` directory. this contains `api`,`entities`,`helper`, `middleware`, `repository`, `service` and `pkg`

    .
    ├── internal
    │   ├── api                 # Contains all api this
    │   │   ├── dto             # Used for processing the struct before send it as response
    │   │   ├── handlers        # Files that handling the requests and the response of endpoint
    │   │   ├── presenter       # Struct that used as response of endpoint
    │   │   ├── request_model   # Struct for request parameter/body
    │   │   └── router.go       # Routes of all endpoints
    │   ├── const               # Constant variable list
    │   ├── entities            # Struct of data that represent from database
    │   ├── middlewares         # List of middlewares runs before handler
    │   ├── pkg                 # List of package we can use  
    │   ├── repositories        # Database interactor
    │   └── service             # Business Logic of the endpoints
    └── ...

# Database Schema

| Column name   | Data Type    | Not Null | Key |
|---------------|--------------|----------|-----|
| id            | Varchar(35)  | ✅        | PK  |
| email         | Varchar(255) | ✅        |     |
| full_name     | Varchar(255) | ✅        |     |
| password      | Varchar(255) | ✅        |     |
| age           | Int(2)       | ✅        |     |
| mobile_number | Varchar(13)  | ✅        |     |
| created_at    | Timestamp    | ✅        |     |
| updated_at    | Timestamp    | ✅        |     |
| deleted_at    | Timestamp    |          |     |