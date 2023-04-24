create table reservations
(
    product_id varchar   not null primary key,
    amount     integer   not null,
    created_at timestamp not null default CURRENT_TIMESTAMP
);
