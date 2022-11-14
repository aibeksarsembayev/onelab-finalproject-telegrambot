CREATE TABLE IF NOT EXISTS "article" (
  "id" BIGSERIAL PRIMARY KEY,
  "article_id" BIGINT NOT NULL UNIQUE,
  "title" TEXT NOT NULL,
  "author" TEXT NOT NULL,
  "category" TEXT NOT NULL,
  "url" TEXT NOT NULL,
  "created_at" TIMESTAMP NOT NULL
);

CREATE TABLE IF NOT EXISTS "article_api" (
  "id" BIGSERIAL PRIMARY KEY,
  "article_id" BIGINT NOT NULL UNIQUE,
  "user_id" BIGINT NOT NULL,
  "category_id" BIGINT NOT NULL,
  "category" TEXT NOT NULL,
  "title" TEXT NOT NULL,
  "created_at" TIMESTAMP NOT NULL,
  "updated_at" TIMESTAMP NOT NULL,
  "author_first_name" TEXT NOT NULL,
  "author_last_name" TEXT NOT NULL,
  "url" TEXT NOT NULL
);