create extension if not exists pgcrypto;

create table if not exists users (
    id uuid primary key default gen_random_uuid(),
    email text not null,
    username text not null,
    password_hash text not null,
    first_name text,
    last_name text,
    phone text,
    status text not null default 'pending'
        check (status in ('pending', 'active', 'blocked', 'deleted')),
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    deleted_at timestamptz,
    constraint users_email_format_chk check (position('@' in email) > 1),
    constraint users_email_uniq unique (email),
    constraint users_username_uniq unique (username)
);

create table if not exists user_status (
    id bigserial primary key,
    user_id uuid not null references users(id) on delete cascade,
    from_status text
        check (from_status is null or from_status in ('pending', 'active', 'blocked', 'deleted')),
    to_status text not null
        check (to_status in ('pending', 'active', 'blocked', 'deleted')),
    reason text,
    changed_by uuid,
    changed_at timestamptz not null default now()
);

create index if not exists
    users_status_idx
on users (status);

create index if not exists
    users_created_at_idx
on users (created_at desc);

create index if not exists
    users_email_lower_idx
on users ((lower(email)));

create index if not exists
    user_status_user_changed_at_idx
on user_status (user_id, changed_at desc);

create index if not exists
    user_status_to_status_idx
on user_status (to_status);