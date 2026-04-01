-- +goose Up

-- =========================
-- Base tables
-- =========================
create table expenses
(
    id        serial primary key,
    date      date,
    gift_for  varchar(500),
    pupil_id  int,
    summ      int
);

create table payments
(
    id        serial primary key,
    pupil_id  int,
    summ      int,
    date      date default now(),
    purpose   varchar(500)
);

create table pupils
(
    id           serial primary key,
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
    op          text        not null,

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

-- +goose StatementBegin
create or replace function trg_expenses_archive()
    returns trigger
    language plpgsql
as $$
begin
    insert into expenses_archive (op, id, date, gift_for, pupil_id, summ)
    values (tg_op, old.id, old.date, old.gift_for, old.pupil_id, old.summ);

    if tg_op = 'UPDATE' then
        return new;
    else
        return old; -- DELETE
    end if;
end;
$$;
-- +goose StatementEnd

-- +goose StatementBegin
create or replace function trg_payments_archive()
    returns trigger
    language plpgsql
as $$
begin
    insert into payments_archive (op, id, pupil_id, summ)
    values (tg_op, old.id, old.pupil_id, old.summ);

    if tg_op = 'UPDATE' then
        return new;
    else
        return old;
    end if;
end;
$$;
-- +goose StatementEnd

-- +goose StatementBegin
create or replace function trg_pupils_archive()
    returns trigger
    language plpgsql
as $$
begin
    insert into pupils_archive (op, id, name, surname, parent_name, parent_phone)
    values (tg_op, old.id, old.name, old.surname, old.parent_name, old.parent_phone);

    if tg_op = 'UPDATE' then
        return new;
    else
        return old;
    end if;
end;
$$;
-- +goose StatementEnd

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
