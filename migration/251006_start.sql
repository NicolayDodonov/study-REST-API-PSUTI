BEGIN;

create table if not exists users
(
    id         uuid    default gen_random_uuid() not null
    constraint users_pk
    primary key,
    first_name text,
    last_name  text,
    user_type  text,
    login      text                              not null,
    password   text                              not null,
    height     integer default 0                 not null,
    weight     integer default 0                 not null,
    age        integer default 0                 not null,
    sex        text    default 'man'::text       not null
    );


COMMIT;