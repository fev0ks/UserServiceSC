-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE table if not exists "user" (
  "id" bigserial PRIMARY KEY,
  "name" varchar NOT null,
  "age" integer NOT NULL,
  "created_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" timestamp
);
CREATE table if not exists "item" (
  "id" bigserial PRIMARY KEY,
  "name" varchar NOT NULL,
  "created_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" timestamp
);
CREATE table if not exists "type" (
  "id" integer PRIMARY KEY,
  "type" varchar NOT NULL
);
CREATE TABLE if not exists "user_type" (
  "user_id" integer UNIQUE NOT NULL REFERENCES "user" ("id") ON delete cascade,
  "type_id" integer NOT NULL REFERENCES "type" ("id")
);
-- actually this is not required table, items may belong to only 1 user...
CREATE TABLE if not exists "user_item" (
  "user_id" integer NOT NULL REFERENCES "user" ("id") ON delete cascade,
  "item_id" integer NOT NULL REFERENCES item ("id") ON delete cascade
);
INSERT into "type"("id", "type") values
        (0, 'INVALID_USER_TYPE'),
        (1, 'EMPLOYEE_USER_TYPE'),
        (2, 'CUSTOMER_USER_TYPE')
    ON CONFLICT DO NOTHING;

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE IF EXISTS "user_item";
DROP TABLE IF EXISTS "user_type";
DROP TABLE IF EXISTS "type";
DROP TABLE IF EXISTS "item";
DROP TABLE IF EXISTS "user";
