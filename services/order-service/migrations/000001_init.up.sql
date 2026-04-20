create extension if not exists pgcrypto;

create table if not exists orders (
    id uuid primary key default gen_random_uuid(),
    user_id uuid not null,
    cart_id uuid,
    order_number text not null,
    status text not null default 'created'
        check (status in ('created', 'awaiting_payment', 'paid', 'cancelled', 'completed')),
    currency char(3) not null default 'RUB',
    subtotal numeric(12, 2) not null default 0
        check (subtotal >= 0),
    shipping_amount numeric(12, 2) not null default 0
        check (shipping_amount >= 0),
    discount_amount numeric(12, 2) not null default 0
        check (discount_amount >= 0),
    total_amount numeric(12, 2) not null default 0
        check (total_amount >= 0),
    placed_at timestamptz,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    constraint orders_order_number_uniq unique (order_number)
);

create table if not exists order_items (
    id uuid primary key default gen_random_uuid(),
    order_id uuid not null references orders(id)
        on delete cascade,
    product_id uuid not null,
    product_name text not null,
    quantity int not null check (quantity > 0),
    uint_price numeric(12, 2) not null check (uint_price >= 0),
    line_total numeric(12, 2) not null check (line_total >= 0),
    created_at timestamptz not null default now(),
    constraint order_items_order_product_uniq unique (order_id, product_id)
);

create index if not exists
    orders_user_status_idx
on orders (user_id, status);

create index if not exists
    orders_status_created_idx
on orders (status, created_at desc);

create index if not exists
    orders_placed_at_idx
on orders (placed_at desc);

create index if not exists
    order_item_order_id_idx
on order_items (order_id);

create index if not exists
    order_items_product_id_idx
on order_items (product_id);