CREATE TABLE IF NOT EXISTS areas (
    id bigserial not null primary key,
    name varchar(255)
);

CREATE TABLE IF NOT EXISTS users (
    id bigserial not null primary key,
    username varchar(255) not null,
    firstname varchar(255) not null,
    lastname varchar(255) not null,
    fullname varchar(255) not null,
    password varchar(255) not null,
    last_online timestamp default NOW(),
    area_id bigint not null,
    CONSTRAINT fk_area_id FOREIGN KEY (area_id) REFERENCES areas(id)
);

CREATE TABLE IF NOT EXISTS roles (
    id bigserial not null primary key,
    name varchar(255) not null
);

CREATE TABLE IF NOT EXISTS user_in_roles (
    user_id bigint not null,
    role_id bigint not null,
    CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES users(id),
    CONSTRAINT fk_role_id FOREIGN KEY (role_id) REFERENCES roles(id)
);

CREATE TABLE IF NOT EXISTS info (
    id bigserial not null primary key,
    phone varchar(255) not null,
    address varchar(1024) not null,
    user_id bigint not null,
    CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES users(id)
);

