CREATE TABLE users
(
    id serial primary KEY,
    name varchar(255)           not null,
    username varchar(255)       not null unique,
    password_hash varchar(255)  not null
);

CREATE TABLE todo_lists
(
    id serial primary KEY,
    title varchar(255) not null,
    description varchar(255) 
);

CREATE TABLE users_lists
(
    id serial not null unique,
    user_id integer not null,
    list_id integer not null,
    FOREIGN KEY(user_id) REFERENCES users(id) on delete cascade,
    FOREIGN KEY(list_id) REFERENCES todo_lists(id) on delete cascade
);


CREATE TABLE todo_items
(
    id serial primary KEY,
    title varchar(255) not null,
    description varchar(255),
    done boolean not null default false
);


CREATE TABLE list_items
(
    id serial not null unique,
    item_id integer not null,
    list_id integer not null,
    FOREIGN KEY(item_id) REFERENCES todo_items(id) on delete cascade,
    FOREIGN KEY(list_id) REFERENCES todo_lists(id) on delete cascade
);
