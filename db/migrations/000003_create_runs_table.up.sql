CREATE TABLE "public"."runs" ( 
  "id" UUID NOT NULL DEFAULT gen_random_uuid() ,
  "status" VARCHAR(255) NOT NULL,
  "started_at" TIMESTAMP NOT NULL,
  "updated_at" TIMESTAMP NOT NULL,
  "crawler_id" UUID NOT NULL,

  CONSTRAINT "runs_pkey" PRIMARY KEY ("id"),
  CONSTRAINT "runs_crawler_id_fkey" FOREIGN KEY ("crawler_id") REFERENCES "public"."crawlers"("id")
);
