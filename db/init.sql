CREATE DATABASE tripdb
    WITH
    ENCODING = 'UTF8'
    LC_COLLATE = 'en_US.utf8'
    LC_CTYPE = 'en_US.utf8';

\c tripdb;

create table rates
(
    id         bigserial primary key,
    from_point point,
    to_point   point,
    status     smallint,
    distance   int8,
    duration   int8,
    created_at timestamp,
    updated_at timestamp
);