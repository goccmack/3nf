
CREATE SCHEMA IF NOT EXISTS "binary";


CREATE TABLE IF NOT EXISTS "binary"."table1"
(
	"id" bigserial,
	"blob" bytea NULL,
	PRIMARY KEY ("id")
); 


