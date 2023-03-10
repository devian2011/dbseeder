CREATE TABLE IF NOT EXISTS areas (
    id int not null primary key auto_increment,
    name varchar(255)
);

CREATE TABLE IF NOT EXISTS users (
    id int not null primary key auto_increment,
    username varchar(255) not null,
    firstname varchar(255) not null,
    lastname varchar(255) not null,
    fullname varchar(255) not null,
    password varchar(255) not null,
    last_online timestamp default NOW(),
    area_id bigint not null
);

CREATE TABLE IF NOT EXISTS roles (
    id int not null primary key auto_increment,
    name varchar(255) not null
);

CREATE TABLE IF NOT EXISTS user_in_roles (
    user_id bigint not null,
    role_id bigint not null
);

CREATE TABLE IF NOT EXISTS info (
    id int not null primary key auto_increment,
    phone varchar(255) not null,
    address varchar(1024) not null,
    user_id bigint not null
);
