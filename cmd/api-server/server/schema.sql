CREATE TABLE "public"."users" (
    "id" serial NOT NULL,
    "name" character varying(50) NOT NULL,
    "created_at" timestamp DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT users_pkey PRIMARY KEY ("id")
);
