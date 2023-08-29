CREATE DATABASE avito;
\c avito;
CREATE SCHEMA public;
CREATE TABLE public.users_segments (
    USER_ID BIGINT,
    SEGMENT_ID INT
);
CREATE TABLE public.id_name_segments (
    SEGMENT_ID INT PRIMARY KEY,
    SEGMENT_NAME TEXT
);
