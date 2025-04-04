
CREATE SCHEMA IF NOT EXISTS "silly";


CREATE TABLE IF NOT EXISTS "silly"."silly_type"
(
	id bigint,
	name text UNIQUE NOT NULL,
	description text NULL,
	Primary Key (id)
); 
INSERT INTO "silly"."silly_type" VALUES (1,'Funny','');
INSERT INTO "silly"."silly_type" VALUES (2,'Strange','');
INSERT INTO "silly"."silly_type" VALUES (3,'Dangerous','');


CREATE TABLE IF NOT EXISTS "silly"."actor"
(
	"id" bigserial,
	"name" text NOT NULL,
	PRIMARY KEY ("id")
); 


CREATE TABLE IF NOT EXISTS "silly"."movie"
(
	"id" bigserial,
	"name" text NOT NULL,
	"silly" bigint NOT NULL,
	PRIMARY KEY ("id"),
	FOREIGN KEY ("silly") REFERENCES "silly"."silly_type" ("id") 
); 


CREATE TABLE IF NOT EXISTS "silly"."movie_actor"
(
	"id" bigserial,
	"actor" bigint NOT NULL,
	"movie" bigint NOT NULL,
	PRIMARY KEY ("id"),
	FOREIGN KEY ("actor") REFERENCES "silly"."actor" ("id"), 
	FOREIGN KEY ("movie") REFERENCES "silly"."movie" ("id") 
); 


