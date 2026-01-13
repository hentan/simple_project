-- +goose Up
SELECT 'up SQL query';
create table expenses
(
    id       serial,
    date     date,
    gift_for varchar(500),
    pupil_id  int,
    summ     int
);
create table payments
(
    id serial,
    pupil_id int,
    summ int
);

create table pupils(
    id serial,
    name varchar(100),
    surname varchar(100),
    parent_name varchar(100),
    parent_phone varchar(50)
);
-- +goose Down
SELECT 'down SQL query';
drop table expenses;
drop table payments;
drop table pupils;