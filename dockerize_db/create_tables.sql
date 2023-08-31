CREATE DATABASE avito;

\c avito;

CREATE TABLE public.users_segments (
                                      user_id       BIGINT,
                                      segment_id    INT,
                                      endtime       TIMESTAMP
);

CREATE TABLE public.id_name_segments (
                                      segment_id    SERIAL PRIMARY KEY,
                                      segment_name  TEXT NOT NULL
);

CREATE TABLE public.user_segment_history (
                                      id            SERIAL PRIMARY KEY,
                                      user_id       BIGINT,
                                      segment_name  TEXT,
                                      operation     VARCHAR(10),
                                      timestamp     TIMESTAMP
);

CREATE TABLE public.users (
                                      user_id       BIGINT PRIMARY KEY
)