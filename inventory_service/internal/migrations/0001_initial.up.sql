create table if not exists transactions(
                                           id           bigserial primary key,
                                           created_at   timestamp with time zone not null default NOW(),
                                           updated_at   timestamp with time zone not null default NOW(),
                                           deleted_at   timestamp with time zone,
                                           "product_id" varchar(255)             not null,
                                           operation    integer                  not null,
                                           amount       integer                  not null,
                                           cart_id      varchar(255)
);
create index if not exists idx_transactions_deleted_at on transactions (deleted_at);
create index if not exists idx_transactions_product_id on transactions (product_id);
