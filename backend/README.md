# PTracker Backend

## Database Schema

### Database type

- **Database system:** PostgreSQL

### `sessions`

| Column             | Type        | Constraints / Notes               |
| ------------------ | ----------- | --------------------------------- |
| id                 | UUID        | PK, default `gen_random_uuid()`   |
| user_id            | UUID        | FK → users(id), ON DELETE CASCADE |
| refresh_token_hash | TEXT        | UNIQUE, NOT NULL                  |
| user_agent         | TEXT        |                                   |
| ip_address         | INET        |                                   |
| device_name        | TEXT        |                                   |
| created_at         | TIMESTAMPTZ | NOT NULL, default `NOW()`         |
| last_active_at     | TIMESTAMPTZ | NOT NULL, default `NOW()`         |
| revoked_at         | TIMESTAMPTZ | NULL                              |
| expires_at         | TIMESTAMPTZ | NOT NULL                          |

**Indexes**

| Index Name           | Columns | Condition          |
| -------------------- | ------- | ------------------ |
| idx_sessions_user_id | user_id | —                  |
| idx_sessions_active  | user_id | revoked_at IS NULL |

### `projects`

| Column      | Type         | Constraints / Notes             |
| ----------- | ------------ | ------------------------------- |
| id          | UUID         | PK, default `gen_random_uuid()` |
| skills      | TEXT         |                                 |
| owner       | UUID         | FK → users(id)                  |
| name        | VARCHAR(255) | NOT NULL                        |
| description | TEXT         |                                 |
| created_at  | TIMESTAMPTZ  | default CURRENT_TIMESTAMP       |
| updated_at  | TIMESTAMPTZ  |                                 |
| deleted_at  | TIMESTAMPTZ  |                                 |

### Enum: `request_status`

| Possible Values |
| --------------- |
| Pending         |
| Accepted        |

### `join_requests`

| Column     | Type           | Constraints / Notes        |
| ---------- | -------------- | -------------------------- |
| project_id | UUID           | FK → projects(id)          |
| user_id    | UUID           | FK → users(id)             |
| status     | request_status | NOT NULL, default "Pendig" |
| created_at | TIMESTAMPTZ    | default CURRENT_TIMESTAMP  |
| updated_at | TIMESTAMPTZ    |                            |
| deleted_at | TIMESTAMPTZ    |                            |

**Indexes**

| Index Name               | Columns               |
| ------------------------ | --------------------- |
| ux_project_join_requests | (project_id, user_id) |

### Enum: `task_status`

| Possible Values |
| --------------- |
| Unassigned      |
| Ongoing         |
| Completed       |
| Abandoned       |

### `tasks`

| Column      | Type         | Constraints / Notes             |
| ----------- | ------------ | ------------------------------- |
| id          | UUID         | PK, default `gen_random_uuid()` |
| project_id  | UUID         | FK → projects(id)               |
| title       | VARCHAR(255) | NOT NULL                        |
| description | TEXT         |                                 |
| status      | task_status  | NOT NULL                        |
| created_at  | TIMESTAMPTZ  | default CURRENT_TIMESTAMP       |
| updated_at  | TIMESTAMPTZ  |                                 |
| deleted_at  | TIMESTAMPTZ  |                                 |

### Enum: `user_role`

| Possible Values |
| --------------- |
| Owner           |
| Member          |

### `roles`

| Column     | Type        | Constraints / Notes       |
| ---------- | ----------- | ------------------------- |
| project_id | UUID        | FK → projects(id)         |
| user_id    | UUID        | FK → users(id)            |
| role       | user_role   | NOT NULL                  |
| created_at | TIMESTAMPTZ | default CURRENT_TIMESTAMP |
| updated_at | TIMESTAMPTZ |                           |
| deleted_at | TIMESTAMPTZ |                           |

**Indexes**

| Index Name       | Columns               |
| ---------------- | --------------------- |
| ux_project_roles | (project_id, user_id) |

### `assignees`

| Column     | Type        | Constraints / Notes       |
| ---------- | ----------- | ------------------------- |
| project_id | UUID        | FK → projects(id)         |
| task_id    | UUID        | FK → tasks(id)            |
| user_id    | UUID        | FK → users(id)            |
| created_at | TIMESTAMPTZ | default CURRENT_TIMESTAMP |
| updated_at | TIMESTAMPTZ |                           |
| deleted_at | TIMESTAMPTZ |                           |

