create table order_items
(
    id         bigserial
        primary key,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    order_id   bigint
        constraint fk_orders_items
            references orders,
    product_id text,
    quantity   bigint
);

alter table order_items
    owner to postgres;

create index idx_order_items_product_id
    on order_items (product_id);

create index idx_order_items_order_id
    on order_items (order_id);

create index idx_order_items_deleted_at
    on order_items (deleted_at);

