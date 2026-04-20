create extension if not exists pgcrypto;

create table if not exists notifications (
    id uuid primary key default gen_random_uuid(),
    user_id uuid not null,
    channel text not null check (channel in ('email', 'sms', 'push', 'webhook')),
    type text not null,
    status text not null default 'pending',
        check (status in ('pending', 'sent', 'delivered', 'failed', 'read')),
    subject text,
    message text not null,
    metadata jsonb,
    sent_at timestamptz,
    read_at timestamptz,
    created_at timestamptz not null default now()
);

create index if not exists
    notifications_user_status_idx
on notifications (user_id, status);

create index if not exists
    notifications_channel_status_idx
on notifications (channel, status);

create index if not exists
    notifications_type_created_idx
on notifications (type, created_at desc);

create index if not exists
    notifications_created_at_idx
on notifications (created_at desc);