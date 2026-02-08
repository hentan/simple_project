-- +goose Up
SELECT 'up SQL query';

-- =========================
-- Base tables
-- =========================
create table expenses
(
    id        serial,
    date      date,
    gift_for  varchar(500),
    pupil_id  int,
    summ      int
);

create table payments
(
    id        serial,
    pupil_id  int,
    summ      int
);

create table pupils
(
    id           serial,
    name         varchar(100),
    surname      varchar(100),
    parent_name  varchar(100),
    parent_phone varchar(50)
);

-- =========================
-- Archive tables
-- =========================
create table expenses_archive
(
    archived_at timestamptz not null default now(),
    op          text        not null, -- 'UPDATE' | 'DELETE'

    id        int,
    date      date,
    gift_for  varchar(500),
    pupil_id  int,
    summ      int
);

create table payments_archive
(
    archived_at timestamptz not null default now(),
    op          text        not null,

    id        int,
    pupil_id  int,
    summ      int
);

create table pupils_archive
(
    archived_at  timestamptz not null default now(),
    op           text        not null,

    id           int,
    name         varchar(100),
    surname      varchar(100),
    parent_name  varchar(100),
    parent_phone varchar(50)
);

-- =========================
-- Trigger functions
-- =========================
create or replace function trg_expenses_archive()
returns trigger
language plpgsql
as $$
begin
insert into expenses_archive (op, id, date, gift_for, pupil_id, summ)
values (tg_op, old.id, old.date, old.gift_for, old.pupil_id, old.summ);
return old;
end;
$$;

create or replace function trg_payments_archive()
returns trigger
language plpgsql
as $$
begin
insert into payments_archive (op, id, pupil_id, summ)
values (tg_op, old.id, old.pupil_id, old.summ);
return old;
end;
$$;

create or replace function trg_pupils_archive()
returns trigger
language plpgsql
as $$
begin
insert into pupils_archive (op, id, name, surname, parent_name, parent_phone)
values (tg_op, old.id, old.name, old.surname, old.parent_name, old.parent_phone);
return old;
end;
$$;

-- =========================
-- Triggers
-- =========================
create trigger expenses_archive_bu_bd
    before update or delete on expenses
for each row
execute function trg_expenses_archive();

create trigger payments_archive_bu_bd
    before update or delete on payments
for each row
execute function trg_payments_archive();

create trigger pupils_archive_bu_bd
    before update or delete on pupils
for each row
execute function trg_pupils_archive();


-- +goose Down
SELECT 'down SQL query';

-- Drop triggers first, then functions, then tables
drop trigger if exists expenses_archive_bu_bd on expenses;
drop trigger if exists payments_archive_bu_bd on payments;
drop trigger if exists pupils_archive_bu_bd on pupils;

drop function if exists trg_expenses_archive();
drop function if exists trg_payments_archive();
drop function if exists trg_pupils_archive();

drop table if exists expenses_archive;
drop table if exists payments_archive;
drop table if exists pupils_archive;

drop table if exists expenses;
drop table if exists payments;
drop table if exists pupils;
