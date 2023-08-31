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
CREATE TABLE user_segment_history (
                                      id         serial primary key,
                                      user_id    integer,
                                      segment_id integer,
                                      operation  varchar(10),
                                      timestamp  timestamp
);

