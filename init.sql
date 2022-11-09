CREATE TABLE IF NOT EXISTS "article" (
  "id" BIGSERIAL PRIMARY KEY,
  "title" TEXT NOT NULL,
  "author" TEXT NOT NULL,
  "category" TEXT NOT NULL,
  "url" TEXT NOT NULL,
  "created_at" TIMESTAMP NOT NULL
);