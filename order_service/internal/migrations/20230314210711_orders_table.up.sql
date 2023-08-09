create table public.orders
(
    id                             bigserial
        primary key,
    created_at                     timestamp with time zone,
    updated_at                     timestamp with time zone,
    deleted_at                     timestamp with time zone,
    user_id                        text,
    status                         bigint,
    shipping_address_first_name    text,
    shipping_address_last_name     text,
    shipping_address_address_line1 text,
    shipping_address_address_line2 text,
    shipping_address_zip_code      text,
    shipping_address_region        text,
    shipping_address_state         text,
    shipping_address_country       text,
    invoice_address_first_name     text,
    invoice_address_last_name      text,
    invoice_address_address_line1  text,
    invoice_address_address_line2  text,
    invoice_address_zip_code       text,
    invoice_address_region         text,
    invoice_address_state          text,
    invoice_address_country        text,
    total                          bigint
);

alter table public.orders
    owner to postgres;

create index idx_orders_status
    on public.orders (status);

create index idx_orders_user_id
    on public.orders (user_id);

create index idx_orders_deleted_at
    on public.orders (deleted_at);

