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

Before starting the server, make sure you have `.env` file inside the `backend` folder.

For dev, you can COPY the following `.env` and it should work fine.

**.env**

```sh
HOST=localhost
PORT=8081
HOME_URL=http://localhost:5173
ENCRYPTION_SECRET=6dee83baf4eb0bea602a632c4eed37ff
PG_HOST=localhost
PG_USER=postgres
PG_PORT=5432
PG_PASS=1234
PG_DB=ptracker
KC_URL=http://localhost:8080
KC_REALM=ptracker
KC_CLIENT_ID=api
KC_CLIENT_SECRET=cp50avHQeX18cESEraheJvr3RhUBMq2A
KC_REDIRECT_URI=http://localhost:8081/api/v1/auth/callback
```

Export the environment variables to your shell(Linux).

```sh
export (cat .env | xargs)
```

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
