# Simple Bank

Simple Bank is a simple banking application built with Go. This project provides RESTful and gRPC APIs for managing bank accounts, transactions, user authentication, and more.

## Table of Contents

- [Introduction](#introduction)
- [Features](#features)
- [System Requirements](#system-requirements)
- [Installation](#installation)
- [Running the Project](#running-the-project)
- [API Documentation](#api-documentation)
- [Technologies Used](#technologies-used)

---

## Introduction

Simple Bank is a mock banking application where users can create accounts, perform transactions, and manage their accounts. The project is designed for learning and practicing modern software development techniques such as:

- Building RESTful and gRPC APIs.
- Using PostgreSQL as the database.
- Implementing JWT authentication.
- Integrating Redis for asynchronous task processing.

---

## Features

- **User Management**: Register, login, and update user information.
- **Account Management**: Create accounts and view balances.
- **Transactions**: Transfer money between accounts.
- **Authentication**: Protect APIs using JWT.
- **Swagger Integration**: Auto-generated API documentation.
- **gRPC Support**: Efficient communication between services.

---

## System Requirements

- **Go**: Version 1.20 or higher.
- **PostgreSQL**: Version 13 or higher.
- **Redis**: Version 6 or higher.
- **Protobuf**: To compile `.proto` files.
- **Docker**: (Optional) To run the application in a container.

---

## Installation

1. **Clone the project**:
   ```bash
   git clone https://github.com/phongnd2802/simple-bank.git
   cd simple-bank
   ```

2. **Install dependencies**:
   ```bash
   go mod tidy
   ```

3. **Initialize the database**:
     ```bash
     make network
     make postgres
     make redis
     make migrate-up
     ```

---

## Running the Project

1. **Run the application**:
   ```bash
   make dev
   ```

2. **Run with Docker**:
   - Build and run the container:
     ```bash
     docker-compose up --build
     ```

3. **Access the APIs**:
   - REST API: `http://localhost:8080`
   - Swagger UI: `http://localhost:8080/swagger/`
   - gRPC: `localhost:9090`

---

## API Documentation

- **Swagger**: Auto-generated RESTful API documentation is available at `/swagger/`.
- **gRPC**: gRPC definitions are located in the `proto/` directory.

---

## Technologies Used

- **Language**: Go
- **Database**: PostgreSQL
- **Cache**: Redis
- **Authentication**: JWT
- **gRPC**: Service-to-service communication
- **Swagger**: API documentation
- **Asynq**: Asynchronous task processing
