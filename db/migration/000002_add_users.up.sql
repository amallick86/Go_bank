CREATE TABLE "users" (
  "username" varchar PRIMARY KEY,
  "hashed_password" varchar NOT NULL,
  "full_name" varchar NOT NULL,
  "email" varchar UNIQUE NOT NULL,
  "ctitizianship" varchar UNIQUE NOT NULL,
  "password_changed_at" timestamp NOT NULL DEFAULT '0001-01-01 00:00:00+00',
  "created_at" timestamp NOT NULL DEFAULT 'now()'
);

ALTER TABLE "accounts" ADD FOREIGN KEY ("owner") REFERENCES "users" ("username");

ALTER TABLE "accounts" ADD FOREIGN KEY ("ctitizianship") REFERENCES "users" ("ctitizianship");

--CREATE UNIQUE INDEX ON "accounts" ("ctitizianship", "currency");
ALTER TABLE "accounts" ADD CONSTRAINT "citizen_currency_key" UNIQUE ("ctitizianship", "currency");