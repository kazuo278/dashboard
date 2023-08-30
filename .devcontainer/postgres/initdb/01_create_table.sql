DROP TABLE IF EXISTS "repositories" CASCADE;
CREATE TABLE IF NOT EXISTS "repositories" (
  "repository_id" varchar(100) NOT NULL PRIMARY KEY,
  "repository_name" varchar(100) NOT NULL
);

DROP TABLE IF EXISTS "jobs" CASCADE;
CREATE TABLE IF NOT EXISTS "jobs" (
  "job_id" varchar(100) NOT NULL,
  "repository_id" varchar(100) NOT NULL,
  "run_id" varchar(100) NOT NULL,
  "workflow_ref"  varchar(100) NOT NULL,
  "job_name" varchar(100) NOT NULL,
  "run_attempt" varchar(10) NOT NULL,
  "status" varchar(15) NOT NULL,
  "started_at" timestamp WITH TIME ZONE,
  "finished_at" timestamp WITH TIME ZONE,
  PRIMARY KEY("job_id")
);

DROP TABLE IF EXISTS "job_details" CASCADE;
CREATE TABLE IF NOT EXISTS "job_details" (
  "job_detail_id" SERIAL NOT NULL,
  "job_id" varchar(100) NOT NULL,
  "type" varchar(20) NOT NULL,
  "using_path" varchar(100) NOT NULL,
  "using_ref"  varchar(100) NOT NULL,
  PRIMARY KEY("job_detail_id")
);