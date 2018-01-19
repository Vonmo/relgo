
-- +migrate Up
DROP TABLE IF EXISTS public.counters;

CREATE TABLE public.counters
(
  name character varying(50) NOT NULL,
  value integer NOT NULL DEFAULT 0,
  updated timestamp with time zone DEFAULT now(),
  CONSTRAINT counters_pk PRIMARY KEY (name)
);

-- +migrate Down
DROP TABLE IF EXISTS public.counters;
