# PTracker - A collaborative project management application

Web application for tracking and managing projects between multiple people.

**Tech Stack**

Go, PostgreSQL, KeyCloak, ReactJS, TypeScript

## Models

### **User**

| Field    | Description       |
| -------- | ----------------- |
| ID       | Unique identifier |
| Username | Name of the user  |

---

### **Project**

| Field            | Description              |
| ---------------- | ------------------------ |
| ID               | Unique identifier        |
| Name             | Project name             |
| Description      | Project description      |
| Owner (Username) | Owner of the project     |
| CreatedAt        | Timestamp of creation    |
| UpdatedAt        | Timestamp of last update |

---

### **Member**

| Field      | Description          |
| ---------- | -------------------- |
| User ID    | Reference to user    |
| Project ID | Reference to project |
| Joined At  | Timestamp of joining |
| Is Deleted | Soft delete flag     |

---

### **Task**

| Field       | Description                                            |
| ----------- | ------------------------------------------------------ |
| ID          | Unique identifier                                      |
| Project ID  | Reference to project                                   |
| Title       | Task title                                             |
| Description | Task description                                       |
| Status      | Unassigned / Ongoing / Completed / Abandoned / Deleted |
| Created At  | Timestamp of creation                                  |
| Updated At  | Timestamp of last update                               |

---

### **Assignee**

| Field       | Description             |
| ----------- | ----------------------- |
| Task ID     | Reference to task       |
| User ID     | Reference to user       |
| Assigned At | Timestamp of assignment |
| Is Deleted  | Soft delete flag        |

---

### **Comment**

| Field      | Description           |
| ---------- | --------------------- |
| ID         | Unique identifier     |
| Task ID    | Reference to task     |
| User ID    | Reference to user     |
| Content    | Comment text          |
| Created At | Timestamp of creation |

---

## API Endpoints

---

## **Auth (Managed by Keycloak)**

### **POST /register**

Create a user.

### **POST /login**

Login a user.

---

## **Projects**

### **POST /projects**

Create a project.

**Payload:**

```json
{
  "name": "PROJECT A",
  "description": "Blah blah blah"
}
```

---

### **POST /projects/code**

Join a project using a code.

**Payload:**

```json
{
  "code": "PROJECT_A",
  "user_id": "USER_A"
}
```

---

### **PATCH /projects/<project_id>**

Update project name or description.

**Payload:**

```json
{
  "name": "PROJECT A (UPDATED)",
  "description": "Blah blah blah (UPDATED)"
}
```

---

### **GET /projects/<project_id>**

Get project details.

**Response:**

```json
{
  "id": "PROJECT_A",
  "name": "PROJECT A",
  "description": "Blah blah blah",
  "owner": "USER_A"
}
```

---

## **Tasks**

### **POST /projects/<project_id>/tasks**

Create an **Unassigned** task.

**Payload:**

```json
{
  "title": "TASK A",
  "description": "Blah blah"
}
```

---

### **GET /projects/<project_id>/tasks**

Get all tasks under the project.
Supports pagination.

**Query:** `?page=4&limit=10`

**Response:**

```json
[
  {
    "id": "TASK_A",
    "title": "Task A title",
    "description": "Blah blah blah",
    "status": "Unassigned",
    "assignees": [],
    "created_at": "2025/12/02 10:40",
    "updated_at": "2025/12/04 22:30"
  }
]
```

---

### **PATCH /projects/<project_id>/tasks/<task_id>**

Update task fields and/or status.

**Payload:**

```json
{
  "title": "A different title",
  "description": "A different description",
  "status": "Completed"
}
```

---

### **GET /projects/<project_id>/tasks/<task_id>**

Get task details.

**Response:**

```json
{
  "id": "TASK_A",
  "title": "Task A title",
  "description": "Blah blah blah",
  "status": "Unassigned",
  "assignees": [],
  "created_at": "2025/12/02 10:40",
  "updated_at": "2025/12/04 22:30"
}
```

---

### **POST /projects/<project_id>/tasks/<task_id>/assignees**

Add assignee(s) to a project task.

**Payload:**

```json
{
  "members": ["USER_A"]
}
```

---

### **PATCH /projects/<project_id>/tasks/<task_id>/assignees**

Change assignee(s) of a project task.

**Payload:**

```json
{
  "members": ["USER_A"]
}
```

---

### **DELETE /projects/<project_id>/tasks/<task_id>/assignees/<assignee_id>**

Remove an assignee from a task.

---

## **Members**

### **GET /projects/<project_id>/members**

List all members of a project.

---

### **DELETE /projects/<project_id>/members/<member_id>**

Remove a member from the project.
