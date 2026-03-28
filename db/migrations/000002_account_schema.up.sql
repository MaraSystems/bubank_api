CREATE TABLE "accounts" (
  "id" bigserial PRIMARY KEY,
  "owner" varchar NOT NULL,
  "balance" bigint NOT NULL,
  "currency" varchar NOT NULL,
  "created_at" timestamptz DEFAULT (now())
);

CREATE INDEX ON "accounts" ("owner");

COMMENT ON COLUMN "accounts"."balance" IS 'can be positive or negative';

ALTER TABLE "accounts" ADD FOREIGN KEY ("owner") REFERENCES "users" ("username") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "accounts" ADD CONSTRAINT "owner_currency_key" UNIQUE ("owner", "currency");
