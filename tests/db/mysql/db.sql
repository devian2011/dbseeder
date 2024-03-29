CREATE TABLE IF NOT EXISTS areas
(
    id   int not null primary key auto_increment,
    name varchar(255)
) ENGINE = InnoDB;

CREATE TABLE IF NOT EXISTS users
(
    id          int          not null primary key auto_increment,
    username    varchar(255) not null,
    firstname   varchar(255) not null,
    lastname    varchar(255) not null,
    fullname    varchar(255) not null,
    password    varchar(255) not null,
    last_online timestamp default NOW(),
    area_id     int          not null,
    CONSTRAINT user_areas_fk FOREIGN KEY (area_id) REFERENCES areas (id)
) ENGINE = InnoDB;

CREATE TABLE IF NOT EXISTS roles
(
    id   int          not null primary key auto_increment,
    name varchar(255) not null
) ENGINE = InnoDB;

CREATE TABLE IF NOT EXISTS user_in_roles
(
    user_id int not null,
    role_id int not null,
    CONSTRAINT user_role_user_id_fk FOREIGN KEY (user_id) REFERENCES users (id),
    CONSTRAINT user_role_role_id_fk FOREIGN KEY (role_id) REFERENCES roles (id)
) ENGINE = InnoDB;

CREATE TABLE IF NOT EXISTS info
(
    id      int           not null primary key auto_increment,
    phone   varchar(255)  not null,
    address varchar(1024) not null,
    user_id int           not null,
    CONSTRAINT user_info_fk FOREIGN KEY (user_id) REFERENCES users (id)
) ENGINE = InnoDB;
