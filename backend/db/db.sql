CREATE TABLE IF NOT EXISTS "users" (
	"id" VARCHAR(255) NOT NULL UNIQUE,
	"username" VARCHAR(255) NOT NULL UNIQUE,
	PRIMARY KEY("id")
);




CREATE TABLE IF NOT EXISTS "projects" (
	"id" VARCHAR(255) NOT NULL UNIQUE,
	-- Human readable alphanumeric project code
	"code" VARCHAR(255) NOT NULL UNIQUE,
	"owner" VARCHAR(255) NOT NULL,
	"name" VARCHAR(255) NOT NULL,
	"description" VARCHAR(255),
	"created_at" TIMESTAMPTZ NOT NULL,
	"updated_at" TIMESTAMPTZ,
	"deleted_at" TIMESTAMPTZ,
	PRIMARY KEY("id")
);


COMMENT ON COLUMN "projects"."code" IS 'Human readable alphanumeric project code';


CREATE TABLE IF NOT EXISTS "tasks" (
	"id" VARCHAR(255) NOT NULL UNIQUE,
	"project_id" VARCHAR(255) NOT NULL,
	"title" VARCHAR(255) NOT NULL,
	"description" VARCHAR(255),
	-- "Unassigned" | "Ongoing" | "Completed" | "Abandoned"
	"status" VARCHAR(255) NOT NULL,
	"created_at" TIMESTAMPTZ NOT NULL,
	"updated_at" TIMESTAMPTZ,
	"deleted_at" TIMESTAMPTZ,
	PRIMARY KEY("id")
);


COMMENT ON COLUMN "tasks"."status" IS '"Unassigned" | "Ongoing" | "Completed" | "Abandoned"';


CREATE TABLE IF NOT EXISTS "roles" (
	"project_id" VARCHAR(255) NOT NULL,
	"user_id" VARCHAR(255) NOT NULL,
	-- "Owner" | "Member"
	"role" VARCHAR(255) NOT NULL,
	"created_at" TIMESTAMPTZ,
	"updated_at" TIMESTAMPTZ NOT NULL,
	"deleted_at" VARCHAR(255) NOT NULL
);


COMMENT ON COLUMN "roles"."role" IS '"Owner" | "Member"';


CREATE TABLE IF NOT EXISTS "assignees" (
	"project_id" VARCHAR(255) NOT NULL,
	"task_id" VARCHAR(255) NOT NULL,
	"user_id" VARCHAR(255) NOT NULL,
	"created_at" VARCHAR(255) NOT NULL,
	"updated_at" VARCHAR(255),
	"deleted_at" VARCHAR(255)
);




CREATE TABLE IF NOT EXISTS "comments" (
	"id" VARCHAR(255) NOT NULL UNIQUE,
	"project_id" VARCHAR(255) NOT NULL,
	"task_id" VARCHAR(255) NOT NULL,
	"user_id" VARCHAR(255) NOT NULL,
	"content" TEXT NOT NULL,
	"created_at" TIMESTAMPTZ NOT NULL,
	"updated_at" TIMESTAMPTZ,
	"deleted_at" TIMESTAMPTZ,
	PRIMARY KEY("id")
);



ALTER TABLE "projects"
ADD FOREIGN KEY("owner") REFERENCES "users"("id")
ON UPDATE NO ACTION ON DELETE NO ACTION;
ALTER TABLE "tasks"
ADD FOREIGN KEY("project_id") REFERENCES "projects"("id")
ON UPDATE NO ACTION ON DELETE NO ACTION;
ALTER TABLE "roles"
ADD FOREIGN KEY("user_id") REFERENCES "users"("id")
ON UPDATE NO ACTION ON DELETE NO ACTION;
ALTER TABLE "roles"
ADD FOREIGN KEY("project_id") REFERENCES "projects"("id")
ON UPDATE NO ACTION ON DELETE NO ACTION;
ALTER TABLE "users"
ADD FOREIGN KEY("id") REFERENCES "assignees"("user_id")
ON UPDATE NO ACTION ON DELETE NO ACTION;
ALTER TABLE "projects"
ADD FOREIGN KEY("id") REFERENCES "assignees"("project_id")
ON UPDATE NO ACTION ON DELETE NO ACTION;
ALTER TABLE "tasks"
ADD FOREIGN KEY("id") REFERENCES "assignees"("task_id")
ON UPDATE NO ACTION ON DELETE NO ACTION;
ALTER TABLE "projects"
ADD FOREIGN KEY("id") REFERENCES "comments"("project_id")
ON UPDATE NO ACTION ON DELETE NO ACTION;
ALTER TABLE "users"
ADD FOREIGN KEY("id") REFERENCES "comments"("user_id")
ON UPDATE NO ACTION ON DELETE NO ACTION;
ALTER TABLE "tasks"
ADD FOREIGN KEY("id") REFERENCES "comments"("task_id")
ON UPDATE NO ACTION ON DELETE NO ACTION;