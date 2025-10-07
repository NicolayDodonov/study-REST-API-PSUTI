BEGIN;

create table if not exists users
(
    id         uuid default gen_random_uuid() not null
    constraint users_pk
    primary key,
    first_name text,
    last_name  text,
    user_type  text,
    login      text,
    password   text,
    height     integer,
    weight     integer,
    age        integer,
    sex        text
)

COMMIT;