create table cart_items
(
    id         bigserial
        primary key,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    cart_id    bigint not null
        constraint fk_carts_items
            references carts
            on update cascade on delete cascade,
    product_id text   not null,
    quantity   bigint not null
        constraint chk_cart_items_quantity
            check (quantity > 0)
);

create index idx_cart_items_product_id
    on cart_items (product_id);

create index idx_cart_items_cart_id
    on cart_items (cart_id);

create index idx_cart_items_deleted_at
    on cart_items (deleted_at);