create table carts
(
    id         bigserial
        primary key,
    created_at timestamp with time zone default NOW(),
    updated_at timestamp with time zone default NOW(),
    deleted_at timestamp with time zone,
    user_id    text   not null,
    status     bigint not null default 1,
    total      bigint not null default 0
        constraint chk_carts_total
            check (total >= 0)
);

create index idx_carts_user_id
    on carts (user_id);

create index idx_carts_deleted_at
    on carts (deleted_at);

create index idx_carts_status
    on carts (status);

create unique index carts_user_id_uindex
    on carts (user_id)
    where (status = 1);
