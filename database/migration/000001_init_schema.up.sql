CREATE TABLE "accounts" (
  "id" BIGSERIAL PRIMARY KEY,
  "owner" bigint NOT NULL,
  "balance" bigint NOT NULL,
  "currency" varchar NOT NULL,
  "created_at" timestamp NOT NULL DEFAULT (now())
);

CREATE TABLE "users" (
  "id" BIGSERIAL PRIMARY KEY,
  "username" varchar UNIQUE NOT NULL,
  "created_at" timestamp NOT NULL DEFAULT (now())
);

CREATE TABLE "entries" (
  "id" BIGSERIAL PRIMARY KEY,
  "amount" bigint NOT NULL,
  "account_id" bigint NOT NULL,
  "created_at" timestamp NOT NULL DEFAULT (now())
);

CREATE TABLE "transfars" (
  "id" BIGSERIAL PRIMARY KEY,
  "amount" bigint NOT NULL,
  "status" varchar NOT NULL DEFAULT 'pending',
  "from_account_id" bigint NOT NULL,
  "to_account_id" bigint NOT NULL,
  "created_at" timestamp NOT NULL DEFAULT (now())
);

CREATE INDEX ON "accounts" ("owner");

CREATE INDEX ON "entries" ("account_id");

CREATE INDEX ON "transfars" ("from_account_id");

CREATE INDEX ON "transfars" ("to_account_id");

CREATE INDEX ON "transfars" ("from_account_id", "to_account_id");

COMMENT ON COLUMN "transfars"."amount" IS 'Must be positive';

ALTER TABLE "accounts" ADD FOREIGN KEY ("owner") REFERENCES "users" ("id");

ALTER TABLE "entries" ADD FOREIGN KEY ("account_id") REFERENCES "accounts" ("id");

ALTER TABLE "transfars" ADD FOREIGN KEY ("from_account_id") REFERENCES "accounts" ("id");

ALTER TABLE "transfars" ADD FOREIGN KEY ("to_account_id") REFERENCES "accounts" ("id");
