CREATE TABLE "users" (
    "id" serial primary key,
    "name" varchar(255) not null,
    "username" varchar(255) not null unique,
    "password_hash" varchar(255) not null
);

CREATE TABLE "todo_lists"
(
    "id" serial primary key,
    "title" varchar(255) not null,
    "description" text
);

CREATE TABLE "users_lists"
(
    "id" serial primary key,
    "user_id" int references "users"(id) on delete cascade on update cascade not null,
    "list_id" int references "todo_lists"(id) on delete cascade on update cascade not null
);

CREATE TABLE "todo_items"
(
    "id" serial primary key,
    "title" varchar(255) not null,
    "description" text,
    "done" boolean not null default false
);

CREATE TABLE "lists_items"
(
    "id" serial primary key,
    "list_id" int references "todo_lists"(id) on delete cascade on update cascade not null,
    "item_id" int references "todo_items"(id) on delete cascade on update cascade not null
);