**Indexes**

| Index Name               | Columns                        |
| ------------------------ | ------------------------------ |
| ux_project_task_assignee | (project_id, task_id, user_id) |

### `comments`

| Column     | Type        | Constraints / Notes             |
| ---------- | ----------- | ------------------------------- |
| id         | UUID        | PK, default `gen_random_uuid()` |
| project_id | UUID        | FK → projects(id)               |
| task_id    | UUID        | FK → tasks(id)                  |
| user_id    | UUID        | FK → users(id)                  |
| content    | TEXT        | NOT NULL                        |
| created_at | TIMESTAMPTZ | default CURRENT_TIMESTAMP       |
| updated_at | TIMESTAMPTZ |                                 |
| deleted_at | TIMESTAMPTZ |                                 |

### Relationships

- **users to projects**: one-many
- **projects to tasks**: one-many
- **users to roles**: one-many
- **projects to roles**: one-many
- **users to assignees**: one-many
- **projects to assignees**: one-many
- **tasks to assignees**: one-many
- **projects to comments**: one-many
- **users to comments**: one-many
- **tasks to comments**: one-many

### Database Diagram

[Database Diagram PDF](./db/schema_design.pdf)

### Migration

For the first time running server, you would need the necessary tables. There is already a `migrations` folder containing all the necessary changes. You can perform the migration with the following command:

> `n` is the number of migrations.

```sh
migrate -source file://migrations -database postgres://[user]:[password]@[host]:[port]/[database]?sslmode=disable up n
```

For some reason if the tables are not created or modification is needed, use the following command to take the tables down.

```sh
migrate -source file://migrations -database postgres://[user]:[password]@[host]:[port]/[database]?sslmode=disable down n
```

And fix the migrations, then try again.

