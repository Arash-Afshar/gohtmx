CREATE TABLE IF NOT EXISTS "samples" (
	"id"   TEXT NOT NULL UNIQUE,
	"name" TEXT,
	PRIMARY KEY("id")
);

CREATE TABLE IF NOT EXISTS "posts" (
	"id"   TEXT NOT NULL UNIQUE,
	"title" TEXT,
	PRIMARY KEY("id")
);

