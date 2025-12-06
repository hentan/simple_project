-- +goose Up
-- +goose StatementBegin
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

create table pupils_info(
    id serial,
    name varchar(100),
    surname varchar(100),
    parent_name varchar(100),
    parent_phone varchar(50)
);
-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
drop table expenses;
drop table payments;
drop table pupils;