CREATE TABLE IF NOT EXISTS "samples" (
	"id"   INTEGER NOT NULL UNIQUE,
	"name" TEXT,
	PRIMARY KEY("id" AUTOINCREMENT)
);

CREATE TABLE IF NOT EXISTS "posts" (
	"id"   TEXT NOT NULL UNIQUE,
	"title" TEXT
);
