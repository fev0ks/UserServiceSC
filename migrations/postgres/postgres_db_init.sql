-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE table if not exists "user" (
  "id" bigserial PRIMARY KEY,
  "name" varchar NOT null,
  "age" int8 NOT NULL,
  "created_at" timestamp NOT NULL DEFAULT 'now()',
  "updated_at" timestamp NOT NULL DEFAULT 'now()'
);
CREATE table if not exists "item" (
  "id" bigserial PRIMARY KEY,
  "name" varchar NOT NULL,
  "created_at" timestamp NOT NULL DEFAULT 'now()',
  "updated_at" timestamp NOT NULL DEFAULT 'now()'
);
CREATE table if not exists "type" (
  "id" bigserial PRIMARY KEY,
  "type" varchar NOT NULL
);
CREATE TABLE if not exists "user_type" (
  "user_id" bigint NOT NULL REFERENCES "user" ("id"),
  "type_id" bigint NOT NULL REFERENCES "type" ("id")
);
CREATE TABLE if not exists "user_item" (
  "user_id" bigint NOT NULL REFERENCES "user" ("id") ON delete cascade,
  "item_id" bigint NOT NULL REFERENCES item ("id")
);

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE IF EXISTS "user_item";
DROP TABLE IF EXISTS "user_type";
DROP TABLE IF EXISTS "type";
DROP TABLE IF EXISTS "item";
DROP TABLE IF EXISTS "user";