Before running the `migrate` command, you would need the CLI. You can install the CLI following the official golang-migrate [instructions](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate#unversioned).

## API Endpoints

**All endpoints have `/api` prefix**

## **Auth**

Authentication is managed by Keycloak.

Endpoints:

1. `/auth/login`: Start keycloak login process.
2. `/auth/callback`: Keycloak will redirect authorization code to this endpoint where sessions cookie is created.
3. `/auth/refresh`: Refresh token with keycloak.
4. `/auth/logout`: Logout by cleaning cookie, session data.

### **Tokens and session_id**

1. User tries to login.
2. Backend redirects to Keycloak login UI.
3. User logs in successfully.
4. Backend receives the `access_token`, `id_token` and `refresh_token`.
5. Backend creates a new session (entry in the `sessions` tables) in the database with hashed `refresh_token`. Backend also saves the `access_token` with `session_id` in-memory.
6. Backend sends `session_id` as HTTPOnly, secure cookie.

### **Per Request Check**

For every request from frontend,

1. Backend checks the `session_id`, looks up for the `access_token`.
2. If `access_token` verification is successful - get the `sub` field.
3. Fetch user details from `users` table using `sub` and save the user in the request context.
4. If `access_token` verification fails return `401` response.
5. Retry with `/refresh_token` endpoint to refresh the access token using the hashed `refresh_token` to get new `access_token` from keycloak.
6. Send the new `access_token` to frontend.

### **refresh_token expires**

In case `refresh_token` expires, we

- log out the user
- revoke the session
- ask user to login again

## **Projects**

API supports following actions for project:

- Create a new project
- Get a project with `project_id`
- Update a project name/description
- Delete a project. It is a soft delete.

### **POST /projects**

Create a project. Any authenticated user can create a project. If `name` is missing or _empty_ in the payload, user will get `invalid_body` error. If `description` is missing or _empty_ in the payload, it is assumed to be `""` and same for `skills` field. No extra field is allowed in the payload, otherwise it will throw an `invalid_body` error.

**Payload:**

```json
{
  "name": "PROJECT A",
  "description": "Blah blah blah",
  "skills": "Java"
}
```

**Response**

`201 Created`

```json
{
  "status": "success",
  "data": {
    "id": "a3f8b1ce-92d4-4c1b-b8df-5e71a2c6c901",
    "name": "PROJECT A",
    "description": "Blah blah blah",
    "skills": "C++, Python",
    "created_at": "2025-12-02T10:40:00Z",
    "updated_at": null
  }
}
```

`400 Bad Request`

```json
{
  "status": "error",
  "error": {
    "id": "invalid_body",
    "message": "Project 'name' is missing"
  }
}
```

`401 Unauthorized`

```json
{
  "status": "error",
  "error": {
    "id": "unauthorized",
    "message": "User is not authorized"
  }
}
```

### **GET /projects**

Get all projects with details of the project and the task statistics inside them. The request query takes `page` and `limit`. E.g. `/projects?page=1&limit=10`. Without the query parameters the API assumes `page=1` and `limit=10`.

**Response**

`200 OK`

```json
  "status": "success",
  "data": {
    "projects":[
      {
        "id": "a3f8b1ce-92d4-4c1b-b8df-5e71a2c6c901",
        "name": "PROJECT A",
        "description": "Blah blah blah",
        "skills": "C++, Python",
        "role": "Owner",
        "unassigned_tasks": 1,
        "ongoing_tasks": 2,
        "completed_tasks": 0,
        "abandoned_tasks": 0,
        "created_at": "2025-12-02T10:40:00Z",
        "updated_at": "2025-12-02T10:40:00Z"
      },
      {
        "id": "edf8b2da-10e6-4c1b-b8df-5e71a2c6c901",
        "name": "PROJECT B",
        "description": "Blah blah blah",
        "skills": "C, Java",
        "role": "Member",
        "unassigned_tasks": 0,
        "ongoing_tasks": 3,
        "completed_tasks": 1,
        "abandoned_tasks": 0,
        "created_at": "2025-12-02T10:40:00Z",
        "updated_at": "2025-12-02T10:40:00Z"
      },
      // ...
    ],
    "page": 4,
    "limit": 10,
    "total": 50,
    "has_next": true
  }
```

`401 Unauthorized`

```json
{
  "status": "error",
  "error": {
    "id": "unauthorized",
    "message": "User is not authorized"
  }
}
```

### **GET /projects/<project_id>**

Get project details of project ID `project_id`. Project owners and members can get the project details. If user is not part of the project, they will get `access_denied` error. If the project with ID is not found, they will get `resource_not_found` error.

**Response:**

`200 OK`

```json
{
  "status": "success",
  "data": {
    "id": "a3f8b1ce-92d4-4c1b-b8df-5e71a2c6c901",
    "name": "PROJECT A",
    "code": "XYZ",
    "description": "Blah blah blah",
    "owner": {
      "id": "f92a07cd-1e6f-4df3-9e70-9ad94f0d0ed3",
      "username": "USER_A"
    },
    "created_at": "2025-12-02T10:40:00Z",
    "updated_at": "2025-12-02T10:40:00Z"
  }
}
```

`401 Unauthorized`

```json
{
  "status": "error",
  "error": {
    "id": "unauthorized",
    "message": "User is not authorized"
  }
}
```

`403 Forbidden`

```json
{
  "status": "error",
  "error": {
    "id": "access_denied",
    "message": "User is not a member of the project"
  }
}
```

`404 Not Found`

```json
{
  "status": "error",
  "error": {
    "id": "resource_not_found",
    "message": "Project not found"
  }
}
```

### **PATCH /projects/<project_id>**

Update the name or description of the project with ID `project_id`. Request without any payload properties will result in `invalid_body` error. Payload should have atleast one property. If user is not the project owner then they will get `access_denied` error. If the project with ID is not found, they will get `resource_not_found` error.

**Payload:**

```json
{
  "name": "PROJECT A (UPDATED)",
  "description": "Blah blah blah (UPDATED)"
}
```

**Response**

`200 OK`

```json
{
  "status": "success",
  "data": {
    "id": "a3f8b1ce-92d4-4c1b-b8df-5e71a2c6c901",
    "code": "XYZ",
    "name": "PROJECT A (UPDATED)",
    "description": "Blah blah blah (UPDATED)",
    "created_at": "2025-12-02T10:40:00Z",
    "updated_at": "2025-12-02T10:40:00Z"
  }
}
```

`400 Bad Request`

```json
{
  "status": "error",
  "error": {
    "id": "invalid_body",
    "message": "Payload can't be empty"
  }
}
```

`401 Unauthorized`

```json
{
  "status": "error",
  "error": {
    "id": "unauthorized",
    "message": "User is not authorized"
  }
}
```

`403 Forbidden`

```json
{
  "status": "error",
  "error": {
    "id": "access_denied",
    "message": "User does not have permission to update the project"
  }
}
```

`404 Not Found`

```json
{
  "status": "error",
  "error": {
    "id": "resource_not_found",
    "message": "Project not found"
  }
}
```

### **DELETE /projects/<project_id>**

Delete a project with project ID `project_id`. Only project owner can delete a project. If user is not a project owner, they will get `access_denied` error. If the project with ID is not found, they will get `resource_not_found` error.

**Response**

`204 No Content`

`401 Unauthorized`

```json
{
  "status": "error",
  "error": {
    "id": "unauthorized",
    "message": "User is not authorized"
  }
}
```

`403 Forbidden`

```json
{
  "status": "error",
  "error": {
    "id": "access_denied",
    "message": "User does not have permission to update the project"
  }
}
```

`404 Not Found`

```json
{
  "status": "error",
  "error": {
    "id": "resource_not_found",
    "message": "Project not found"
  }
}
```

## **Tasks**

API has endpoints for the following actions:

- Create a new task for the project
- Get all the tasks in the project
- Get a particular task with `task_id`
- Update a task title/description
- Update a task status
- Delete a task. It is a soft delete.

### **POST /projects/<project_id>/tasks**

Create a task. Only project owner can create task for a project. Initially the task will have `status` `Unassigned`. If user is not the project owner, they will get `access_denied` error. If `title` is missing or _empty_ in the payload, user will get `invalid_body` error. If `description` is missing or _empty_ in the payload, it is assumed to be `""`. If the project with ID is not found, they will get `resource_not_found` error.

**Payload:**

```json
{
  "title": "TASK A",
  "description": "Blah blah"
}
```

**Response:**

`200 OK`

```json
{
  "status": "success",
  "data": {
    "id": "c14b5e6a-71bf-4d8d-ae93-8c2420f4c5b2",
    "title": "Task A title",
    "description": "Blah blah blah",
    "status": "Unassigned",
    "assignees": [],
    "created_at": "2025-12-02T10:40:00Z",
    "updated_at": null
  }
}
```

`400 Bad Request`

```json
{
  "status": "error",
  "error": {
    "id": "invalid_body",
    "message": "Payload property 'name' is missing or empty"
  }
}
```

`401 Unauthorized`

```json
{
  "status": "error",
  "error": {
    "id": "unauthorized",
    "message": "User is not authorized"
  }
}
```

`403 Forbidden`

```json
{
  "status": "error",
  "error": {
    "id": "access_denied",
    "message": "User does not have permission to add task in the project"
  }
}
```

`404 Not Found`

```json
{
  "status": "error",
  "error": {
    "id": "resource_not_found",
    "message": "Project not found"
  }
}
```

### **GET /projects/<project_id>/tasks**

Get all tasks under the project. If user is not project owner or project member, they will get `access_denied` error. If the project with ID is not found, they will get `resource_not_found` error.

API endpoint Supports pagination.

**Query:** `?page=4&limit=10`

**Response:**

`200 OK`

```json
{
  "status": "success",
  "data": [
    {
      "id": "c14b5e6a-71bf-4d8d-ae93-8c2420f4c5b2",
      "title": "Task A title",
      "description": "Blah blah blah",
      "status": "Unassigned",
      "assignees": [],
      "created_at": "2025-12-02T10:40:00Z",
      "updated_at": null
    },
    {
      "id": "9c379d84-97c1-4a19-bdd7-d64fae91ab15",
      "title": "Task B title",
      "description": "Blah blah blah",
      "status": "Ongoing",
      "assignees": [
        {
          "id": "f92a07cd-1e6f-4df3-9e70-9ad94f0d0ed3",
          "username": "USER_A"
        }
      ],
      "created_at": "2025-12-02T10:40:00Z",
      "updated_at": "2025-12-02T10:40:00Z"
    }
    // ...
  ],
  "page": 4,
  "limit": 10,
  "total": 50,
  "has_next": true
}
```

`401 Unauthorized`

```json
{
  "status": "error",
  "error": {
    "id": "unauthorized",
    "message": "User is not authorized"
  }
}
```

`403 Forbidden`

```json
{
  "status": "error",
  "error": {
    "id": "access_denied",
    "message": "User does not have permission to view tasks in the project"
  }
}
```

`404 Not Found`

```json
{
  "status": "error",
  "error": {
    "id": "resource_not_found",
    "message": "Project not found"
  }
}
```

### **GET /projects/<project_id>/tasks/<task_id>**

Get task details of project task with ID `task_id`. If user is not project owner or project member, they will get `access_denied` error. If the project `project_id` or task `task_id` is not found, they will get `resource_not_found` error.

**Response:**

`200 OK`

```json
{
  "status": "success",
  "data": {
    "id": "c14b5e6a-71bf-4d8d-ae93-8c2420f4c5b2",
    "title": "Task A title",
    "description": "Blah blah blah",
    "status": "Unassigned",
    "assignees": [],
    "created_at": "2025-12-02T10:40:00Z",
    "updated_at": null
  }
}
```

`401 Unauthorized`

```json
{
  "status": "error",
  "error": {
    "id": "unauthorized",
    "message": "User is not authorized"
  }
}
```

`403 Forbidden`

```json
{
  "status": "error",
  "error": {
    "id": "access_denied",
    "message": "User does not have permission to view the task in the project"
  }
}
```

`404 Not Found`

```json
{
  "status": "error",
  "error": {
    "id": "resource_not_found",
    "message": "Project/Task not found"
  }
}
```

### **PATCH /projects/<project_id>/tasks/<task_id>**

Update the name or description of the project task with ID `task_id`. Request without any payload properties will result in `invalid_body` error. Payload should have atleast one property. If user is not the project owner then they will get `access_denied` error. If the project or task with the given ID is not found, they will get `resource_not_found` error.

**Payload:**

```json
{
  "title": "A different title",
  "description": "A different description"
}
```

**Response**

`200 OK`

```json
{
  "status": "success",
  "data": {
    "id": "c14b5e6a-71bf-4d8d-ae93-8c2420f4c5b2",
    "title": "A different title",
    "description": "A different description",
    "status": "Unassigned",
    "assignees": [],
    "created_at": "2025-12-02T10:40:00Z",
    "updated_at": "2025-12-02T10:40:00Z"
  }
}
```

`400 Bad Request`

```json
{
  "status": "error",
  "error": {
    "id": "invalid_body",
    "message": "Payload can't be empty"
  }
}
```

`401 Unauthorized`

```json
{
  "status": "error",
  "error": {
    "id": "unauthorized",
    "message": "User is not authorized"
  }
}
```

`403 Forbidden`

```json
{
  "status": "error",
  "error": {
    "id": "access_denied",
    "message": "User does not have permission to update the project task"
  }
}
```

`404 Not Found`

```json
{
  "status": "error",
  "error": {
    "id": "resource_not_found",
    "message": "Project/Task not found"
  }
}
```

### **PATCH /projects/<project_id>/tasks/<task_id>/status**

Update the `status` of the project task `task_id`. `status` has the following values:

- _Unassigned_: A task with no assignee yet. When a new task is created they are in this `status`.
- _Ongoing_: When project owner assigns someone to the task. If project owner removes the assignee then the task will go back to _Unassigned_ `status` again.
- _Completed_: When the task is completed, project owner will mark the task as _Completed_. Project owner can bring back the task to _Ongoing_ if he thinks necessary.
- _Abandoned_: Project owner can also mark a task as _Abandoned_ due to various reasons. The task might become irrelevant or not necessary for the project. Assignee was not found for the task. Or, Assignee was not able to continue the task for any reason. An _Abandoned_ task may again become _Ongoing_ if someone volunteer to continue that or project owner deems the task necessary for the project again.

If payload `status` is not any of the above, it will return `invalid_body` error. If user is not the project owner then they will get `access_denied` error. If the project or task with the given ID is not found, they will get `resource_not_found` error.

**Payload**

```json
{
  "status": "Ongoing" // "Unassigned"/"Completed"/"Abandoned"
}
```

**Response**

`200 OK`

```json
{
  "status": "success",
  "data": {
    "id": "c14b5e6a-71bf-4d8d-ae93-8c2420f4c5b2",
    "title": "A different title",
    "description": "A different description",
    "status": "Ongoing",
    "assignees": [
      {
        "id": "f92a07cd-1e6f-4df3-9e70-9ad94f0d0ed3",
        "username": "USER_A"
      }
    ],
    "created_at": "2025-12-02T10:40:00Z",
    "updated_at": "2025-12-02T10:40:00Z"
  }
}
```

`400 Bad Request`

```json
{
  "status": "error",
  "error": {
    "id": "invalid_body",
    "message": "Payload contains invalid 'status' value"
  }
}
```

`401 Unauthorized`

```json
{
  "status": "error",
  "error": {
    "id": "unauthorized",
    "message": "User is not authorized"
  }
}
```

`403 Forbidden`

```json
{
  "status": "error",
  "error": {
    "id": "access_denied",
    "message": "User does not have permission to change task status"
  }
}
```

`404 Not Found`

```json
{
  "status": "error",
  "error": {
    "id": "resource_not_found",
    "message": "Project/Task not found"
  }
}
```

### **DELETE /projects/<project_id>/tasks/<task_id>**

Delete a task in project `project_id` with ID `task_id`. Only project owner can delete a task. If user is not a project owner, they will get `access_denied` error. If the project with ID is not found, they will get `resource_not_found` error.

**Response**

`204 No Content`

`401 Unauthorized`

```json
{
  "status": "error",
  "error": {
    "id": "unauthorized",
    "message": "User is not authorized"
  }
}
```

`403 Forbidden`

```json
{
  "status": "error",
  "error": {
    "id": "access_denied",
    "message": "User does not have permission to delete a task"
  }
}
```

`404 Not Found`

```json
{
  "status": "error",
  "error": {
    "id": "resource_not_found",
    "message": "Project/Task not found"
  }
}
```

## **Assignees**

API supports following assignee related actions:

- Add new assignee(s) to the task.
- Remove an assignee from the task. It is a soft delete.

API does not support the following assignee actions:

- Get all assignees. You can get all assignees in the `GET /projects/<project_id>/tasks/<task_id>` endpoint.

### **POST /projects/<project_id>/tasks/<task_id>/assignees**

Add assignee(s) to a project task. The payload contains list of assignees with their user_id. Empty list will return `invalid_body` error. If user is not the project owner then they will get `access_denied` error. If the project or task with the given ID is not found, they will get `resource_not_found` error.

**Payload:**

```json
{
  "assignees": ["9c379d84-97c1-4a19-bdd7-d64fae91ab15"]
}
```

**Response**

`200 OK`

```json
{
  "status": "success",
  "data": {
    "id": "c14b5e6a-71bf-4d8d-ae93-8c2420f4c5b2",
    "title": "Task A",
    "description": "Task A description",
    "status": "Ongoing",
    "assignees": [
      {
        "id": "f92a07cd-1e6f-4df3-9e70-9ad94f0d0ed3",
        "username": "USER_A"
      },
      {
        "id": "9c379d84-97c1-4a19-bdd7-d64fae91ab15",
        "username": "USER_C"
      }
    ],
    "created_at": "2025-12-02T10:40:00Z",
    "updated_at": "2025-12-02T10:40:00Z"
  }
}
```

`400 Bad Request`

```json
{
  "status": "error",
  "error": {
    "id": "invalid_body",
    "message": "Payload can't be empty"
  }
}
```

`401 Unauthorized`

```json
{
  "status": "error",
  "error": {
    "id": "unauthorized",
    "message": "User is not authorized"
  }
}
```

`403 Forbidden`

```json
{
  "status": "error",
  "error": {
    "id": "access_denied",
    "message": "User does not have permission to update the project task"
  }
}
```

`404 Not Found`

```json
{
  "status": "error",
  "error": {
    "id": "resource_not_found",
    "message": "Project/Task not found"
  }
}
```

### **DELETE /projects/<project_id>/tasks/<task_id>/assignees/<assignee_id>**

Remove an assignee `assignee_id` from a task. Only project owner can remove an assignee. If user is not a project owner, they will get `access_denied` error. If the project or task with ID is not found, they will get `resource_not_found` error.

**Response**

`200 OK`

```json
{
  "status": "success",
  "data": {
    "id": "c14b5e6a-71bf-4d8d-ae93-8c2420f4c5b2",
    "title": "Task A",
    "description": "Task A description",
    "status": "Ongoing",
    "assignees": [
      {
        "id": "f92a07cd-1e6f-4df3-9e70-9ad94f0d0ed3",
        "username": "USER_A"
      }
    ],
    "created_at": "2025-12-02T10:40:00Z",
    "updated_at": "2025-12-02T10:40:00Z"
  }
}
```

`401 Unauthorized`

```json
{
  "status": "error",
  "error": {
    "id": "unauthorized",
    "message": "User is not authorized"
  }
}
```

`403 Forbidden`

```json
{
  "status": "error",
  "error": {
    "id": "access_denied",
    "message": "User does not have permission to delete a task"
  }
}
```

`404 Not Found`

```json
{
  "status": "error",
  "error": {
    "id": "resource_not_found",
    "message": "Project/Task not found"
  }
}
```

## **Members**

API supports following actions for members:

- Get all members of the project.
- Remove a member of the project. It is a soft delete.

### **GET /projects/<project_id>/members**

List all members of a project with ID `project_id`. If user is not project owner or project member, they will get `access_denied` error. If the project with ID is not found, they will get `resource_not_found` error.

API endpoint Supports pagination.

**Query:** `?page=4&limit=10`

**Response:**

`200 OK`

```json
{
  "status": "success",
  "data": [
    {
      "id": "f92a07cd-1e6f-4df3-9e70-9ad94f0d0ed3",
      "username": "USER_A",
      "role": "Owner",
      "created_at": "2025-12-02T10:40:00Z",
      "updated_at": "2025-12-02T10:40:00Z"
    },
    {
      "id": "84f9c26e-154a-4d93-946e-64c082031273",
      "username": "USER_B",
      "role": "Member",
      "created_at": "2025-12-02T10:40:00Z",
      "updated_at": "2025-12-02T10:40:00Z"
    }
    // ...
  ],
  "page": 4,
  "limit": 10,
  "total": 50,
  "has_next": true
}
```

`401 Unauthorized`

```json
{
  "status": "error",
  "error": {
    "id": "unauthorized",
    "message": "User is not authorized"
  }
}
```

`403 Forbidden`

```json
{
  "status": "error",
  "error": {
    "id": "access_denied",
    "message": "User does not have permission to view members in the project"
  }
}
```

`404 Not Found`

```json
{
  "status": "error",
  "error": {
    "id": "resource_not_found",
    "message": "Project not found"
  }
}
```

### **DELETE /projects/<project_id>/members/<member_id>**

Remove a member `member_id` from the project `project_id`. It is a soft delete. Only project owner can remove a member. If user is not a project owner, they will get `access_denied` error. If the project or member with ID is not found, they will get `resource_not_found` error.

**Response**

`204 No Content`

`401 Unauthorized`

```json
{
  "status": "error",
  "error": {
    "id": "unauthorized",
    "message": "User is not authorized"
  }
}
```

`403 Forbidden`

```json
{
  "status": "error",
  "error": {
    "id": "access_denied",
    "message": "User does not have permission to remove a member"
  }
}
```

`404 Not Found`

```json
{
  "status": "error",
  "error": {
    "id": "resource_not_found",
    "message": "Project/Member not found"
  }
}
```

## **Comments**

API supports following actions for comments:

- Add a comment to a task
- Get all comments of a task
- Update a comment
- Delete a comment

### **POST /projects/<project_id>/tasks/<task_id>/comments**

Add a comment in a task `task_id` of a project `project_id`. Only project members can comment. If user is not a project member, they will get `access_denied` error. Payload must include `content` with some text otherwise it will return `invalid_body` error. If the project or task with ID is not found, they will get `resource_not_found` error.

**Payload:**

```json
{
  "content": "Blah blah blah"
}
```

**Response**

`201 OK`

```json
{
  "status": "success",
  "data": {
    "id": "5b4e1fca-d89c-4e3e-8d73-94dcd77b3ef3",
    "content": "Blah blah blah",
    "commenter": {
      "id": "f92a07cd-1e6f-4df3-9e70-9ad94f0d0ed3",
      "username": "USER_A"
    },
    "created_at": "2025-12-02T10:40:00Z",
    "updated_at": null
  }
}
```

`400 Bad Request`

```json
{
  "status": "error",
  "error": {
    "id": "invalid_body",
    "message": "Payload property 'content' is missing or empty"
  }
}
```

`401 Unauthorized`

```json
{
  "status": "error",
  "error": {
    "id": "unauthorized",
    "message": "User is not authorized"
  }
}
```

`403 Forbidden`

```json
{
  "status": "error",
  "error": {
    "id": "access_denied",
    "message": "User does not have permission to add comment to the task"
  }
}
```

`404 Not Found`

```json
{
  "status": "error",
  "error": {
    "id": "resource_not_found",
    "message": "Project/Task not found"
  }
}
```

### **GET /projects/<project_id>/tasks/<task_id>/comments**

Get task comments in paginated form. If user is not project owner or project member, they will get `access_denied` error. If the project or task with ID is not found, they will get `resource_not_found` error.

API endpoint Supports pagination.

**Query:** `?page=4&limit=10`

**Response:**

`200 OK`

```json
{
  "status": "success",
  "data": [
    {
      "id": "b0c4a9ac-8d63-489b-af83-91b6bcfe77e0",
      "content": "Blah blah",
      "commenter": {
        "id": "f92a07cd-1e6f-4df3-9e70-9ad94f0d0ed3",
        "username": "USER_A"
      },
      "created_at": "2025-12-02T10:40:00Z",
      "updated_at": null
    },
    {
      "id": "e3b1f3de-c5ef-4608-90ee-2b834c0bbd02",
      "content": "Blah blah",
      "commenter": {
        "id": "84f9c26e-154a-4d93-946e-64c082031273",
        "username": "USER_B"
      },
      "created_at": "2025-12-02T10:40:00Z",
      "updated_at": "2025-12-02T10:40:00Z"
    }
    // ...
  ],
  "page": 4,
  "limit": 10,
  "total": 50,
  "has_next": true
}
```

`401 Unauthorized`

```json
{
  "status": "error",
  "error": {
    "id": "unauthorized",
    "message": "User is not authorized"
  }
}
```

`403 Forbidden`

```json
{
  "status": "error",
  "error": {
    "id": "access_denied",
    "message": "User does not have permission to view comments in the task"
  }
}
```

`404 Not Found`

```json
{
  "status": "error",
  "error": {
    "id": "resource_not_found",
    "message": "Project/Task not found"
  }
}
```

### **PATCH /projects/<project_id>/tasks/<task_id>/comments/<comment_id>**

Update the content of the comment with ID `comment_id`. Request without any payload properties will result in `invalid_body` error. Payload should have non-empty `content` value. If user is not the project owner or commenter then they will get `access_denied` error. If the project/task/comment with the given ID is not found, they will get `resource_not_found` error.

**Payload:**

```json
{
  "content": "An updated comment"
}
```

**Response**

`200 OK`

```json
{
  "status": "success",
  "data": {
    "id": "e3b1f3de-c5ef-4608-90ee-2b834c0bbd02",
    "content": "Blah blah",
    "commenter": {
      "id": "84f9c26e-154a-4d93-946e-64c082031273",
      "username": "USER_B"
    },
    "created_at": "2025-12-02T10:40:00Z",
    "updated_at": "2025-12-02T10:40:00Z"
  }
}
```

`400 Bad Request`

```json
{
  "status": "error",
  "error": {
    "id": "invalid_body",
    "message": "Payload 'content' is missing or empty"
  }
}
```

`401 Unauthorized`

```json
{
  "status": "error",
  "error": {
    "id": "unauthorized",
    "message": "User is not authorized"
  }
}
```

`403 Forbidden`

```json
{
  "status": "error",
  "error": {
    "id": "access_denied",
    "message": "User does not have permission to update the comment"
  }
}
```

`404 Not Found`

```json
{
  "status": "error",
  "error": {
    "id": "resource_not_found",
    "message": "Project/Task/Comment not found"
  }
}
```

### **DELETE /projects/<project_id>/tasks/<task_id>/comments/<comment_id>**

Delete a comment in project `project_id` task `task_id` with comment ID `comment_id`. Only project owner or commenter can delete a comment. If user is not a project owner or the commenter, they will get `access_denied` error. If the project/task/comment with ID is not found, they will get `resource_not_found` error.

**Response**

`204 No Content`

`401 Unauthorized`

```json
{
  "status": "error",
  "error": {
    "id": "unauthorized",
    "message": "User is not authorized"
  }
}
```

`403 Forbidden`

```json
{
  "status": "error",
  "error": {
    "id": "access_denied",
    "message": "User does not have permission to delete a comment"
  }
}
```

`404 Not Found`

```json
{
  "status": "error",
  "error": {
    "id": "resource_not_found",
    "message": "Project/Task/Comment not found"
  }
}
```

## Rules

| **Action**            | **Allowed Roles**              |
| --------------------- | ------------------------------ |
| Update project        | Project owner                  |
| Add task              | Project owner                  |
| Update task           | Project owner, Assignee        |
| Add assignee          | Project owner                  |
| Change assignee       | Project owner                  |
| Add comment           | Project owner, Project members |
| Delete comment        | Project owner, Commenter       |
| Remove assignee       | Project owner                  |
| Remove project member | Project owner                  |
