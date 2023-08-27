CREATE DATABASE avito;
\c avito;

CREATE SCHEMA avito.public;

CREATE TABLE users_segments (
                                       USER_ID INT PRIMARY KEY,
                                       SEGMENT_ID INT
);

CREATE TABLE id_name_segments (
                                         SEGMENT_ID INT PRIMARY KEY,
                                         SEGMENT_NAME TEXT
);kl
