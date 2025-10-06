BEGIN;

create table if not exists user_data
(
    id              uuid not null
    constraint user_data_pk
    primary key,
    count_read_book integer,
    reader_score    integer,
    favorite_book   text[][]
);

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
    data       uuid
    constraint users_user_data_id_fk
    references user_data
    );

COMMIT;