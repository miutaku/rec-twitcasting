# API Usage Guide

This document provides an overview of the API endpoints implemented in the provided Go code. The server runs on `localhost:8888`, and all endpoints return JSON responses.

---

## Table of Contents

1. [List Casting Users](#1-list-casting-users)
2. [Add Casting User](#2-add-casting-user)
3. [Delete Casting User](#3-delete-casting-user)
4. [Check Recording State](#4-check-recording-state)
5. [Update Recording State](#5-update-recording-state)

---

### 1. List Casting Users

**Endpoint:** `/list-casting-users`

**Method:** `GET`

**Description:** Retrieves a list of all casting users from the database.

**Curl Example:**
```bash
curl "http://localhost:8888/list-casting-users"
```

**Response:**
```json
[
    {
        "recording_state": true,
        "action_date_time": "2025-01-01T12:00:00Z",
        "action": "listCastingUser",
        "target_username": "user1"
    },
    {
        "recording_state": false,
        "action_date_time": "2025-01-01T12:01:00Z",
        "action": "listCastingUser",
        "target_username": "user2"
    }
]
```

---

### 2. Add Casting User

**Endpoint:** `/add-casting-user`

**Method:** `GET`

**Description:** Adds a new user to the database.

**Query Parameters:**
- `username` (required): The username of the new user.

**Curl Example:**
```bash
curl "http://localhost:8888/add-casting-user?username=newuser"
```

**Response (Success):**
```json
{
    "action_date_time": "2025-01-01T12:00:00Z",
    "action": "addCastingUser",
    "target_username": "newuser"
}
```

**Response (User Already Exists):**
```json
{
    "error": "User already exists",
    "action_date_time": "2025-01-01T12:00:00Z",
    "action": "addCastingUser",
    "target_username": "newuser"
}
```

---

### 3. Delete Casting User

**Endpoint:** `/del-casting-user`

**Method:** `GET`

**Description:** Deletes an existing user from the database.

**Query Parameters:**
- `username` (required): The username of the user to be deleted.

**Curl Example:**
```bash
curl "http://localhost:8888/del-casting-user?username=user1"
```

**Response:**
```json
{
    "action_date_time": "2025-01-01T12:00:00Z",
    "action": "deleteCastingUser",
    "target_username": "user1"
}
```

---

### 4. Check Recording State

**Endpoint:** `/check-recording-state`

**Method:** `GET`

**Description:** Retrieves the recording state of a specific user.

**Query Parameters:**
- `username` (required): The username whose recording state is to be checked.

**Curl Example:**
```bash
curl "http://localhost:8888/check-recording-state?username=user1"
```

**Response (Success):**
```json
{
    "recording_state": true,
    "action_date_time": "2025-01-01T12:00:00Z",
    "action": "checkRecordingState",
    "target_username": "user1"
}
```

**Response (User Not Found):**
```text
User not found
```

---

### 5. Update Recording State

**Endpoint:** `/update-recording-state`

**Method:** `GET`

**Description:** Updates the recording state of a specific user.

**Query Parameters:**
- `username` (required): The username whose recording state is to be updated.
- `recording_state` (required): The new recording state (`true` or `false`).

**Curl Example:**
```bash
curl "http://localhost:8888/update-recording-state?username=user1&recording_state=true"
```

**Response:**
```json
{
    "action_date_time": "2025-01-01T12:00:00Z",
    "action": "updateRecordingState",
    "target_username": "user1"
}
