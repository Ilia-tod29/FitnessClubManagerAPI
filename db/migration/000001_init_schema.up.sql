CREATE TABLE "inventory" (
                             "id" bigserial PRIMARY KEY,
                             "name" varchar NOT NULL,
                             "image" varchar,
                             "created_at" timestamptz DEFAULT (now())
);

CREATE TABLE "users" (
                         "id" bigserial PRIMARY KEY,
                         "email" varchar UNIQUE NOT NULL,
                         "hashed_password" varchar NOT NULL,
                         "role" varchar NOT NULL DEFAULT 'user',
                         "suspended" boolean NOT NULL DEFAULT 'false',
                         "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "subscriptions" (
                                 "id" bigserial PRIMARY KEY,
                                 "user_id" bigint NOT NULL,
                                 "start_date" date NOT NULL,
                                 "end_date" date NOT NULL,
                                 "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "gallery" (
                           "id" bigserial PRIMARY KEY,
                           "image" varchar,
                           "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE INDEX ON "inventory" ("name");

CREATE INDEX ON "users" ("email");

CREATE INDEX ON "subscriptions" ("user_id");

ALTER TABLE "subscriptions" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");