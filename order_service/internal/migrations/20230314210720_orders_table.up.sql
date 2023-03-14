create table orders
(
    id                             bigserial
        primary key,
    created_at                     timestamp with time zone,
    updated_at                     timestamp with time zone,
    deleted_at                     timestamp with time zone,
    user_id                        text   not null,
    status                         bigint not null,
    cart_id                        bigint not null
        constraint fk_orders_cart
            references carts,
    shipping_address_first_name    text   not null,
    shipping_address_last_name     text   not null,
    shipping_address_address_line1 text   not null,
    shipping_address_address_line2 text,
    shipping_address_zip_code      text   not null,
    shipping_address_region        text,
    shipping_address_state         text   not null,
    shipping_address_country       text   not null,
    invoice_address_first_name     text   not null,
    invoice_address_last_name      text   not null,
    invoice_address_address_line1  text   not null,
    invoice_address_address_line2  text,
    invoice_address_zip_code       text   not null,
    invoice_address_region         text,
    invoice_address_state          text   not null,
    invoice_address_country        text   not null,
    total                          bigint
);

create index idx_orders_cart_id
    on orders (cart_id);

create index idx_orders_status
    on orders (status);

create index idx_orders_user_id
    on orders (user_id);

create index idx_orders_deleted_at
    on orders (deleted_at);
