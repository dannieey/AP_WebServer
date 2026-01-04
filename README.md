# AP_WebServer

## Project Overview

**AP_WebServer** is a Go HTTP server implementing in-memory CRUD operations for key-value pairs. It supports concurrent requests, a background worker for statistics, graceful shutdown, and usage of generics. The project fulfills all assignment requirements (Part A–E).

## Project Structure

* **server/** – HTTP API handlers and routes
* **store/** – Generic in-memory store
* **worker/** – Background worker printing server statistics
* **model/** – Structs for JSON responses

## Installation and Run

1. Clone the repository
2. Navigate to the project folder
3. Run `go run .` to start the server on port `:8080`
4. Background worker prints statistics every 5 seconds

## API Overview (Part A)

* **POST /data** – Add a key-value pair; validates input
* **GET /data** – Returns all key-value pairs
* **GET /data/{key}** – Returns a single key-value pair
* **DELETE /data/{key}** – Deletes a key-value pair
* **GET /stats** – Returns server statistics: number of requests, number of keys, uptime

## Concurrency and Thread Safety (Part B)

* In-memory store protected with mutexes
* Request counter is atomic
* Server handles multiple concurrent HTTP requests safely

## Background Worker (Part C)

* Runs in a separate goroutine
* Prints statistics every 5 seconds
* Controlled with a stop channel and select statement
* Stops gracefully on shutdown

## Graceful Shutdown (Part D)

* Catches OS signals (`SIGINT`, `SIGTERM`)
* Stops background worker
* Shuts down HTTP server gracefully, allowing in-flight requests to complete

## Generics Requirement (Part E)

* Generic store implemented
* Supports Set, Get, Delete, Snapshot methods
* Instantiated as `Store[string, string]`

## Summary

* CRUD endpoints are fully functional
* Concurrency is safe and thread-safe
* Background worker works as expected
* Graceful shutdown implemented
* Generics are used correctly


