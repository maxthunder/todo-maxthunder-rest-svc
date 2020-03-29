create schema task collate utf8_general_ci;

create table task
(
    task_id int auto_increment,
    description varchar(255) not null,
    timestamp varchar(255) not null,
    is_completed boolean default false not null,
    constraint task_pk
        primary key (task_id)
);

