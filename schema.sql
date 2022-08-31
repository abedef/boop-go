CREATE TABLE boops (
    id integer GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    text text NOT NULL DEFAULT 'null'::text,
    created timestamp with time zone NOT NULL DEFAULT now()
);