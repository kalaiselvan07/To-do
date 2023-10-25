# TODO Microservices Setup

Welcome to the TODO microservices setup! This setup allows you to quickly deploy and run the TODO services using Docker and Docker Compose.

## Quick Start

**Prerequisite**: Ensure you have `Docker` and `Docker Compose` installed.

1. **Clone or download this repository**.
   
2. **Navigate to the directory** containing the `docker-compose.yml` file.

3. **Run the following command**:

    ```bash
    docker-compose up
    ```

That's it! The images are already available on Docker Hub, so there's no need to build anything locally.

## Services

* **authn** - Authentication service.
* **authz** - Authorization service.
* **app** - Main application service.

Thank you for using the TODO microservices setup! 

# TODO API

**API for managing TODOs**  
_Version: 1.0.0_  
Server: `http://localhost:8087`

---

## Endpoints

### 1. User Login
**Endpoint:** `/login`  
**Method:** `POST`

**Request Body:**
- `Username`: _string_
- `Password`: _string_

**Response (200):**
- `token`: _string_

---

### 2. Retrieve All TODOs
**Endpoint:** `/todo`  
**Method:** `GET`

**Headers:**
- `apikey`: _string (required)_

**Response (200):**
- List of TODOs:
  - `id`: _integer_
  - `item`: _string_
  - `completed`: _boolean_

---

### 3. Add a New TODO
**Endpoint:** `/todo`  
**Method:** `POST`

**Headers:**
- `apikey`: _string (required)_

**Request Body:**
- `id`: _integer_
- `item`: _string_
- `completed`: _boolean_

**Responses:**  
- `200`: Successful operation
- `400`: Bad request

---

### 4. Retrieve a Specific TODO by ID
**Endpoint:** `/todo/{id}`  
**Method:** `GET`

**Parameters:**
- `apikey`: _string (required in header)_
- `id`: _integer (required in path)_

**Response (200):**
- A single TODO:
  - `id`: _integer_
  - `item`: _string_
  - `completed`: _boolean_

---

### 5. Update a TODO by ID
**Endpoint:** `/todo/{id}`  
**Method:** `PATCH`

**Parameters:**
- `apikey`: _string (required in header)_
- `id`: _integer (required in path)_

**Request Body:**
- `id`: _integer_
- `item`: _string_
- `completed`: _boolean_

**Responses:**  
- `200`: Successful operation
- `400`: Bad request
- `404`: Not found

---

## Schema:

**Todo:**
- `id`: _integer_
- `item`: _string_
- `completed`: _boolean_
