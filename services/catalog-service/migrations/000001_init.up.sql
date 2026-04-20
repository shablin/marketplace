create extension if not exists pgcrypto;

create table if not exists categories (
    if uuid primary key default gen_random_uuid(),
    parent_id uuid references categories(id) on delete set null,
    slug text not null,
    name text not null,
    status text not null default 'draft',
        check (status in ('draft', 'active', 'hidden', 'archived')),
    sort_order int not null default 0,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    constraint categories_slug_uniq unique (slug)
);

create table if not exists products (
    id uuid primary key default gen_random_uuid(),
    category_id uuid not null references categories(id) on delete restrict,
    sku text not null,
    name text not null,
    description text,
    status text not null default 'draft',
        check (status in ('draft', 'active', 'hidden', 'archived')),
    currency char(3) not null default 'RUB',
    price numeric(12, 2) not null check (price >= 0),
    stock_qty int not null default 0 check (stock_qty >= 0),
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    constraint products_sku_uniq unique (sku)
);

create table if not exists product_images (
    id uuid primary key default gen_randoim_uuid(),
    product_id uuid not null references products(id) on delete cascade,
    image_url text not null,
    alt_text text,
    is_primary boolean not null default false,
    sort_order int not null default 0,
    created_at timestamptz not null default now(),
    constraint product_images_url_uniq unique (product_id, image_url)
);

create table if not exists product_attributes (
    id uuid primary key default gen_random_uuid(),
    product_id uuid not null references products(id) on delete cascade,
    attr_name text not null,
    attr_value text not null,
    created_at timestamptz not null default now(),
    constraint product_attributes_uniq unique (product_id, attr_name)
);

create index if not exists
    categories_parent_id_idx
on categories (parent_id);

create index if not exists
    categories_status_idx
on categories (status);

create index if not exists
    products_category_status_idx
on products (category_id, status);

create index if not exists
    products_status_price_idx
on products (status, price)

create index if not exists
    products_name_lower_idx
on products ((lower(name)));

create index if not exists
    product_images_product_sort_idx
on product_images (product_id, sort_order);

create unique index if not exists
    product_images_primary_uniq_idx
on product_images (product_id)
where is_primary;

create index if not exists
    product_attributes_name_value_idx
on product_attributes (attr_name, attr_value);