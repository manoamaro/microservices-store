create table carts
(
    id         bigserial
        primary key,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    user_id    text   not null,
    status     bigint not null,
    total      bigint not null
        constraint chk_carts_total
            check (total >= 0)
);

create index idx_carts_user_id
    on carts (user_id);

create index idx_carts_deleted_at
    on carts (deleted_at);

create index idx_carts_status
    on carts (status);

