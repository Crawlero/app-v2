CREATE TABLE "public"."crawlers" ( 
  "id" UUID NOT NULL DEFAULT gen_random_uuid(),
  "name" VARCHAR(255) NOT NULL,
  "description" TEXT NULL,
  "status" VARCHAR(10) NOT NULL,
  "schema" JSON NULL,
  "user_id" UUID NULL,
  "url" VARCHAR(255) NOT NULL,
  CONSTRAINT "crawlers_pkey" PRIMARY KEY ("id"),
  CONSTRAINT "crawlers_userid_fkey" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON DELETE NO ACTION ON UPDATE NO ACTION
);
