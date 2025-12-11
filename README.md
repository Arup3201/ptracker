# PTracker - A collaborative project management application

Web application for tracking and managing projects between multiple people.

**Tech Stack**

Go, PostgreSQL, KeyCloak, ReactJS, TypeScript

## Set up

Change directory to `backend` first to startup the necessary database applications.

```sh
cd backend
```

Use docker-compose to run the instances in the background (`-d`).

```sh
docker compose up -d
```

It should start the necessary databases.

Then you can run the go server with,

```sh
go run main.go
```

It should start the server as well as create all the database tables necessary for the app to run through migrations.

For frontend, open new terminal and change directory to frontend.

```sh
cd frontend
```

And run the following command to start the frontend at `localhost:5173`,

```sh
npm run dev
```
