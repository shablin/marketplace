create extension if not exists pgcrypto;

create table if not exists payments (
    id uuid primary key default gen_random_uuid(),
    order_id uuid not null,
    provider text not null,
    status text not null default 'created'
        check (status in ('created', 'processing', 'paid', 'failed', 'refunded')),
    amount numeric(12, 2) not null check (amount >= 0),
    currency char(3) not null default 'RUB',
    external_payment_id text,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    paid_at timestamptz,
    constraint payments_order_provider_uniq unique (order_id, provider)
);

create table if not exists payment_transactions (
    id uuid primary key default gen_random_uuid(),
    payment_id uuid not null references payments(id) on delete cascade,
    transaction_type text not null
        check (transaction_type in ('authorization', 'capture', 'refund', 'void')),
    status text not null
        check (status in ('created', 'processing', 'succeeded', 'failed')),
    amount numeric(12, 2) not null check (amount >= 0),
    external_transaction_id text,
    gateway_response jsonb,
    created_at timestamptz not null default now(),
    constraint payment_transactions_external_uniq unique (external_transaction_id)
);

create index if not exists
    payments_order_status_idx
on payments (order_id, status);

create index if not exists
    payments_status_created_idx
on payments (status, created_at desc);

create index if not exists
    payment_transactions_payment_created_idx
on payment_transactions (payment_id, created_at desc);

create index if not exists
    payment_transactions_status_idx
on payment_transactions (status);