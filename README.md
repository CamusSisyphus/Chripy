
# **Chirpy API Documentation**  

Chirpy is a simple social media web server built with Go, where users can create and manage “chirps” (posts). Below is the API documentation for the project.  
## Key Features of Chirpy API
### User Management

Create, retrieve, update, and list users.
Users have is_chirpy_red status to indicate premium membership.
Authentication is based on JWT tokens.
### Chirp Management

Users can post, retrieve, and delete "chirps" (posts).
Chirps can be queried by author_id and sorted by ID order.
### Authentication & Token Handling

Login provides both access and refresh tokens.
Refresh tokens can be revoked or used to renew access.


---

## **User Endpoints**  

### **GET /api/users**  
Returns a list of all users, including their IDs, email addresses, and Chirpy Red membership status.  

**Response Example:**  
```json
[
    {
        "id": 1,
        "email": "alice@example.com",
        "is_chirpy_red": false
    }
]
```

---

### **POST /api/users**  
Creates a new user.  

**Request Body Example:**  
```json
{
    "email": "alice@example.com",
    "password": "pass123"
}
```

**Response Example:**  
```json
{
    "id": 1,
    "email": "alice@example.com",
    "is_chirpy_red": false
}
```

---

### **PUT /api/users**  
Updates a user’s details. Requires authentication with a valid **access token**.  

**Headers:**  
```
Authorization: Bearer <access-token>
```

**Request Body Example:**  
```json
{
    "email": "alice@example.com",
    "password": "newPassword!"
}
```

**Response Example:**  
```json
{
    "id": 1,
    "email": "alice@example.com",
    "is_chirpy_red": false
}
```

---

### **GET /api/users/{ID}**  
Retrieves a user by their ID.  

**Example:**  
`GET /api/users/1`  

**Response Example:**  
```json
{
    "id": 1,
    "email": "alice@example.com",
    "is_chirpy_red": false
}
```

**Response:**  
- `404`: If user ID is not found.

---

## **Authentication Endpoints**  

### **POST /api/login**  
Logs in a user with their email and password.  

**Request Body Example:**  
```json
{
    "email": "alice@example.com",
    "password": "newPassword!"
}
```

**Response Example:**  
```json
{
    "id": 1,
    "email": "alice@example.com",
    "is_chirpy_red": false,
    "token": "<access-token>",
    "refresh_token": "<refresh-token>"
}
```

- **401:** If credentials are incorrect.

---

### **POST /api/refresh**  
Exchanges a refresh token for a new access token.  

**Headers:**  
```
Authorization: Bearer <refresh-token>
```

**Response Example:**  
```json
{
    "token": "<new-access-token>"
}
```

- **401:** If refresh token is invalid.

---

### **POST /api/revoke**  
Revokes the provided refresh token.  

**Headers:**  
```
Authorization: Bearer <refresh-token>
```

- **200:** On success.  
- **401:** If the token is invalid.

---

## **Chirp Endpoints**  

### **GET /api/chirps**  
Retrieves all chirps, sorted by ID in ascending order.  

**Response Example:**  
```json
[
    {
        "id": 1,
        "body": "example text",
        "author_id": 1
    }
]
```

**Optional Query Parameters:**  
- `?author_id=<ID>`: Filter by the author’s ID.  
- `?sort=desc`: Sort chirps in descending order.

---

### **POST /api/chirps**  
Creates a new chirp. Requires authentication with a valid **access token**.  

**Headers:**  
```
Authorization: Bearer <access-token>
```

**Request Body Example:**  
```json
{
    "body": "example text"
}
```

**Response Example:**  
```json
{
    "id": 1,
    "body": "example text",
    "author_id": 1
}
```

- **401:** If the access token is missing or invalid.

---

### **GET /api/chirps/{ID}**  
Retrieves a specific chirp by its ID.  

**Example:**  
`GET /api/chirps/1`  

**Response Example:**  
```json
{
    "id": 1,
    "body": "example text",
    "author_id": 1
}
```

- **404:** If the chirp ID is not found.

---

### **DELETE /api/chirps/{ID}**  
Deletes a chirp by its ID. The authenticated user must be the chirp’s author.  

**Headers:**  
```
Authorization: Bearer <access-token>
```

- **200:** On success.  
- **401:** If the token is invalid or the user is not authorized to delete the chirp.

---

## **Authentication Flow Example**  

1. **Login:**  
   ```
   POST /api/login
   ```
   **Request Body:**  
   ```json
   {
       "email": "alice@example.com",
       "password": "newPassword!"
   }
   ```
   **Response:**  
   ```json
   {
       "token": "<access-token>",
       "refresh_token": "<refresh-token>"
   }
   ```

2. **Use Access Token in Requests:**  
   Include in headers:  
   ```
   Authorization: Bearer <access-token>
   ```

3. **Refresh Token if Access Token Expires:**  
   ```
   POST /api/refresh
   ```
   **Headers:**  
   ```
   Authorization: Bearer <refresh-token>
   ```
   **Response:**  
   ```json
   {
       "token": "<new-access-token>"
   }
   ```

4. **Revoke a Refresh Token:**  
   ```
   POST /api/revoke
   ```
   **Headers:**  
   ```
   Authorization: Bearer <refresh-token>
   ```

---

## **Error Handling Summary**  

| **Endpoint**          | **Error**                 | **Response** |
|-----------------------|---------------------------|--------------|
| `/api/users/{ID}`     | User not found            | 404          |
| `/api/login`          | Incorrect credentials     | 401          |
| `/api/refresh`        | Invalid refresh token     | 401          |
| `/api/chirps/{ID}`    | Chirp not found           | 404          |
| `/api/chirps/{ID}`    | Unauthorized delete       | 401          |

---
