create table if not exists flags (
    id bigserial primary key,
    created_at timestamp with time zone not null default NOW(),
    updated_at timestamp with time zone not null default NOW(),
    deleted_at timestamp with time zone,
    "name" varchar(255) not null
);
create index if not exists idx_flags_deleted_at on flags (deleted_at);

create table if not exists domains (
    id bigserial primary key,
    created_at timestamp with time zone not null default NOW(),
    updated_at timestamp with time zone not null default NOW(),
    deleted_at timestamp with time zone,
    "domain" varchar(255) not null
);
create index if not exists idx_domains_deleted_at on domains (deleted_at);

create table if not exists auths (
    id bigserial primary key,
    created_at timestamp with time zone not null default NOW(),
    updated_at timestamp with time zone not null default NOW(),
    deleted_at timestamp with time zone,
    email varchar(255) not null constraint idx_auths_email unique,
    password varchar(255) not null,
    salt varchar(255) default '' not null
);
create index if not exists idx_auths_deleted_at on auths (deleted_at);

create table if not exists auths_domains (
    auth_id bigint constraint fk_audiences_auth references auths,
    domain_id bigint constraint fk_audiences_domain references domains,
    primary key (auth_id, domain_id)
);

create table if not exists auths_flags (
    auth_id bigint not null constraint fk_auths_flags_auth references auths,
    flag_id bigint not null constraint fk_auths_flags_flag references flags,
    primary key (auth_id, flag_id)
);