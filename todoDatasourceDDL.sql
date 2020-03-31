create schema todo collate utf8_general_ci;

create table task
(
    taskId int auto_increment,
    description varchar(255) not null,
    timestamp varchar(255) not null,
    isCompleted boolean default false not null,
    constraint task_pk
        primary key (taskId)
);

