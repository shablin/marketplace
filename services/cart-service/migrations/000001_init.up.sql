create extension if not exists pgcrypto;

create table if not exists carts (
    id uuid primary key default gen_random_uuid(),
    user_id uuid not null,
    status text not null default 'active'
        check (status in ('active', 'checked_out', 'abandoned')),
    currency char(3) not null default 'rub',
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now()
);

create table if not exists cart_items (
    id uuid primary key default gen_random_uuid(),
    cart_id uuid not null references carts(id) on delete cascade,
    product_id uuid not null,
    quantity int not null check (quantity > 0),
    uint_price numeric(12, 2) not null check (uint_price >= 0),
    added_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    constraint cart_items_cart_product_uniq unique (cart_id, product_id)
);

create index if not exists
    carts_user_status_idx on carts (user_id, status);

create unique index if not exists
    carts_user_active_uniq_idx on carts (user_id) where status = 'active';

create index if not exists
    carts_updated_at_idx on carts (updated_at desc);

create index if not exists
    carts_items_cart_id_idx on carts_items (cart_id);

create index if not exists
    cart_items_product_id_idx on cart_items (product_id);