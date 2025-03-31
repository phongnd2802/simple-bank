-- +goose Up
-- +goose StatementBegin
CREATE TABLE "verify_emails" (
    "id" bigserial PRIMARY KEY,
    "username" varchar NOT NULL,
    "email" varchar NOT NULL,
    "secret_code" varchar NOT NULL,
    "is_used" boolean NOT NULL DEFAULT false,
    "created_at" timestamptz NOT NULL DEFAULT (now())
);


ALTER TABLE "verify_emails" ADD FOREIGN KEY ("username") REFERENCES "users" ("username");


ALTER TABLE "users" ADD COLUMN "is_email_verified" boolean NOT NULL DEFAULT false;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS "verify_emails" CASCADE;

ALTER TABLE "users" DROP COLUMN "is_email_verified";
-- +goose StatementEnd
