create table users
(
    id       serial primary key,
    username varchar unique not null,
    password bytea          not null
);

create table resources
(
    id      serial primary key,
    user_id int,
    alias   varchar,
    type    int not null,
    data    bytea,
    meta    bytea,

    CONSTRAINT fk_users FOREIGN KEY (user_id) REFERENCES users (id) on delete cascade
);
-- create unique index on resources (user_id, alias);
---- create above / drop below ----
DROP TABLE IF EXISTS "resources";
DROP TABLE IF EXISTS "users";
